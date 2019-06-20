package util

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestId(t *testing.T) {
	now := time.Now()

	ai := New(1, 1)
	var lastId int
	var max = 1000000
	for i := 0 ; i < max ; i ++ {
		lastId = ai.Id()
	}

	assert.Equal(t, lastId, max)
	diff := time.Now().Sub(now)

	log.Println("use time : ", diff)
}
