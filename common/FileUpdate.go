package common

import (
	"encoding/xml"
	"github.com/radovskyb/watcher"
)

const (
	StartService    = 1
	StopService     = 2
	ExecuteCommands = 3
)

type FileUpdate struct {
	Name               string
	RelativePath       string
	Operation          watcher.Op
	Content            []byte
	ExecutableCommands []string
	PiperOpcode        byte
}

func FromJson(json []byte) FileUpdate {
	var update FileUpdate
	xml.Unmarshal(json, &update)
	return update
}

func (update FileUpdate) ToJson() string {
	e, _ := xml.Marshal(update)
	return string(e)
}
