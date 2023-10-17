package predis

import "github.com/gomodule/redigo/redis"

// NewCommand
//
//	@Description: 构建redis命令
//	@return unc NewCommand(name string, args []interface{},
func NewCommand(name string, args []interface{},
	replyHandler func(err error, reply interface{}, args ...interface{}), replyArgs []interface{}) *Command {
	return &Command{Name: name, Args: args, ReplyHandler: replyHandler, ReplyArgs: replyArgs}
}

// NewScriptCommand
//
//	@Description: 构建redis脚本命令
//	@return unc NewScriptCommand(script *redis.Script, args []interface{},
func NewScriptCommand(script *redis.Script, args []interface{},
	replyHandler func(err error, reply interface{}, args ...interface{}), replyArgs []interface{}) *Command {
	return &Command{Script: script, Args: args, ReplyHandler: replyHandler, ReplyArgs: replyArgs}
}

type Command struct {
	Name         string
	Script       *redis.Script
	Args         []interface{}
	ReplyHandler func(err error, reply interface{}, args ...interface{})
	ReplyArgs    []interface{}
}
