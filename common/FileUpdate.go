package common

import "github.com/radovskyb/watcher"

type FileUpdate struct {
	Name string
	RelativePath string
	Operation watcher.Op
	Content string
}
