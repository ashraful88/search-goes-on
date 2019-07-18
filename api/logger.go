package api

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

// StructuredLog struct for structured loggin
type StructuredLog struct {
	Timestamp string `json:"@timestamp,omitempty"`
	Service   string `json:"service,omitempty"`
	Thread    string `json:"thread,omitempty"`
	IP        string `json:"ip,omitempty"`
	Env       string `json:"env,omitempty"`
	Server    string `json:"server,omitempty"`

	Level        string      `json:"level,omitempty"`
	Event        string      `json:"event,omitempty"`
	Message      string      `json:"message,omitempty"`
	Account      string      `json:"account,omitempty"`
	ID           string      `json:"id,omitempty"`
	Raw          string      `json:"raw,omitempty"`
	RawInterface interface{} `json:"rawInterface,omitempty"`
}

// LogEvent in json format
func LogEvent(thelog StructuredLog, level, event, msg string) {
	thelog.Timestamp = time.Now().Format(time.RFC3339)
	hostname, _ := os.Hostname()
	thelog.Server = hostname
	thelog.Level = level
	thelog.Event = event
	thelog.Message = msg
	//thelog.Service = "search-api"
	logJSON, err := json.Marshal(thelog)
	if err != nil {
		log.Println("Structured Logger: Logger JSON Marshal failed !")
	}
	log.Println(string(logJSON))
}

// LogNew is to log with a new StructuredLog struct
func LogNew(level, event, msg string) {
	var thelog StructuredLog
	thelog.Timestamp = time.Now().Format(time.RFC3339)
	hostname, _ := os.Hostname()
	thelog.Server = hostname
	thelog.Level = level
	thelog.Event = event
	thelog.Message = msg
	thelog.Service = "search-api"
	logJSON, err := json.Marshal(thelog)
	if err != nil {
		log.Println("Structured logger: Logger JSON Marshal failed !")
	}
	log.Println(string(logJSON))
}
