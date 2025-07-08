package main

import (
    "encoding/json"
    "log"
    "os"

    amqp "github.com/rabbitmq/amqp091-go"
    adapters "identity-adapters"
)

const authQueue = "auth_requests"

func main() {
    rabbitURL := os.Getenv("RABBIT_URL")
    if rabbitURL == "" {
        rabbitURL = "amqp://rabbitmq"
    }
    conn, err := amqp.Dial(rabbitURL)
    if err != nil {
        log.Fatalf("amqp dial: %v", err)
    }
    ch, err := conn.Channel()
    if err != nil {
        log.Fatalf("channel: %v", err)
    }
    _, err = ch.QueueDeclare(authQueue, false, false, false, false, nil)
    if err != nil {
        log.Fatalf("declare: %v", err)
    }
    msgs, err := ch.Consume(authQueue, "", false, false, false, false, nil)
    if err != nil {
        log.Fatalf("consume: %v", err)
    }
    log.Println("Auth worker waiting for messages")
    for msg := range msgs {
        provider := string(msg.Body)
        var res adapters.Credentials
        if provider == "aws" {
            res, _ = adapters.GetAwsCredentials()
        } else {
            res, _ = adapters.GetAzureToken()
        }
        body, _ := json.Marshal(res)
        ch.Publish("", msg.ReplyTo, false, false, amqp.Publishing{
            CorrelationId: msg.CorrelationId,
            Body:          body,
        })
        ch.Ack(msg.DeliveryTag, false)
    }
}
