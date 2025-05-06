package competition

import (
	"biathlon/internal/config"
	"biathlon/internal/event"
	"biathlon/internal/logger"
	"bufio"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
)

type Competition struct {
	config.Config
	r             io.Reader
	scheduled     []event.Event
	comp          map[int]Competitor
	l             *logger.EventLogger
	eventHandlers map[event.EventType]func(*event.Event, *Competitor) error
}

func New(r io.Reader, l *logger.EventLogger, cfg config.Config) (*Competition, error) {
	if r == nil || l == nil {
		return nil, errors.New("nil values in constructor")
	}

	c := Competition{r: r, comp: make(map[int]Competitor), Config: cfg, l: l}

	eventHandlers := map[event.EventType]func(*event.Event, *Competitor) error{
		event.Registered:      c.processRegistered,
		event.SetStartTime:    c.processSetStartTime,
		event.OnStartLine:     c.processOnStartLine,
		event.Started:         c.processStarted,
		event.OnFiringRange:   c.processOnFiringRange,
		event.TargetHit:       c.processTargetHit,
		event.LeftFiringRange: c.processLeftFiringRange,
		event.EnteredPenalty:  c.processEnteredPenalty,
		event.LeftPenalty:     c.processLeftPenalty,
		event.EndedMain:       c.processEndedMain,
		event.CanNotContinue:  c.processCanNotContinue,
		event.Disqualified:    c.processDisqualified,
	}
	c.eventHandlers = eventHandlers

	return &c, nil
}

func (c *Competition) Start() error {
	sc := bufio.NewScanner(c.r)

	for sc.Scan() {
		e, err := event.New(sc.Text())
		if err != nil {
			return err
		}

		for len(c.scheduled) > 0 && c.scheduled[0].Time.Before(e.Time) {
			if err := c.ProcessEvent(&c.scheduled[0]); err != nil {
				return err
			}

			c.scheduled = c.scheduled[1:]
		}

		if err := c.ProcessEvent(e); err != nil {
			return err
		}
	}

	if err := sc.Err(); err != nil {
		return err
	}

	return nil
}

func (c *Competition) ProcessEvent(e *event.Event) error {
	if e == nil {
		return errors.New("nil event")
	}

	eventHandler, isFound := c.eventHandlers[e.ID]
	if !isFound {
		return fmt.Errorf("unknown eventID %d", e.ID)
	}

	comp, exists := c.comp[e.Competitor]
	if !exists && e.ID != event.Registered {
		return fmt.Errorf("unknown competitor %d", e.Competitor)
	}

	if err := eventHandler(e, &comp); err != nil {
		return err
	}
	c.comp[e.Competitor] = comp

	return nil
}

func (c *Competition) GenerateReport() (string, error) {
	comps := make([]Competitor, 0, len(c.comp))

	for _, comp := range c.comp {
		comps = append(comps, comp)
	}

	sort.Slice(comps, func(i, j int) bool {
		return comps[i].TotalTime() < comps[j].TotalTime()
	})

	var b strings.Builder
	for i := range comps {
		report, err := comps[i].Report(c.Config)
		if err != nil {
			return "", err
		}

		b.WriteString(report)
		b.WriteRune('\n')
	}
	return b.String(), nil
}
