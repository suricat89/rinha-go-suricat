package utils

import (
	"time"

	"github.com/gofiber/fiber/v3/log"
)

type TraceTime struct {
	module    string
	function  string
	operation string
	startTime time.Time
}

func NewTraceTime(module string, function string) *TraceTime {
	return &TraceTime{
		module:    module,
		function:  function,
		operation: "",
		startTime: time.Now(),
	}
}

func (t *TraceTime) Start(operation string) {
	t.operation = operation
	t.startTime = time.Now()
}

func (t *TraceTime) End() {
	log.Tracef(
		"[%s/%s] %s: %v",
		t.module,
		t.function,
		t.operation,
		time.Since(t.startTime),
	)
}
