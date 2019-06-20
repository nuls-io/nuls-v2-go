package goroutine

import (
	"fmt"
	"testing"
)

func TestExecute(t *testing.T) {

	pool := NewPool(5)
	go pool.Run()

	for i := 0; i < 10; i++ {
		task := NewTask(func() {
			fmt.Println("this is num of ", i)
		})
		pool.Execute(task)
	}
}
