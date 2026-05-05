package logic

import (
	"context"
	"database/sql"

	"order_rpc/internal/svc"
	"order_rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderLogic {
	return &GetOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOrderLogic) GetOrder(in *order.GetOrderReq) (*order.GetOrderResp, error) {
	// todo: add your logic here and delete this line
	var result struct {
		OrderNo    string `db:"order_no"`
		UserId     int64  `db:"user_id"`
		ActId      int64  `db:"act_id"`
		Status     int32  `db:"status"`
		CreateTime int64  `db:"create_time"`
	}
	query := `SELECT order_no, user_id, act_id, status, UNIX_TIMESTAMP(create_time) as create_time 
              FROM seckill_order WHERE order_no = ?`
	err := l.svcCtx.Conn.QueryRow(&result, query, in.OrderNo)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &order.GetOrderResp{
		OrderNo:    result.OrderNo,
		UserId:     result.UserId,
		ActId:      result.ActId,
		Status:     result.Status,
		CreateTime: result.CreateTime,
	}, nil
}
