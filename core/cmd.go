package core

type RedisCmd struct {
	Cmd  string   // PUT
	Args []string // K V
}

type RedisCmds []*RedisCmd
