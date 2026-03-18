# E-commerce ER Diagram (Entity Relationship)

## Database Tables

```mermaid
erDiagram

    User ||--o{ Order : "has"
    User ||--o{ PointsTransactions : "has"
    User ||--o{ Logistics : "has"
    
    Product ||--o{ OrderItems : "includes"
    Product ||--o{ Inventory : "has"
    
    Order ||--o{ OrderItems : "contains"
    Order ||--o{ Logistics : "has"
    Order ||--o{ PointsTransactions : "triggers"
    
    Inventory }o--|| Product : "belongs to"

    User {
        bigint ID PK
        timestamp.createdAt
        timestamp.updatedAt
        timestamp.deletedAt
        string Username "用户名"
        string Email "邮箱"
        string Password "密码"
        string Phone "手机号"
        string Avatar "头像"
        tinyint Status "状态"
        tinyint Level "会员等级"
    }

    Product {
        bigint ID PK
        timestamp.createdAt
        timestamp.updatedAt
        timestamp.deletedAt
        string Name "商品名称"
        string Subtitle "副标题"
        text Desc "商品描述"
        text Images "商品图片(JSON数组)"
        decimal Price "当前价格"
        int Stock "库存数量"
        tinyint Status "状态"
    }

    Order {
        bigint ID PK
        timestamp.createdAt
        timestamp.updatedAt
        timestamp.deletedAt
        string OrderNo "订单号"
        bigint UserID FK "用户ID"
        decimal Total "订单总金额"
        decimal PayAmount "订单支付金额"
        tinyint Status "订单状态"
        tinyint PayType "支付方式"
        tinyint PayUrl "支付URL"
    }

    OrderItems {
        bigint ID PK
        timestamp.createdAt
        timestamp.updatedAt
        timestamp.deletedAt
        bigint OrderId FK "订单ID"
        bigint ProductId FK "商品ID"
        string ProductName "商品名称"
        string ProductImage "商品图片"
        decimal Price "商品单价"
        int Quantity "购买数量"
        decimal Total "商品总金额"
    }

    Inventory {
        bigint ID PK
        timestamp.createdAt
        timestamp.updatedAt
        timestamp.deletedAt
        bigint ProductId FK "商品ID"
        int TotalStock "总库存"
        int FrozenStock "冻结库存"
        int AvailableStock "可用库存"
        bigint WarehouseId "仓库ID"
        string Location "库位"
    }

    Logistics {
        bigint ID PK
        timestamp.createdAt
        timestamp.updatedAt
        timestamp.deletedAt
        bigint OrderId FK "订单ID"
        string OrderNo "订单号"
        bigint UserId FK "用户ID"
        string ExpressCompany "快递公司"
        string ExpressNo "快递单号"
        string ReceiverName "收货人姓名"
        string ReceiverPhone "收货人电话"
        string Province "省"
        string City "市"
        string Area "区/县"
        string Address "详细地址"
        tinyint Status "物流状态"
        bigint LogisticsId "物流ID"
    }

    PointsTransactions {
        bigint ID PK
        timestamp.createdAt
        timestamp.updatedAt
        timestamp.deletedAt
        bigint UserId FK "用户ID"
        bigint Points "积分变动数量"
        tinyint Type "类型"
        tinyint RelatedType "关联类型"
        bigint RelatedId "关联ID"
        bigint Balance "变动后总积分"
        string Description "描述"
    }
```

## ER Diagram Description

### Tables Overview

| Table | Description | Key Fields |
|-------|-------------|------------|
| **User** | 用户表 | ID, Username, Email, Phone, Level |
| **Product** | 商品表 | ID, Name, Price, Stock, Status |
| **Order** | 订单表 | ID, OrderNo, UserID, Total, Status |
| **OrderItems** | 订单明细表 | ID, OrderId, ProductId, Quantity |
| **Inventory** | 库存表 | ID, ProductId, TotalStock, FrozenStock |
| **Logistics** | 物流表 | ID, OrderId, ExpressCompany, Status |
| **PointsTransactions** | 积分明细表 | ID, UserId, Points, Type, Balance |

### Relationships

1. **User → Order** (1:N)
   - One user can place many orders
   - Foreign Key: `Order.UserID`

2. **User → PointsTransactions** (1:N)
   - One user can have many point transactions
   - Foreign Key: `PointsTransactions.UserId`

3. **User → Logistics** (1:N)
   - One user can have many logistics records
   - Foreign Key: `Logistics.UserId`

4. **Product → Order Items** (1:N)
   - One product can appear in many order items
   - Foreign Key: `OrderItems.ProductId`

5. **Product → Inventory** (1:1)
   - One product has one inventory record
   - Foreign Key: `Inventory.ProductId`

6. **Order → Order Items** (1:N)
   - One order can contain many order items
   - Foreign Key: `OrderItems.OrderId`

7. **Order → Logistics** (1:1)
   - One order has one logistics record
   - Foreign Key: `Logistics.OrderId`

8. **Order → PointsTransactions** (1:N)
   - One order can trigger multiple point transactions
   - Foreign Key: `PointsTransactions.RelatedId` (when RelatedType=1)

### Foreign Keys Summary

| Table | FK Column | Ref Table | Description |
|-------|-----------|-----------|-------------|
| Order | UserID | User | User who placed the order |
| OrderItems | OrderId | Order | Order containing this item |
| OrderItems | ProductId | Product | Product in this order item |
| Inventory | ProductId | Product | Product inventory |
| Logistics | OrderId | Order | Order logistics |
| Logistics | UserId | User | User who placed the order |
| PointsTransactions | UserId | User | User who earned/spent points |
| PointsTransactions | RelatedId | Order | Related order (when Type=1) |

### Index Recommendations

```sql
-- User
CREATE INDEX idx_user_email ON User(Email);
CREATE INDEX idx_user_phone ON User(Phone);

-- Product
CREATE INDEX idx_product_name ON Product(Name);
CREATE INDEX idx_product_status ON Product(Status);

-- Order
CREATE INDEX idx_order_user_id ON Order(UserID);
CREATE INDEX idx_order_no ON Order(OrderNo);
CREATE INDEX idx_order_status ON Order(Status);

-- OrderItems
CREATE INDEX idx_orderitems_order_id ON OrderItems(OrderId);
CREATE INDEX idx_orderitems_product_id ON OrderItems(ProductId);

-- Inventory
CREATE INDEX idx_inventory_product_id ON Inventory(ProductId);

-- Logistics
CREATE INDEX idx_logistics_order_id ON Logistics(OrderId);
CREATE INDEX idx_logistics_user_id ON Logistics(UserId);
CREATE INDEX idx_logistics_status ON Logistics(Status);

-- PointsTransactions
CREATE INDEX idx_points_user_id ON PointsTransactions(UserId);
CREATE INDEX idx_points_related ON PointsTransactions(RelatedType, RelatedId);
```

### Stock Flow Logic

```
初始库存: TotalStock = 100
下单冻结: FrozenStock = 20, AvailableStock = 80
订单完成: FrozenStock = 0, TotalStock = 80
订单取消: FrozenStock = 0, AvailableStock = 100
```

### Points Awarding Logic

```
订单支付完成后:
- Points = order.Total * 0.05 (5% return)
- Balance = User.Points + Points
- Type = 1 (获得)
- RelatedType = 1 (订单)
```
