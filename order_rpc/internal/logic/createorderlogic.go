package logic

import (
	"context"
	"strings"

	"order_rpc/internal/model"
	"order_rpc/internal/svc"
	"order_rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateOrderLogic) CreateOrder(in *order.CreateOrderReq) (*order.CreateOrderResp, error) {
	// todo: add your logic here and delete this line
	_, err := l.svcCtx.OrderModel.Insert(l.ctx, &model.SeckillOrder{
		OrderNo: in.OrderNo,
		UserId:  in.UserId,
		ActId:   in.ActId,
		Status:  int64(in.Status),
	})
	if err != nil && strings.Contains(err.Error(), "Duplicate") {
		// 幂等：重复订单号视为成功
		return &order.CreateOrderResp{Success: true}, nil
	}
	if err != nil {
		return nil, err
	}
	return &order.CreateOrderResp{Success: true}, nil
}
