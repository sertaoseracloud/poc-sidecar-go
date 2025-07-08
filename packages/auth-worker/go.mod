module auth-worker

go 1.23

require (
	github.com/rabbitmq/amqp091-go v1.10.0
	identity-adapters v0.0.0
)

replace identity-adapters => ../identity-adapters
