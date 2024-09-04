package task

import (
	"github.com/gy/gosocket/internal/deque"
	"sync"
)

type Task func()

// TaskQueue Task 任务队列
type TaskQueue struct {
	lock  sync.Mutex
	queue deque.Deque[Task]
}

func (tq *TaskQueue) Push(t Task) {
	if t == nil {
		return
	}
	tq.lock.Lock()
	defer tq.lock.Unlock()
	tq.queue.PushBack(t)
}

func (tq *TaskQueue) Execute() {
	t := tq.getTask()
	if t != nil {
		go tq.execute(t)
	}
}

func (tq *TaskQueue) getTask() Task {
	tq.lock.Lock()
	defer tq.lock.Unlock()

	t, ok := tq.queue.PopFront()
	if !ok {
		return nil
	}
	return t
}

func (tq *TaskQueue) execute(t Task) {
	for t != nil {
		t()
		t = tq.getTask()
	}
}
