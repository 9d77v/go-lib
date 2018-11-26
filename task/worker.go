package task

import (
	"sync"
)

type empty struct{}

//Worker 带并发控制的任务执行器
type Worker struct {
	wg        *sync.WaitGroup
	tasks     []func()   //执行的任务
	max       int        //最大并发数
	limitChan chan empty //控制并发的channel
}

//NewWorker 新的worker
func NewWorker(tasks []func(), max int) *Worker {
	if max < 1 {
		max = 1
	}
	return &Worker{
		wg:        new(sync.WaitGroup),
		tasks:     tasks,
		max:       max,
		limitChan: make(chan empty, max),
	}
}

//Run 并发执行任务
func (w *Worker) Run() {
	for _, task := range w.tasks {
		w.limitChan <- empty{}
		w.wg.Add(1)
		go w.do(task)
	}
	w.wg.Wait()
}

//do 执行某个任务
func (w *Worker) do(task func()) {
	task()
	<-w.limitChan
	w.wg.Done()
}
