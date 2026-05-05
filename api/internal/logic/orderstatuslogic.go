// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"fmt"
	"order_rpc/orderclient"
	"seckill/api/internal/svc"
	"seckill/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderStatusLogic {
	return &OrderStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderStatusLogic) OrderStatus(req *types.OrderStatusReq) (resp *types.OrderStatusResp, err error) {
	// todo: add your logic here and delete this line
	rpcResp, err := l.svcCtx.OrderRpc.GetOrder(l.ctx, &orderclient.GetOrderReq{
		OrderNo: req.OrderNo,
	})
	if err != nil {
		return nil, err
	}
	if rpcResp == nil {
		return nil, fmt.Errorf("订单不存在")
	}
	return &types.OrderStatusResp{
		OrderNo: rpcResp.OrderNo,
		Status:  rpcResp.Status,
	}, nil
}
