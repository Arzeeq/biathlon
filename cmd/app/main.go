package main

import (
	"biathlon/internal/competition"
	"biathlon/internal/config"
	"biathlon/internal/logger"
	"flag"
	"log"
	"os"
)

func main() {
	cfgPath := flag.String("config", "./configs/config.json", "path to your json config file")
	eventsPath := flag.String("events", "./events/events", "path to your events file")
	logPath := flag.String("log", "output.log", "path to log file")
	reportPath := flag.String("report", "report.txt", "path to report file")
	flag.Parse()

	if cfgPath == nil || eventsPath == nil || logPath == nil || reportPath == nil {
		log.Fatal("failed to parse arguments")
	}

	// load config
	cfg, err := config.MustLoad(*cfgPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// create logger
	logFile, err := os.Create(*logPath)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}

	l, err := logger.New(logFile)
	if err != nil {
		log.Fatalf("failed to create event logger: %v", err)
	}

	// open events eventsFile
	eventsFile, err := os.Open(*eventsPath)
	if err != nil {
		log.Fatalf("failed to open events: %v", err)
	}

	// create and start competition
	comp, err := competition.New(eventsFile, l, *cfg)
	if err != nil {
		log.Fatalf("failed to create competition: %v", err)
	}

	if err := comp.Start(); err != nil {
		log.Fatalf("failed to process incoming events: %v", err)
	}

	// generate report
	reportFile, err := os.Create(*reportPath)
	if err != nil {
		log.Fatalf("failed to open report file: %v", err)
	}

	report, err := comp.GenerateReport()
	if err != nil {
		log.Fatalf("failed to generate report: %v", err)
	}

	_, err = reportFile.WriteString(report)
	if err != nil {
		log.Fatalf("failed to write report: %v", err)
	}
}
