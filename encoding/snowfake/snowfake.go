package snowfake

import (
	"github.com/bwmarrin/snowflake"
	"sync"
	"time"
)

var snow *snowflake.Node
var once sync.Once

// NewSnowflakeID 使用当前时间戳作为 Seed 创建 Snowflake 的实例
func NewSnowflakeID() int64 {
	once.Do(func() {
		snow, _ = snowflake.NewNode(time.Now().UnixNano())
	})
	// Generate a snowflake ID.
	id := snow.Generate().Int64()
	return id
}
