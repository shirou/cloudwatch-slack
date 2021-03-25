package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
)

const (
	ColorGood   = "#2EB886"
	ColorWarn   = "#DAA038"
	ColorDanger = "#A30100"
)

type Event interface {
	readTemplate() (*template.Template, error)
	genMessage(*template.Template) (string, error)
}

func GenerateMessage(source string, detailType string, detail json.RawMessage) (string, error) {
	event, err := NewEvent(source, detailType, detail)
	if err != nil {
		return "", fmt.Errorf("NewEvent, %w", err)
	}
	tmpl, err := event.readTemplate()
	if err != nil {
		return "", fmt.Errorf("readTemplate, %w", err)
	}
	return event.genMessage(tmpl)
}

func NewEvent(source, detailType string, detail json.RawMessage) (Event, error) {
	switch source {
	case "aws.cloudwatch":
		if detailType == "CloudWatch Alarm State Change" {
			return NewEventCloudWatchAlarm(source, detailType, detail), nil
		}
	}
	return nil, fmt.Errorf("can not find matched template: %s, %s", source, detailType)
}

func readTemplate(paths ...string) (*template.Template, error) {
	// root := filepath.Dir(os.Args[0])

	path := []string{"templates"}
	for _, p := range paths {
		path = append(path, p)
	}
	p := strings.Join(path, "/")
	return template.ParseFiles(p)
}
