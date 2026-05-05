package svc

import (
	"order_rpc/orderclient"
	"seckill/api/internal/config"
	"seckill/api/model"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config   config.Config
	Redis    *redis.Redis
	KqPusher *kq.Pusher
	OrderRpc orderclient.Order
	ActModel model.SeckillActivityModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:   c,
		Redis:    redis.New(c.Redis.Host, redis.WithPass(c.Redis.Pass)),
		KqPusher: kq.NewPusher(c.Kafka.Brokers, c.Kafka.Topic),
		OrderRpc: orderclient.NewOrder(zrpc.MustNewClient(c.OrderRpc)),
		ActModel: model.NewSeckillActivityModel(conn),
	}
}
