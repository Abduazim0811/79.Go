package rabbitmq

import (
	"Task/internal/models"
	"Task/internal/mongodb"
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

func StartConsumer() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"Tasks",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var task models.Tasks
			if err := json.Unmarshal(d.Body, &task); err != nil {
				log.Printf("Error decoding JSON: %s", err)
				continue
			}
			mongoDb, err := mongodb.NewTask()
			if err != nil {
				log.Fatal(err)
			}
			if err = mongoDb.StoreNewTask(task); err != nil {
				log.Printf("Error saving task to MongoDB: %s", err)
                continue
			}
			log.Printf("Received a task: %s", task.Id.Hex())
		}
	}()

	log.Printf("Waiting for messages. To exit press CTRL+C")
	<-forever
}
