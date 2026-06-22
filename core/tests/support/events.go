package support

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

const (
	EventSuiteStart    = "suite_start"
	EventSuiteEnd      = "suite_end"
	EventScenarioStart = "scenario_start"
	EventScenarioEnd   = "scenario_end"
	EventStep          = "step"
	EventHTTP          = "http"
)

type Event struct {
	Type      string      `json:"type"`
	Timestamp string      `json:"timestamp"`
	Data      interface{} `json:"data"`
}

var (
	eventsFile *os.File
	eventsMu   sync.Mutex
	eventsOn   = os.Getenv("EVENTS_FILE") != ""
)

func InitEvents() {
	if !eventsOn {
		return
	}
	var err error
	eventsFile, err = os.OpenFile(os.Getenv("EVENTS_FILE"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		eventsOn = false
	}
}

func CloseEvents() {
	if eventsFile != nil {
		eventsFile.Close()
	}
}

func EmitEvent(eventType string, data interface{}) {
	if !eventsOn {
		return
	}
	eventsMu.Lock()
	defer eventsMu.Unlock()

	event := Event{
		Type:      eventType,
		Timestamp: time.Now().Format(time.RFC3339Nano),
		Data:      data,
	}
	line, _ := json.Marshal(event)
	eventsFile.Write(append(line, '\n'))
}
