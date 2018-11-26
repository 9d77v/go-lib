package task

import (
	"testing"
)

func TestWorker_Wait(t *testing.T) {
	var p = func(i int) {
		println(i)
	}
	var tasks = make([]func(), 0, 0)
	for i := 0; i < 10; i++ {
		t := i
		tasks = append(tasks, func() {
			p(t)
		})
	}
	worker := NewWorker(tasks, 5)
	worker.Run()
	t.Fail()
}
