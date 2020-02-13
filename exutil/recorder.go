package exutil

import (
	"fmt"
	"sync"
)

// Recorder interface for implementing various data recording, like kafka, pusher....
type Recorder interface {
	Record(data interface{}) error
}

// GenericRecorder is a general implement of Recorder interface
type GenericRecorder struct {
	sync.Mutex
}

// NewGenericRecorder generates a GenericRecorder instance
func NewGenericRecorder() Recorder {
	rec := &GenericRecorder{}
	return rec
}

// Record is an method implement for GenericRecorder
func (rec *GenericRecorder) Record(data interface{}) error {
	rec.Lock()
	fmt.Printf("record data %v\n", data)
	rec.Unlock()
	return nil
}
