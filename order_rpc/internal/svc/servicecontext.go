package svc

import (
	"order_rpc/internal/config"
	"order_rpc/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config     config.Config
	OrderModel model.SeckillOrderModel
	Conn       sqlx.SqlConn
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:     c,
		OrderModel: model.NewSeckillOrderModel(conn),
		Conn:       conn,
	}
}
