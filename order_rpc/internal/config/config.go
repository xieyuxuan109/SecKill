package config

import "github.com/zeromicro/go-zero/zrpc"

type MysqlConf struct {
	DataSource string
}

type Config struct {
	zrpc.RpcServerConf
	Mysql MysqlConf
}
