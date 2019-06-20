package goroutine

import (
	"github.com/nuls-io/go-nuls/nuls-frame/util"
	"time"
)

type Task struct {
	f func()
}

func NewTask(f func()) *Task {
	t := Task{
		f : f,
	}
	return &t
}

func (t *Task) execute() {
	t.f()
}

type Pool struct {
	// 协程池最大worker数量,限定Goroutine的个数
	workerNum 	int
	// 协程池内部的任务队列
	taskQueue 	chan *Task
	// 对外接收Task的入口
	workChannel chan chan *Task
	// 停止标志
	quit		chan bool
}

func NewPool(cap int) *Pool {
	pool := Pool{
		workerNum: cap,
		taskQueue: make(chan *Task),
		workChannel: make(chan chan *Task),
	}
	return &pool
}

func (p *Pool) worker(workID int) {
	go func() {
		for {
			p.workChannel <- p.taskQueue
			select {
			case task := <- p.taskQueue:
				task.execute()
			case <- p.quit:
				break
			}
		}
	}()
}

func (p *Pool) Run() {
	for i := 0; i < p.workerNum; i++ {
		p.worker(i)
	}
}

func (p *Pool) Execute(task *Task) {

	now := time.Now()
	taskWork := <- p.workChannel
	util.T2 += time.Now().Sub(now)
	now = time.Now()
	taskWork <- task
	util.T3 += time.Now().Sub(now)
}

func (p *Pool) Shutdown() {
	p.quit <- true
}