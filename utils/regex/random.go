package regex

import (
	"encoding/hex"
	"math/rand"
	"sync"
	"time"
)

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

// RandomHex 生成随机字符串
//
// Deprecated: 该接口后续将废弃!!
func RandomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GenerateHex(n int) string {
	bytes := make([]byte, n)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
