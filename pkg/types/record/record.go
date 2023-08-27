package record

import (
	"encoding/json"
	"time"
)

type Record struct {
	Hash      string    `json:"hash"`
	Version   int       `json:"version"`
	Path      string    `json:"path"`
	Diff      string    `json:"diff"`
	Time      time.Time `json:"time"`
	Editor    string    `json:"editor"`
	Size      uint64    `json:"size"`
	Insert    uint64    `json:"insert"`
	Delete    uint64    `json:"delete"`
	Extension string    `json:"extension"`
}

func NewFromJson(j []byte) *Record {
	record := &Record{}
	err := json.Unmarshal(j, record)
	if err != nil {
		panic(err)
	}
	return record
}

func (l *Record) MarshalJson() []byte {
	json, err := json.Marshal(l)
	if err != nil {
		panic(err)
	}
	return json
}
