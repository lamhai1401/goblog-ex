package util

import "github.com/bwmarrin/snowflake"

var snowflakeGenerator *snowflake.Node

// InitSnowflakeGenerator linter
func InitSnowflakeGenerator(id int64) error {
	var err error
	snowflakeGenerator, err = snowflake.NewNode(id)
	if err != nil {
		return err
	}
	return nil
}

// GetSnowflakeID linter
func GetSnowflakeID() int64 {
	return snowflakeGenerator.Generate().Int64()
}

// ParseSnowflakeID litner
func ParseSnowflakeID(data string) (snowflake.ID, error) {
	id, err := snowflake.ParseString(data)
	if err != nil {
		return 0, err
	}
	return id, nil
}
