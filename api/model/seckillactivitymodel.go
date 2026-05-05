package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SeckillActivityModel = (*customSeckillActivityModel)(nil)

type (
	// SeckillActivityModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSeckillActivityModel.
	SeckillActivityModel interface {
		seckillActivityModel
		withSession(session sqlx.Session) SeckillActivityModel
	}

	customSeckillActivityModel struct {
		*defaultSeckillActivityModel
	}
)

// NewSeckillActivityModel returns a model for the database table.
func NewSeckillActivityModel(conn sqlx.SqlConn) SeckillActivityModel {
	return &customSeckillActivityModel{
		defaultSeckillActivityModel: newSeckillActivityModel(conn),
	}
}

func (m *customSeckillActivityModel) withSession(session sqlx.Session) SeckillActivityModel {
	return NewSeckillActivityModel(sqlx.NewSqlConnFromSession(session))
}
