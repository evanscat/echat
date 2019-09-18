package echat

import (
	"github.com/golang/glog"
	"time"
)

type Logger interface {
	Log([]byte) error
	Query(id string, from time.Time, to time.Time) [][]byte
}

type DefaultLogger struct {
}

func (*DefaultLogger) Log(msg []byte) error {
	glog.Info(string(msg))
	return nil
}

func (*DefaultLogger) Query(id string, from time.Time, to time.Time) [][]byte {
	return nil
}
