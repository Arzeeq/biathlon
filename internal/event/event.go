package event

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const timeFormat = "15:04:05.000"

type EventType int

const (
	Registered EventType = iota + 1
	SetStartTime
	OnStartLine
	Started
	OnFiringRange
	TargetHit
	LeftFiringRange
	EnteredPenalty
	LeftPenalty
	EndedMain
	CanNotContinue
)

const (
	Disqualified EventType = iota + 32
)

type Event struct {
	Time       time.Time
	ID         EventType
	Competitor int
	Extra      string
}

func New(desc string) (*Event, error) {
	strs := strings.Split(desc, " ")

	if len(strs) < 3 {
		return nil, errors.New("not enough params")
	}

	t, err := time.Parse(timeFormat, strings.Trim(strs[0], "[]"))
	if err != nil {
		return nil, err
	}

	id, err := strconv.Atoi(strs[1])
	if err != nil {
		return nil, err
	}

	competitor, err := strconv.Atoi(strs[2])
	if err != nil {
		return nil, err
	}

	event := Event{
		Time:       t,
		ID:         EventType(id),
		Competitor: competitor,
	}

	if len(strs) > 3 {
		event.Extra = strings.Join(strs[3:], " ")
	}

	return &event, nil
}
