package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"
)

type EventCloudWatchAlarmMessage struct {
	Color string

	AlarmName string `json:"alarmName"`
	State     struct {
		Value      string `json:"value"`
		Reason     string `json:"reason"`
		ReasonData string `json:"reasonData"`
		Timestamp  string `json:"timestamp"`
	} `json:"state"`
	PreviousState struct {
		Value      string `json:"value"`
		Reason     string `json:"reason"`
		ReasonData string `json:"reasonData"`
		Timestamp  string `json:"timestamp"`
	} `json:"previousState"`
	Configuration struct {
		Description string `json:"description"`
		Metrics     []struct {
			ID         string `json:"id"`
			MetricStat struct {
				Metric struct {
					Namespace  string `json:"namespace"`
					Name       string `json:"name"`
					Dimensions struct {
						InstanceID string `json:"InstanceId"`
					} `json:"dimensions"`
				} `json:"metric"`
				Period int    `json:"period"`
				Stat   string `json:"stat"`
			} `json:"metricStat,omitempty"`
			ReturnData bool   `json:"returnData"`
			Expression string `json:"expression,omitempty"`
			Label      string `json:"label,omitempty"`
		} `json:"metrics"`
	} `json:"configuration"`
}

type EventCloudWatchAlarm struct {
	Source     string
	DetailType string
	Detail     json.RawMessage
	Parsed     EventCloudWatchAlarmMessage
}

func NewEventCloudWatchAlarm(source string, detailType string, detail json.RawMessage) EventCloudWatchAlarm {
	return EventCloudWatchAlarm{
		Source:     source,
		DetailType: detailType,
		Detail:     detail,
	}
}

func (m EventCloudWatchAlarm) readTemplate() (*template.Template, error) {
	return readTemplate("alarm", "change.json")
}

func (m EventCloudWatchAlarm) genMessage(tmpl *template.Template) (string, error) {
	if err := json.Unmarshal(m.Detail, &m.Parsed); err != nil {
		return "", fmt.Errorf("Unmarshal, %w", err)
	}
	m.Parsed.Color = m.getColor(m.Parsed.State.Value)

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, m.Parsed); err != nil {
		return "", fmt.Errorf("Execute, %w", err)
	}

	return buf.String(), nil
}

func (m EventCloudWatchAlarm) getColor(value string) string {
	switch value {
	case "ALARM":
		return ColorDanger
	case "OK":
		return ColorGood
	case "INSUFFICIENT_DATA":
		return ColorWarn
	default:
		return ColorWarn
	}
}
