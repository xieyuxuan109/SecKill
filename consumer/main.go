package main

import (
	"context"
	"encoding/json"
	"log"

	"order_rpc/orderclient"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
)

type OrderMsg struct {
	OrderNo string `json:"orderNo"`
	UserId  int64  `json:"userId"`
	ActId   int64  `json:"actId"`
	Status  int32  `json:"status"`
}

type Config struct {
	Kafka    kq.KqConf
	OrderRpc zrpc.RpcClientConf
}

func main() {
	var c Config
	conf.MustLoad("etc/consumer.yaml", &c)

	orderRpc := orderclient.NewOrder(zrpc.MustNewClient(c.OrderRpc))

	// 创建消息处理器
	handler := kq.WithHandle(func(ctx context.Context, key, value string) error {
		var msg OrderMsg
		if err := json.Unmarshal([]byte(value), &msg); err != nil {
			log.Printf("unmarshal error: %v", err)
			return err
		}
		_, err := orderRpc.CreateOrder(ctx, &orderclient.CreateOrderReq{
			OrderNo: msg.OrderNo,
			UserId:  msg.UserId,
			ActId:   msg.ActId,
			Status:  msg.Status,
		})
		if err != nil {
			log.Printf("create order error: %v", err)
			return err
		}
		log.Printf("order created: %s", msg.OrderNo)
		return nil
	})

	// 传入 handler
	consumer := kq.MustNewQueue(c.Kafka, handler)
	consumer.Start()
}
