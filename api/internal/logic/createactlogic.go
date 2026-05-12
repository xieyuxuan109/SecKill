// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"fmt"
	"seckill/api/internal/svc"
	"seckill/api/internal/types"
	"seckill/api/model"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateActLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateActLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateActLogic {
	return &CreateActLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateActLogic) CreateAct(req *types.CreateActReq) (resp *types.CreateActResp, err error) {
	// 1. 写入MySQL
	act := &model.SeckillActivity{
		Name:    req.Name,
		Stock:   int64(req.Stock),
		StartAt: req.StartAt,
		EndAt:   req.EndAt,
	}
	result, err := l.svcCtx.ActModel.Insert(l.ctx, act)
	if err != nil {
		return nil, err
	}
	actId, _ := result.LastInsertId()
	// 2. 初始化 Redis 库存
	stockKey := fmt.Sprintf("stock:act:%d", actId)
	l.svcCtx.Redis.Set(stockKey, strconv.Itoa(req.Stock))
	// 3. 存储活动时间（转换为毫秒级）
	timeKey := fmt.Sprintf("act:time:%d", actId)
	l.svcCtx.Redis.Hset(timeKey, "start", strconv.FormatInt(req.StartAt*1000, 10))
	l.svcCtx.Redis.Hset(timeKey, "end", strconv.FormatInt(req.EndAt*1000, 10))
	return &types.CreateActResp{}, nil
}
