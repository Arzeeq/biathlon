package competition

import (
	"biathlon/internal/event"
	"fmt"
	"strconv"
	"time"
)

const timeFormat = "15:04:05.000"

func (c *Competition) processRegistered(e *event.Event, comp *Competitor) error {
	comp.ID = e.Competitor
	comp.Status = NotStarted

	msg := fmt.Sprintf("The competitor(%d) registered", e.Competitor)
	c.l.Log(e.Time, msg)
	return nil
}

func (c *Competition) processSetStartTime(e *event.Event, comp *Competitor) error {
	startTime, err := time.Parse(timeFormat, e.Extra)
	if err != nil {
		return err
	}

	comp.ScheduledStart = startTime

	c.scheduled = append(c.scheduled, event.Event{
		ID:         event.Disqualified,
		Time:       startTime.Add(time.Duration(c.StartDelta)),
		Competitor: e.Competitor,
	})

	msg := fmt.Sprintf("The start time for the competitor(%d) was set by a draw to %s",
		e.Competitor,
		startTime.Format(timeFormat))
	c.l.Log(e.Time, msg)

	return nil
}

func (c *Competition) processOnStartLine(e *event.Event, comp *Competitor) error {
	msg := fmt.Sprintf("The competitor(%d) is on the start line", e.Competitor)
	c.l.Log(e.Time, msg)
	return nil
}

func (c *Competition) processStarted(e *event.Event, comp *Competitor) error {
	comp.Status = Started
	comp.Laps = append(comp.Laps, Lap{Start: comp.ScheduledStart})

	msg := fmt.Sprintf("The competitor(%d) has started", e.Competitor)
	c.l.Log(e.Time, msg)

	return nil
}

func (c *Competition) processOnFiringRange(e *event.Event, comp *Competitor) error {
	firingRange, err := strconv.Atoi(e.Extra)
	if err != nil {
		return err
	}

	comp.Hits = append(comp.Hits, 0)
	comp.Penalty = append(comp.Penalty, Lap{})

	msg := fmt.Sprintf("The competitor(%d) is on the firing range(%d)", e.Competitor, firingRange)
	c.l.Log(e.Time, msg)

	return nil
}

func (c *Competition) processTargetHit(e *event.Event, comp *Competitor) error {
	target, err := strconv.Atoi(e.Extra)
	if err != nil {
		return err
	}

	comp.Hits[len(comp.Hits)-1]++

	msg := fmt.Sprintf("The target(%d) has been hit by competitor(%d)", target, e.Competitor)
	c.l.Log(e.Time, msg)

	return nil
}

func (c *Competition) processLeftFiringRange(e *event.Event, comp *Competitor) error {
	msg := fmt.Sprintf("The competitor(%d) left the firing range", e.Competitor)
	c.l.Log(e.Time, msg)

	return nil
}

func (c *Competition) processEnteredPenalty(e *event.Event, comp *Competitor) error {
	comp.Penalty[len(comp.Penalty)-1].Start = e.Time

	msg := fmt.Sprintf("The competitor(%d) entered the penalty laps", e.Competitor)
	c.l.Log(e.Time, msg)

	return nil
}

func (c *Competition) processLeftPenalty(e *event.Event, comp *Competitor) error {
	comp.Penalty[len(comp.Penalty)-1].End = e.Time

	msg := fmt.Sprintf("The competitor(%d) left the penalty laps", e.Competitor)
	c.l.Log(e.Time, msg)

	return nil
}

func (c *Competition) processEndedMain(e *event.Event, comp *Competitor) error {
	comp.Laps[len(comp.Laps)-1].End = e.Time

	var msg string
	if len(comp.Laps) < c.Laps {
		comp.Laps = append(comp.Laps, Lap{Start: e.Time})
		msg = fmt.Sprintf("The competitor(%d) ended the main lap", e.Competitor)
	} else {
		comp.Status = Finished
		msg = fmt.Sprintf("The competitor(%d) has finished", e.Competitor)
	}

	c.l.Log(e.Time, msg)

	return nil
}

func (c *Competition) processCanNotContinue(e *event.Event, comp *Competitor) error {
	comp.Status = NotFinished

	msg := fmt.Sprintf("The competitor(%d) can`t continue: %s", e.Competitor, e.Extra)
	c.l.Log(e.Time, msg)

	return nil
}

func (c *Competition) processDisqualified(e *event.Event, comp *Competitor) error {
	if comp.Status != NotStarted {
		return nil
	}

	msg := fmt.Sprintf("The competitor(%d) is disqualified", e.Competitor)
	c.l.Log(e.Time, msg)

	return nil
}
