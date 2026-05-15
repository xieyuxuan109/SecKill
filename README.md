# 秒杀服务

基于 Go-Zero 框架实现的高性能秒杀微服务系统。

## 技术栈

- Go 1.25+
- Go-Zero 框架
- MySQL 数据库
- Redis 缓存
- Kafka 消息队列
- ETCD 服务发现

## 功能特性

- 秒杀活动管理：创建、查询活动状态
- 库存管理：Redis 原子扣减，高并发支持
- 订单管理：异步创建订单，支持订单状态查询
- 消息队列：Kafka 异步解耦秒杀请求
- 防重复下单：Redis 标记用户已购买状态
- 活动时间校验：精确控制秒杀开始和结束时间（毫秒级精度）

## 快速开始

### 环境要求

- Docker & Docker Compose
- Go 1.25+

### 启动依赖服务

```bash
cd SecKill
docker compose up -d
```

### 启动服务

#### 1. 启动 Order RPC 服务

```bash
cd order_rpc
go run order.go -f etc/order.yaml
```

#### 2. 启动 Consumer 服务

```bash
cd consumer
go run main.go
```

#### 3. 启动 API 服务

```bash
cd api
go run api.go -f etc/api.yaml
```

## API 接口

### 创建秒杀活动

**POST** `/api/act/create`

请求参数：

| 参数名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| name | string | 是 | 活动名称 |
| stock | int | 是 | 活动库存 |
| startAt | int64 | 是 | 开始时间戳（秒） |
| endAt | int64 | 是 | 结束时间戳（秒） |

请求示例：

```json
{
    "name": "限时秒杀",
    "stock": 100,
    "startAt": 1779000000,
    "endAt": 1779086400
}
```

说明：`startAt` 和 `endAt` 为秒级时间戳。**注意：请使用未来的时间戳**，否则会提示"活动已结束"。可使用以下命令获取时间戳：

```bash
# 获取当前时间戳
date +%s

# 获取24小时后的时间戳（推荐用于 endAt）
echo $(( $(date +%s) + 86400 ))

# 获取1小时后的时间戳（推荐用于 startAt）
echo $(( $(date +%s) + 3600 ))
```

响应示例：

```json
{
    "activityId": 1
}
```

> **注意**：创建活动后会返回 `activityId`，用于后续参与秒杀。

### 参与秒杀

**POST** `/api/seckill`

请求参数：

| 参数名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| activityId | int64 | 是 | 活动ID |

请求示例：

```json
{
    "activityId": 1
}
```

响应示例：

```json
{
    "orderNo": "1747200000123456",
    "message": "排队中，请稍后查询订单状态"
}
```

### 查询订单状态

**GET** `/api/order/status?orderNo=xxx`

请求参数：

| 参数名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| orderNo | string | 是 | 订单号 |

响应示例：

```json
{
    "orderNo": "1747200000123456",
    "status": 1
}
```

状态码说明：
- 0: 未支付
- 1: 已支付
- 2: 已取消
- 3: 已完成

### 查询活动状态

**GET** `/api/act/status?activityId=xxx`

请求参数：

| 参数名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| activityId | int64 | 是 | 活动ID |

响应示例：

```json
{
    "activityId": 1,
    "name": "限时秒杀",
    "status": 1,
    "stock": 95,
    "startAt": 1747200000,
    "endAt": 1747286400
}
```

状态码说明：
- 0: 未开始
- 1: 进行中
- 2: 已结束

### 健康检查

**GET** `/health`

响应示例：

```json
{
    "status": "ok"
}
```

## 项目结构

```
SecKill/
├── api/                    # REST API 服务
│   ├── etc/               # 配置文件
│   ├── internal/          # 内部代码
│   │   ├── config/        # 配置定义
│   │   ├── handler/       # HTTP 处理器
│   │   ├── logic/         # 业务逻辑
│   │   ├── svc/           # 服务上下文
│   │   └── types/         # 请求响应类型
│   ├── model/             # 数据模型
│   ├── api.api            # API 定义文件
│   ├── api.go             # 入口文件
│   └── go.mod             # 依赖管理
├── consumer/              # Kafka 消费者
│   ├── etc/               # 配置文件
│   ├── orderclient/       # RPC 客户端
│   └── main.go            # 入口文件
├── order_rpc/             # 订单 RPC 服务
│   ├── etc/               # 配置文件
│   ├── internal/          # 内部代码
│   ├── order/             # RPC 定义
│   ├── orderclient/       # RPC 客户端
│   ├── order.proto        # Proto 定义
│   └── order.go           # 入口文件
├── test/                   # 测试脚本
│   └── load_test.py       # 压测脚本
├── sql/                   # 数据库初始化脚本
├── dockercompose.yaml     # Docker 配置
├── go.mod                 # 根模块
└── go.work                # Go 工作区
```

## 配置说明

### API 配置 (api/etc/api.yaml)

```yaml
Name: api
Host: 0.0.0.0
Port: 8888

Redis:
  Host: localhost:6379
  Pass: ""

Kafka:
  Name: order-pusher
  Group: order-group
  Brokers:
    - localhost:9092
  Topic: seckill_orders

OrderRpc:
  Etcd:
    Hosts:
      - localhost:2379
    Key: order.rpc

Mysql:
  DataSource: root:123456@tcp(localhost:3306)/seckill?charset=utf8mb4&parseTime=true
```

### Consumer 配置 (consumer/etc/consumer.yaml)

```yaml
Kafka:
  Name: order-consumer
  Group: order-consumer-group
  Brokers:
    - localhost:9092
  Topic: seckill_orders

OrderRpc:
  Etcd:
    Hosts:
      - localhost:2379
    Key: order.rpc
```

### Order RPC 配置 (order_rpc/etc/order.yaml)

```yaml
Name: order.rpc
ListenOn: 0.0.0.0:8081
Etcd:
  Hosts:
    - localhost:2379
  Key: order.rpc
Mysql:
  DataSource: root:123456@tcp(localhost:3306)/seckill?charset=utf8mb4&parseTime=true
```

## 数据库结构

### seckill_activity 表

| 字段名 | 类型 | 说明 |
| :--- | :--- | :--- |
| id | bigint | 主键，自增 |
| name | varchar(100) | 活动名称 |
| stock | int | 库存数量 |
| start_at | bigint | 开始时间戳（秒） |
| end_at | bigint | 结束时间戳（秒） |

### seckill_order 表

| 字段名 | 类型 | 说明 |
| :--- | :--- | :--- |
| id | bigint | 主键，自增 |
| order_no | varchar(64) | 订单号（唯一） |
| user_id | bigint | 用户ID |
| act_id | bigint | 活动ID |
| status | tinyint | 订单状态 |
| create_time | datetime | 创建时间 |

## 秒杀流程

1. 用户请求 `/api/seckill` 参与秒杀
2. 系统校验活动是否存在、时间是否有效（毫秒级精度）
3. 检查用户是否已参与过该活动（防重复）
4. Redis 原子扣减库存（Lua 脚本）
5. 记录用户已购买标记
6. 发送 Kafka 消息，异步创建订单
7. Consumer 消费消息，调用 Order RPC 创建订单
8. 用户通过 `/api/order/status` 查询订单状态

## 测试

### 单元测试

运行所有单元测试：

```bash
cd api
go test -v ./internal/logic/...
```

时间校验测试用例：

| 测试场景 | 输入 | 预期结果 |
| :--- | :--- | :--- |
| 活动进行中 | now=10000, start=0, end=20000 | 通过 |
| 活动尚未开始 | now=0, start=10000, end=20000 | 返回错误"活动尚未开始" |
| 活动已结束 | now=30000, start=0, end=20000 | 返回错误"活动已结束" |
| 边界测试-刚好开始 | now=10000, start=10000, end=20000 | 通过 |
| 边界测试-刚好结束 | now=20000, start=0, end=20000 | 通过 |

### 压力测试

使用压测脚本进行高并发压测：

```bash
cd SecKill
python test/load_test.py
```

压测配置参数：

| 参数 | 默认值 | 说明 |
| :--- | :--- | :--- |
| BASE_URL | http://localhost:8888 | API 服务地址 |
| ACTIVITY_ID | 1 | 活动ID |
| THREAD_COUNT | 50 | 并发线程数 |
| REQUESTS_PER_THREAD | 100 | 每线程请求数 |

压测结果指标：
- 总请求数
- 成功数
- 失败数
- 成功率
- 总耗时
- QPS（每秒请求数）

## 注意事项

1. 确保 Docker 服务正常运行
2. 修改配置文件中的 IP 地址为实际服务器地址
3. Kafka 需提前创建 `seckill_orders` 主题
4. Redis 建议开启持久化配置

## License

MIT License
