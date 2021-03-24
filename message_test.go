package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestGenerateMessage(t *testing.T) {
	tests := []struct {
		filePath string
		output   string
	}{
		{"testdata/alarm/single_metrics.json", ""},
		{"testdata/alarm/multi_metrics.json", ""},
		{"testdata/alarm/anomary.json", ""},
	}

	for _, test := range tests {
		jsonFromFile, err := ioutil.ReadFile(test.filePath)
		if err != nil {
			t.Fatal(err)
		}
		var event CloudwatchEvent
		if err := json.Unmarshal(jsonFromFile, &event); err != nil {
			t.Fatal(err)
		}
		out, err := GenerateMessage(event.Source, event.DetailType, event.Detail)
		if err != nil {
			t.Error(err)
		}
		t.Log(out)
	}
}
