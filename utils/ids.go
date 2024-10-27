package utils

import (
	"github.com/bwmarrin/snowflake"
	"math/rand"
	"sync"
	"time"
)

func NewSnowflakeID() int64 {
	node, _ := snowflake.NewNode(1)
	// Generate a snowflake ID.
	id := node.Generate().Int64()
	return id
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var R *rand.Rand = nil
var initR = sync.Once{}

func RandomString(length int) string {

	initR.Do(func() {
		R = rand.New(rand.NewSource(time.Now().UnixNano()))
	})

	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[R.Intn(len(letterBytes))]
	}
	return string(b)
}
