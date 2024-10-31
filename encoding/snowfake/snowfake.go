package snowfake

import (
	"github.com/bwmarrin/snowflake"
)

func NewSnowflakeID() int64 {
	node, _ := snowflake.NewNode(1)
	// Generate a snowflake ID.
	id := node.Generate().Int64()
	return id
}
