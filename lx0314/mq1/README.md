# Add RabbitMQ and Redis dependencies to go.mod
github.com/streadway/amqp v1.1.0
github.com/go-redis/redis/v8 v8.11.5

# mq package implements queue producer and consumer with Redis idempotency
# Key features:
# - sendMsg(topic, msg string) for message publishing
# - subscribeMsg(topic, handler) for message consumption
# - Redis SETNX for idempotency check
# - Inventory deduction message flow
# - Points awarding message flow