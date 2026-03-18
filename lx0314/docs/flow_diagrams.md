# E-commerce System Flow Diagrams

## 系统架构

```
┌─────────────────────────────────────────────────────────────────────┐
│                         Client (浏览器/APP)                          │
└────────────────────────────────┬────────────────────────────────────┘
                                 │ HTTP Request
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        bff-api/ (Gin HTTP)                          │
│                    port: 8080 (HTTP Server)                         │
└────────────────────────────────┬────────────────────────────────────┘
                                 │ gRPC
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        srv/ (gRPC Services)                         │
│                  port: 50051 (gRPC Server)                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐               │
│  │  OrderService│  │InventoryService││ PointsService│               │
│  └──────────────┘  └──────────────┘  └──────────────┘               │
└────────────────────────────────┬────────────────────────────────────┘
                                 │
        ┌────────────────────────┼────────────────────────┐
        │                        │                        │
        ▼                        ▼                        ▼
┌──────────────────┐    ┌──────────────────┐    ┌──────────────────┐
│   MySQL (GORM)   │    │  Redis (幂等性)  │    │   RabbitMQ       │
│   - User         │    │  - 消息判重      │    │   - order_queue  │
│   - Product      │    │  - 临时缓存      │    │   - inventory_q  │
│   - Order        │    │                  │    │   - points_queue │
│   - Inventory    │    │                  │    │                  │
│   - OrderItems   │    │                  │    │                  │
└──────────────────┘    └──────────────────┘    └──────────────────┘
```

---

## 1. 订单全流程 (Order Full Process)

### Sequence Diagram

```mermaid
sequenceDiagram
    participant C as Customer
    participant HTTP as bff-api (HTTP)
    participant gRPC as srv/ (gRPC)
    participant MySQL as MySQL (GORM)
    participant Redis as Redis
    participant MQ as RabbitMQ

    C->>HTTP: POST /api/order (product_id, quantity)
    HTTP->>gRPC: CreateOrderRequest
    gRPC->>MySQL: 查看库存 (SELECT for update)
    
    alt 库存不足
        MySQL-->>gRPC: 错误：库存不足
        gRPC-->>HTTP: 错误响应
        HTTP-->>C: 400 库存不足
    else 库存充足
        MySQL-->>gRPC: 库存数据
        gRPC->>Redis: 检查幂等性SETNX(order_id)
        Redis-->>gRPC: 未处理 (true)
        
        gRPC->>MySQL: 开启事务
        gRPC->>MySQL: 减少库存 (UPDATE inventory)
        gRPC->>MySQL: 创建订单 (INSERT orders)
        gRPC->>MySQL: 创建订单项 (INSERT order_items)
        
        gRPC->>Redis: 标记已处理 (SET order_id, 24h TTL)
        gRPC->>MQ: 发布消息到order_queue
        
        gRPC->>MySQL: 提交事务
        MySQL-->>gRPC: 成功
        gRPC-->>HTTP: CreateOrderResponse (order_id, total_amount)
        HTTP-->>C: 200 OK (order_id, success)
    end
```

### Data Flow

```
HTTP Layer (bff-api/handler/)
├── request/
│   └── order_create_request.go
├── response/
│   └── order_create_response.go
└── service/
    └── order_service.go (调用gRPC)

gRPC Layer (srv/handler/)
├── order_handler.go
│   ├── 1. 检查库存 (MySQL SELECT)
│   ├── 2. 幂等性检查 (Redis SETNX)
│   ├── 3. 事务处理
│   │   ├── 减少库存 (UPDATE)
│   │   ├── 创建订单 (INSERT)
│   │   └── 创建订单项 (INSERT)
│   └── 4. 发布消息 (RabbitMQ)

MySQL (srv/model/)
├── user.go
├── product.go
├── order.go (gorm.Model + 订单字段)
├── order_items.go (订单项)
└── inventory.go (库存)

Redis: mq:processed:{order_id_md5} (24h TTL)
```

---

## 2. 库存扣减流程 (Inventory Deduction Flow)

### Sequence Diagram

```mermaid
sequenceDiagram
    participant Order as Order Service
    participant Inv as Inventory Service
    participant MySQL as MySQL
    participant Redis as Redis
    participant MQ as RabbitMQ
    participant Log as Log Service

    Order->>Inv: DeductInventoryRequest (product_id, quantity)
    Inv->>MySQL: SELECT库存 FOR UPDATE
    
    alt 库存不足
        MySQL-->>Inv: 0库存
        Inv-->>Order: 错误：库存不足
    else 库存充足
        MySQL-->>Inv: 库存数据
        
        Inv->>Redis: SETNX(deduct_key, 1, 24h)
        alt 未处理
            Redis-->>Inv: true (SET成功)
            
            Inv->>MySQL: UPDATE库存 SET total_stock-quantity
            
            Inv->>MQ: 发布库存扣减消息
            MQ-->>Log: 记录日志
            
            Inv->>Redis: SET(deduct_key, done, 24h)
            Inv-->>Order: 成功
        else 已处理
            Redis-->>Inv: false (已存在)
            Inv-->>Order: 重复请求
        end
    end
```

### Data Flow

```
库存扣减流程:
1. 幂等性检查 (Redis SETNX)
   - Key: "mq:processed:inventory:{product_id}_{quantity}_{timestamp}"
   - TTL: 24小时

2. 数据库事务
   - SELECT库存 FOR UPDATE (悲观锁)
   - 检查库存≥需求数量
   - UPDATE库存减少
   - 记录库存变更日志

3. RabbitMQ消息
   - Topic: inventory_deduction_queue
   - Payload: {"product_id", " quantity", "before_stock", "after_stock"}

4. 日志记录
   - Kafka/文件日志
   - 包含: 产品ID,扣减数量,扣减前/后库存,操作时间
```

---

## 3. 积分奖励流程 (Points Awarding Flow)

### Sequence Diagram

```mermaid
sequenceDiagram
    participant MQ as RabbitMQ
    participant Points as Points Service
    participant MySQL as MySQL
    participant Redis as Redis

    MQ->>Points: PointsAwardMessage (order_id, points)
    Points->>Redis: 检查幂等性SETNX(points_key)
    
    alt 未处理
        Redis-->>Points: true
        
        Points->>MySQL: 开启事务
        Points->>MySQL: 增加用户积分 (UPDATE user SET points+=xxx)
        Points->>MySQL: 记录积分明细 (INSERT points_transactions)
        
        Points->>Redis: 标记已处理SET(points_key, 24h)
        MySQL-->>Points: 成功
        Points-->>MQ: ACK (确认消费)
    else 已处理
        Redis-->>Points: false
        Points-->>MQ: ACK (跳过重复)
    end
```

### Data Flow

```
积分奖励流程:
1. RabbitMQ消息监听
   - Queue: points_award_queue
   - Message: {"order_id", "user_id", "points_earned"}

2. 幂等性检查 (Redis SETNX)
   - Key: "mq:processed:points:{order_id}"
   - TTL: 24小时

3. 数据库事务
   - SELECT用户 (FOR UPDATE)
   - UPDATE用户积分
   - INSERT积分明细 (points_transactions表)

4. Redis幂等性标记
   - 确保订单只奖励一次积分
```

---

## 4. 完整数据流图 (Data Flow Diagram)

```mermaid
graph TD
    subgraph "Layer 1: Client"
        C[Customer]
    end
    
    subgraph "Layer 2: bff-api (Gin)"
        H1[HTTP Handler]
        H2[Request Validation]
        H3[Request Transform]
        H4[gRPC Client]
    end
    
    subgraph "Layer 3: srv/ (gRPC)"
        S1[Order Handler]
        S2[Inventory Handler]
        S3[Points Handler]
        G1[MySQL Transaction]
        G2[Redis Cache]
        G3[RabbitMQ Publisher]
    end
    
    subgraph "Layer 4: Storage"
        M[(MySQL)]
        R[(Redis)]
        Q[(RabbitMQ)]
    end
    
    C --> H1
    H1 --> H2
    H2 --> H3
    H3 --> H4
    H4 --> S1
    
    S1 -->|库存检查| S2
    S2 --> G1
    S1 --> G1
    S1 -->|消息发布| G3
    S1 -->|幂等性| G2
    
    G1 --> M
    G2 --> R
    G3 --> Q
    
    Q -->|消费消息| S3
    S3 --> G2
    S3 --> G1
    S3 -->|日志| Log[Log Service]
```

---

## 5. 消息队列流程 (RabbitMQ Flow)

```mermaid
sequenceDiagram
    participant Order as Order Service
    participant MQ as RabbitMQ
    participant Inv as Inventory Service
    participant Points as Points Service
    participant Log as Log Service

    Order->>MQ: 发布消息 (order_created)
    MQ-->>Inv: 订阅消费
    Inv->>Log: 记录库存变更日志
    
    Order->>MQ: 发布消息 (points_award)
    MQ-->>Points: 订阅消费
    Points->>Log: 记录积分奖励日志
    
    Inv->>MQ: 发布消息 (inventory_reverted)
    MQ-->>Points: 订阅消费
```

---

## 6. 错误处理流程 (Error Handling Flow)

```mermaid
flowchart TD
    A[开始] --> B{数据库操作}
    B -->|失败| C[回滚事务]
    C --> D[记录错误日志]
    D --> E[返回错误给客户端]
    
    B -->|成功| F{Redis幂等性检查}
    F -->|已存在| G[跳过重复操作]
    G --> H[ACK消息]
    
    F -->|不存在| I[设置Redis标记]
    I --> J{消息发布}
    J -->|失败| K[重试机制]
    K --> L{重试成功?}
    L -->|否| M[记录失败队列]
    L -->|是| N[继续]
    
    J -->|成功| N[完成]
```

---

## 7. 关键文件结构

```
lx0314/
├── bff-api/
│   ├── handler/
│   │   ├── request/
│   │   │   ├── order_create_request.go    # 订单创建请求
│   │   │   ├── inventory_deduct_request.go # 库存扣减请求
│   │   │   └── points_award_request.go    # 积分奖励请求
│   │   ├── response/
│   │   │   ├── order_create_response.go
│   │   │   └── ...
│   │   └── service/
│   │       └── order_service.go           # 调用gRPC服务
│   └── router/
│       └──router.go                       # HTTP路由定义
│
├── srv/
│   ├── handler/
│   │   ├── order_handler.go               # gRPC Order Service
│   │   ├── inventory_handler.go           # gRPC Inventory Service
│   │   └── points_handler.go              # gRPC Points Service
│   └── model/
│       ├── order.go                       # 订单模型
│       ├── order_items.go                 # 订单项模型
│       ├── user.go                        # 用户模型
│       ├── product.go                     # 商品模型
│       ├── inventory.go                   # 库存模型
│       └── points_transactions.go         # 积分明细模型
│
└── mq/
    └── rabbitmq.go                        # RabbitMQ封装
        ├── QueueProducer                   # 生产者
        │   ├── SendMsg()                   # 发送消息
        │   ├── SubscribeMsg()              # 订阅消息
        │   └── Close()                     # 关闭连接
        ├── RedisIdempotency                # Redis幂等性
        │   ├── CheckAndMark()              # 检查并标记
        │   └── GenerateMessageKey()        # 生成唯一键
        └── rabbitmq_test.go               # 单元测试
```

---

## 8. 注释说明

1. **幂等性实现**: 使用Redis SETNX命令确保消息消费的幂等性，key过期时间为24小时
2. **库存扣减**: 使用数据库悲观锁(SELECT FOR UPDATE)防止超卖
3. **事务处理**: MySQL事务确保数据一致性
4. **消息队列**: RabbitMQ异步处理，解耦订单、库存、积分服务
5. **日志记录**: 所有重要操作记录到日志系统
6. **重试机制**: 消息发布失败时有自动重试机制
