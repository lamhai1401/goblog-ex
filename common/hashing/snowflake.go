package hashing

import "github.com/bwmarrin/snowflake"

var idGenerator *snowflake.Node

// InitIDGenerator linter
func InitIDGenerator(numNode int) {
	idGenerator, _ = snowflake.NewNode(int64(numNode))
}

// GenID return in64
func GenID() int64 {
	return idGenerator.Generate().Int64()
}
