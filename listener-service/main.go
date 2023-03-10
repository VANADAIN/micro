package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"
	"listener/event"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// try to connect to rabbitmw
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	defer rabbitConn.Close()

	// start listening fro messages
	log.Println("Listening for and consuming RabbitMQ messages...")

	// create consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	// watch teh queue and consume events
	err = consumer.Listen([]string{
		"log.INFO",
		"log.WARNING",
		"log.ERROR",
	})

	if err != nil {
		log.Println(err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continur until rabbit
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMq is not yet ready...")
			counts++
		} else {
			connection = c
			log.Println("Connected to RabbitMQ")
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
