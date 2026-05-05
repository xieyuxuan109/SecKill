package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"seckill/api/internal/svc"
	"seckill/api/internal/types"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type SeckillLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSeckillLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SeckillLogic {
	return &SeckillLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func generateOrderNo() string {
	return fmt.Sprintf("%d%d", time.Now().UnixNano(), rand.Intn(10000))
}

func (l *SeckillLogic) Seckill(req *types.SeckillReq) (resp *types.SeckillResp, err error) {
	// 模拟用户ID
	userId := int64(1001)
	actId := req.ActivityId

	// 1. 活动时间校验
	timeKey := fmt.Sprintf("act:time:%d", actId)
	startStr, _ := l.svcCtx.Redis.Hget(timeKey, "start")
	endStr, _ := l.svcCtx.Redis.Hget(timeKey, "end")
	now := time.Now().Unix()
	start, _ := strconv.ParseInt(startStr, 10, 64)
	end, _ := strconv.ParseInt(endStr, 10, 64)
	if now < start || now > end {
		return nil, fmt.Errorf("不在活动时间内")
	}

	// 2. 防重复下单
	boughtKey := fmt.Sprintf("bought:act:%d:user:%d", actId, userId)
	if ok, _ := l.svcCtx.Redis.Exists(boughtKey); ok {
		return nil, fmt.Errorf("您已参与过该活动")
	}

	// 3. 原子扣减库存 (Lua)
	stockKey := fmt.Sprintf("stock:act:%d", actId)
	luaScript := `
        local stock = redis.call('DECR', KEYS[1])
        if stock < 0 then
            redis.call('INCR', KEYS[1])
            return -1
        end
        return stock
    `
	res, err := l.svcCtx.Redis.Eval(luaScript, []string{stockKey})
	if err != nil || res.(int64) < 0 {
		return nil, fmt.Errorf("库存不足")
	}

	// 4. 记录用户已购买
	l.svcCtx.Redis.Setex(boughtKey, "1", 7200)

	// 5. 发送 Kafka 消息
	orderNo := generateOrderNo()
	msg := map[string]interface{}{
		"orderNo": orderNo,
		"userId":  userId,
		"actId":   actId,
		"status":  1,
	}
	// 将 msg 序列化为 JSON 字符串
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		// 序列化失败，直接返回；此时库存和标记尚未扣减，无需回滚
		return nil, fmt.Errorf("系统错误")
	}

	// 使用 context.Background() 或传入请求的 ctx
	if err := l.svcCtx.KqPusher.Push(context.Background(), string(msgBytes)); err != nil {
		// 发送失败，回滚已扣减的库存和用户标记
		l.svcCtx.Redis.Incr(stockKey) // 恢复库存
		l.svcCtx.Redis.Del(boughtKey) // 删除用户抢购标记
		return nil, fmt.Errorf("系统繁忙，请稍后重试")
	}

	return &types.SeckillResp{
		OrderNo: orderNo,
		Message: "排队中，请稍后查询订单状态",
	}, nil
}
