package competition

import (
	"biathlon/internal/config"
	"fmt"
	"strings"
	"time"
)

// suppose that any competition can not last more than a year
const InfTime = 365 * 24 * time.Hour

type Lap struct {
	Start, End time.Time
}

func (l *Lap) Time() time.Duration {
	return l.End.Sub(l.Start)
}

type Status string

const (
	NotStarted  Status = "NotStarted"
	Started     Status = "Started"
	NotFinished Status = "NotFinished"
	Finished    Status = "Finished"
)

type Competitor struct {
	ID             int
	Status         Status
	ScheduledStart time.Time
	Laps           []Lap
	Penalty        []Lap
	Hits           []int
}

func (c *Competitor) TotalTime() time.Duration {
	if c.Status != Finished || len(c.Laps) == 0 {
		return InfTime
	}

	return c.Laps[len(c.Laps)-1].End.Sub(c.ScheduledStart)
}

func (c *Competitor) Report(cfg config.Config) (string, error) {
	var totalTime string
	if c.Status == Finished {
		totalTime = formatDuration(c.TotalTime())
	} else {
		totalTime = string(c.Status)
	}

	reportLaps := c.reportLaps(cfg.Laps, cfg.LapLen)
	reportPenalty := c.reportPenalty(cfg.PenaltyLen)
	reportHits := c.reportHits(cfg.FiringLines)

	return fmt.Sprintf("[%s] %d %s %s %s", totalTime, c.ID, reportLaps, reportPenalty, reportHits), nil
}

func formatDuration(t time.Duration) string {
	return time.Time{}.Add(t).Format(timeFormat)
}

func (c *Competitor) reportLaps(laps, lapLen int) string {
	var b strings.Builder
	b.WriteByte('[')
	for i := range laps {
		if i < len(c.Laps) && c.Laps[i].End != (time.Time{}) {
			t := c.Laps[i].End.Sub(c.Laps[i].Start)
			speed := float64(lapLen) / t.Seconds()
			b.WriteString(fmt.Sprintf("{%s, %.3f}", formatDuration(t), speed))
		} else {
			b.WriteString("{,}")
		}
		if i != len(c.Laps)-1 {
			b.WriteString(", ")
		}
	}
	b.WriteByte(']')

	return b.String()
}

func (c *Competitor) reportPenalty(penaltyLen int) string {
	var penaltyTime time.Duration
	for i := range c.Penalty {
		if c.Penalty[i].End != (time.Time{}) {
			penaltyTime += c.Penalty[i].Time()
		}
	}

	lapsCount := 5*len(c.Penalty) - c.hits()
	var speed float64
	if lapsCount != 0 {
		speed = float64(penaltyLen*lapsCount) / penaltyTime.Seconds()
	}
	return fmt.Sprintf("{%s, %.3f}", formatDuration(penaltyTime), speed)
}

func (c *Competitor) reportHits(firingRanges int) string {
	return fmt.Sprintf("%d/%d", c.hits(), 5*firingRanges)
}

func (c *Competitor) hits() int {
	hits := 0
	for _, hit := range c.Hits {
		hits += hit
	}
	return hits
}
