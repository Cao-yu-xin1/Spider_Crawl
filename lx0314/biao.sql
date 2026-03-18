-- E-commerce Database Schema
-- Database: ecommerce
-- MySQL 8.0+

-- 1. Users (注册会员/会员)
CREATE TABLE users (
                       id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '用户ID',
                       username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
                       email VARCHAR(100) UNIQUE COMMENT '邮箱',
                       password VARCHAR(255) NOT NULL COMMENT '密码(加密存储)',
                       phone VARCHAR(20) UNIQUE COMMENT '手机号',
                       avatar VARCHAR(255) COMMENT '头像',
                       status TINYINT DEFAULT 1 COMMENT '状态: 0-禁用, 1-启用, 2-注销',
                       level INT DEFAULT 0 COMMENT '会员等级',
                       points BIGINT DEFAULT 0 COMMENT '当前积分',
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                       INDEX idx_username(username),
                       INDEX idx_phone(phone),
                       INDEX idx_email(email),
                       INDEX idx_status(status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 2. Product Categories (商品分类)
CREATE TABLE categories (
                            id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '分类ID',
                            name VARCHAR(100) NOT NULL COMMENT '分类名称',
                            parent_id BIGINT DEFAULT 0 COMMENT '父分类ID(0-顶级分类)',
                            level INT DEFAULT 1 COMMENT '分类层级',
                            sort_order INT DEFAULT 0 COMMENT '排序',
                            status TINYINT DEFAULT 1 COMMENT '状态: 0-禁用, 1-启用',
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                            INDEX idx_parent_id(parent_id),
                            INDEX idx_status(status),
                            INDEX idx_level(level)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品分类表';

-- 3. Products (电商商品)
CREATE TABLE products (
                          id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '商品ID',
                          name VARCHAR(200) NOT NULL COMMENT '商品名称',
                          subtitle VARCHAR(200) COMMENT '副标题',
                          description TEXT COMMENT '商品描述',
                          main_image VARCHAR(500) COMMENT '主图',
                          images TEXT COMMENT '商品图片(JSON数组)',
                          price DECIMAL(10,2) NOT NULL COMMENT '当前价格',
                          original_price DECIMAL(10,2) COMMENT '原价',
                          cost_price DECIMAL(10,2) COMMENT '成本价',
                          stock INT DEFAULT 0 COMMENT '库存数量',
                          sold_count INT DEFAULT 0 COMMENT '销量',
                          category_id BIGINT NOT NULL COMMENT '分类ID',
                          brand_id BIGINT COMMENT '品牌ID',
                          unit VARCHAR(20) DEFAULT '件' COMMENT '单位',
                          specifications JSON COMMENT '规格(JSON)',
                          status TINYINT DEFAULT 1 COMMENT '状态: 0-下架, 1-上架, 2-删除',
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                          FOREIGN KEY (category_id) REFERENCES categories(id),
                          INDEX idx_name(name),
                          INDEX idx_category_id(category_id),
                          INDEX idx_status(status),
                          INDEX idx_created_at(created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品表';

-- 4. Inventory (库存)
CREATE TABLE inventory (
                           id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '库存ID',
                           product_id BIGINT NOT NULL UNIQUE COMMENT '商品ID',
                           total_stock INT DEFAULT 0 COMMENT '总库存',
                           frozen_stock INT DEFAULT 0 COMMENT '冻结库存(下单占用)',
                           available_stock INT DEFAULT 0 COMMENT '可用库存',
                           warehouse_id BIGINT COMMENT '仓库ID',
                           location VARCHAR(100) COMMENT '库位',
                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                           FOREIGN KEY (product_id) REFERENCES products(id),
                           INDEX idx_warehouse_id(warehouse_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='库存表';

-- 5. Points Transactions (积分流水)
CREATE TABLE points_transactions (
                                     id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '积分流水ID',
                                     user_id BIGINT NOT NULL COMMENT '用户ID',
                                     points INT NOT NULL COMMENT '积分变动数量(正数增加,负数减少)',
                                     type TINYINT NOT NULL COMMENT '类型: 1-获得, 2-消费, 3-过期, 4-管理员调整',
                                     related_type TINYINT COMMENT '关联类型: 1-订单, 2-评价, 3-签到, 4-活动',
                                     related_id BIGINT COMMENT '关联ID(订单ID/评价ID等)',
                                     balance INT NOT NULL COMMENT '变动后总积分',
                                     description VARCHAR(200) COMMENT '描述',
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                     FOREIGN KEY (user_id) REFERENCES users(id),
                                     INDEX idx_user_id(user_id),
                                     INDEX idx_created_at(created_at),
                                     INDEX idx_type(type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='积分流水表';

-- 6. Orders (订单表)
CREATE TABLE orders (
                        id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '订单ID',
                        order_no VARCHAR(30) NOT NULL UNIQUE COMMENT '订单号',
                        user_id BIGINT NOT NULL COMMENT '用户ID',
                        total_amount DECIMAL(10,2) NOT NULL COMMENT '订单总金额',
                        pay_amount DECIMAL(10,2) NOT NULL COMMENT '支付金额',
                        points_deduct INT DEFAULT 0 COMMENT '抵扣积分',
                        discount_amount DECIMAL(10,2) DEFAULT 0.00 COMMENT '优惠金额',
                        freight_amount DECIMAL(10,2) DEFAULT 0.00 COMMENT '运费',
                        status TINYINT NOT NULL DEFAULT 0 COMMENT '状态: 0-待付款, 1-待发货, 2-待收货, 3-已完成, 4-已取消, 5-退款中, 6-已退款',
                        payment_method TINYINT COMMENT '支付方式: 1-微信, 2-支付宝, 3-余额, 4-混合支付',
                        payment_time TIMESTAMP NULL COMMENT '支付时间',
                        delivery_time TIMESTAMP NULL COMMENT '发货时间',
                        finish_time TIMESTAMP NULL COMMENT '完成时间',
                        receiver_name VARCHAR(50) NOT NULL COMMENT '收货人姓名',
                        receiver_phone VARCHAR(20) NOT NULL COMMENT '收货人电话',
                        receiver_address VARCHAR(500) NOT NULL COMMENT '收货地址',
                        receiver_zip VARCHAR(10) COMMENT '邮编',
                        invoice_title VARCHAR(200) COMMENT '发票抬头',
                        invoice_content VARCHAR(500) COMMENT '发票内容',
                        notes VARCHAR(500) COMMENT '订单备注',
                        cancel_reason VARCHAR(200) COMMENT '取消原因',
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                        FOREIGN KEY (user_id) REFERENCES users(id),
                        INDEX idx_order_no(order_no),
                        INDEX idx_user_id(user_id),
                        INDEX idx_status(status),
                        INDEX idx_created_at(created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单表';

-- 7. Order Items (订单明细)
CREATE TABLE order_items (
                             id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '订单项ID',
                             order_id BIGINT NOT NULL COMMENT '订单ID',
                             product_id BIGINT NOT NULL COMMENT '商品ID',
                             product_name VARCHAR(200) NOT NULL COMMENT '商品名称',
                             product_image VARCHAR(500) COMMENT '商品图片',
                             product_sku VARCHAR(100) COMMENT '商品规格',
                             price DECIMAL(10,2) NOT NULL COMMENT '商品单价',
                             original_price DECIMAL(10,2) COMMENT '商品原价',
                             quantity INT NOT NULL COMMENT '购买数量',
                             total_amount DECIMAL(10,2) NOT NULL COMMENT '商品总金额',
                             pointsEarned INT DEFAULT 0 COMMENT '获得积分',
                             status TINYINT DEFAULT 1 COMMENT '状态: 1-正常, 2-已退款',
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                             updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                             FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
                             FOREIGN KEY (product_id) REFERENCES products(id),
                             INDEX idx_order_id(order_id),
                             INDEX idx_product_id(product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单明细表';

-- 8. Logistics/Shipping (物流配送)
CREATE TABLE logistics (
                           id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '物流ID',
                           order_id BIGINT NOT NULL UNIQUE COMMENT '订单ID',
                           order_no VARCHAR(30) NOT NULL COMMENT '订单号',
                           user_id BIGINT NOT NULL COMMENT '用户ID',
                           express_company VARCHAR(100) COMMENT '快递公司',
                           express_no VARCHAR(100) COMMENT '快递单号',
                           receiver_name VARCHAR(50) COMMENT '收货人姓名',
                           receiver_phone VARCHAR(20) COMMENT '收货人电话',
                           province VARCHAR(50) COMMENT '省',
                           city VARCHAR(50) COMMENT '市',
                           area VARCHAR(50) COMMENT '区/县',
                           address VARCHAR(500) COMMENT '详细地址',
                           status TINYINT DEFAULT 0 COMMENT '状态: 0-待发货, 1-已发货, 2-已签收, 3-已退回',
                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                           FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
                           FOREIGN KEY (user_id) REFERENCES users(id),
                           INDEX idx_order_id(order_id),
                           INDEX idx_user_id(user_id),
                           INDEX idx_status(status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='物流表';

-- 9. Logistics Track (物流轨迹)
CREATE TABLE logistics_track (
                                 id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '物流轨迹ID',
                                 logistics_id BIGINT NOT NULL COMMENT '物流ID',
                                 status VARCHAR(100) COMMENT '状态描述',
                                 location VARCHAR(200) COMMENT '位置',
                                 time TIMESTAMP COMMENT '时间',
                                 note VARCHAR(500) COMMENT '备注',
                                 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                 FOREIGN KEY (logistics_id) REFERENCES logistics(id) ON DELETE CASCADE,
                                 INDEX idx_logistics_id(logistics_id),
                                 INDEX idx_time(time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='物流轨迹表';

-- 10. Inventory Log (库存变动日志)
CREATE TABLE inventory_log (
                               id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '日志ID',
                               product_id BIGINT NOT NULL COMMENT '商品ID',
                               order_id BIGINT COMMENT '订单ID',
                               order_item_id BIGINT COMMENT '订单项ID',
                               change_type TINYINT NOT NULL COMMENT '变动类型: 1-入库, 2-出库, 3-冻结, 4-解冻, 5-盘盈, 6-盘亏',
                               original_stock INT NOT NULL COMMENT '变动前库存',
                               change_stock INT NOT NULL COMMENT '变动数量(正数增加,负数减少)',
                               final_stock INT NOT NULL COMMENT '变动后库存',
                               related_no VARCHAR(50) COMMENT '关联单号(订单号/出库单号)',
                               operator VARCHAR(50) COMMENT '操作人',
                               remark VARCHAR(500) COMMENT '备注',
                               created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                               FOREIGN KEY (product_id) REFERENCES products(id),
                               INDEX idx_product_id(product_id),
                               INDEX idx_order_id(order_id),
                               INDEX idx_created_at(created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='库存变动日志表';

-- 11. User Address (收货地址)
CREATE TABLE user_addresses (
                                id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '地址ID',
                                user_id BIGINT NOT NULL COMMENT '用户ID',
                                receiver_name VARCHAR(50) NOT NULL COMMENT '收货人姓名',
                                receiver_phone VARCHAR(20) NOT NULL COMMENT '收货人电话',
                                province VARCHAR(50) NOT NULL COMMENT '省',
                                city VARCHAR(50) NOT NULL COMMENT '市',
                                area VARCHAR(50) NOT NULL COMMENT '区/县',
                                address VARCHAR(500) NOT NULL COMMENT '详细地址',
                                is_default TINYINT DEFAULT 0 COMMENT '是否默认: 0-否, 1-是',
                                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                                INDEX idx_user_id(user_id),
                                INDEX idx_is_default(is_default)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户收货地址表';

-- 12. Cart (购物车)
CREATE TABLE carts (
                       id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '购物车ID',
                       user_id BIGINT NOT NULL COMMENT '用户ID',
                       product_id BIGINT NOT NULL COMMENT '商品ID',
                       product_name VARCHAR(200) NOT NULL COMMENT '商品名称',
                       product_image VARCHAR(500) COMMENT '商品图片',
                       price DECIMAL(10,2) NOT NULL COMMENT '商品价格',
                       quantity INT NOT NULL DEFAULT 1 COMMENT '数量',
                       added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                       FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                       FOREIGN KEY (product_id) REFERENCES products(id),
                       INDEX idx_user_id(user_id),
                       INDEX idx_product_id(product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='购物车表';

-- 13. Reviews (商品评价)
CREATE TABLE reviews (
                         id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '评价ID',
                         order_item_id BIGINT NOT NULL UNIQUE COMMENT '订单项ID',
                         user_id BIGINT NOT NULL COMMENT '用户ID',
                         product_id BIGINT NOT NULL COMMENT '商品ID',
                         order_id BIGINT NOT NULL COMMENT '订单ID',
                         score TINYINT NOT NULL COMMENT '评分: 1-5',
                         content VARCHAR(1000) COMMENT '评价内容',
                         images JSON COMMENT '评价图片(JSON数组)',
                         is_anonymous TINYINT DEFAULT 0 COMMENT '是否匿名: 0-否, 1-是',
                         reply_content VARCHAR(500) COMMENT '商家回复',
                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                         updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                         FOREIGN KEY (order_item_id) REFERENCES order_items(id),
                         FOREIGN KEY (user_id) REFERENCES users(id),
                         FOREIGN KEY (product_id) REFERENCES products(id),
                         INDEX idx_user_id(user_id),
                         INDEX idx_product_id(product_id),
                         INDEX idx_created_at(created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品评价表';

-- 14. User Points Summary (用户积分统计)
CREATE TABLE user_points_summary (
                                     id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '统计ID',
                                     user_id BIGINT NOT NULL UNIQUE COMMENT '用户ID',
                                     total_earned BIGINT DEFAULT 0 COMMENT '累计获得积分',
                                     total_consumed BIGINT DEFAULT 0 COMMENT '累计消费积分',
                                     total_expired BIGINT DEFAULT 0 COMMENT '累计过期积分',
                                     current_balance BIGINT DEFAULT 0 COMMENT '当前余额',
                                     last_earn_at TIMESTAMP NULL COMMENT '最后获得时间',
                                     last_consume_at TIMESTAMP NULL COMMENT '最后消费时间',
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                     FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户积分统计表';

-- 15. Payment (支付记录)
CREATE TABLE payments (
                          id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '支付ID',
                          order_no VARCHAR(30) NOT NULL COMMENT '订单号',
                          user_id BIGINT NOT NULL COMMENT '用户ID',
                          amount DECIMAL(10,2) NOT NULL COMMENT '支付金额',
                          pay_type TINYINT NOT NULL COMMENT '支付类型: 1-微信, 2-支付宝, 3-余额, 4-混合',
                          transaction_no VARCHAR(100) COMMENT '第三方交易号',
                          status TINYINT DEFAULT 0 COMMENT '状态: 0-待支付, 1-已支付, 2-支付失败, 3-已退款',
                          paid_at TIMESTAMP NULL COMMENT '支付完成时间',
                          refund_amount DECIMAL(10,2) DEFAULT 0.00 COMMENT '退款金额',
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                          FOREIGN KEY (user_id) REFERENCES users(id),
                          INDEX idx_order_no(order_no),
                          INDEX idx_transaction_no(transaction_no),
                          INDEX idx_status(status),
                          INDEX idx_created_at(created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='支付记录表';

-- 16. Brand (品牌表)
CREATE TABLE brands (
                        id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '品牌ID',
                        name VARCHAR(100) NOT NULL COMMENT '品牌名称',
                        logo VARCHAR(255) COMMENT '品牌logo',
                        description VARCHAR(500) COMMENT '品牌描述',
                        sort_order INT DEFAULT 0 COMMENT '排序',
                        status TINYINT DEFAULT 1 COMMENT '状态: 0-禁用, 1-启用',
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='品牌表';