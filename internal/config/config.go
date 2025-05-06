package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const timeFormat = "15:04:05.000"

type StartTime time.Time
type DeltaDuration time.Duration

func (t *StartTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}

	tmp, err := time.Parse(timeFormat, string(bytes.Trim(data, "\"")))
	if err != nil {
		return err
	}

	*t = StartTime(tmp)
	return nil
}

func (t *DeltaDuration) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}

	strs := strings.Split(string(bytes.Trim(data, "\"")), ":")
	if len(strs) != 3 {
		return errors.New("format does not match with expected (15:04:05)")
	}

	hours, err := strconv.Atoi(strs[0])
	if err != nil {
		return err
	}
	minutes, err := strconv.Atoi(strs[1])
	if err != nil {
		return err
	}
	seconds, err := strconv.Atoi(strs[2])
	if err != nil {
		return err
	}

	*t = DeltaDuration(
		time.Duration(hours)*time.Hour +
			time.Duration(minutes)*time.Minute +
			time.Duration(seconds)*time.Second)
	return nil
}

type Config struct {
	Laps        int           `json:"laps"`
	LapLen      int           `json:"lapLen"`
	PenaltyLen  int           `json:"penaltyLen"`
	FiringLines int           `json:"firingLines"`
	Start       StartTime     `json:"start"`
	StartDelta  DeltaDuration `json:"startDelta"`
}

func MustLoad(name string) (*Config, error) {
	var config Config
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalf("failed to close %s: %v", name, err)
		}
	}()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	return &config, err
}
