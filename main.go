package main

import (
    "fmt"
    "os"

    "github.com/streadway/amqp"
)

func main() {
    fmt.Println("Starting server")

    // Get the connection string from the environment variable
    url := os.Getenv("AMQP_URL")

    //If it doesn't exist, use the default connection string.

    if url == "" {
        //Don't do this in production, this is for testing purposes only.
        url = "amqp://guest:guest@localhost:5672"
    }

    // Connect to the rabbitMQ instance
    connection, err := amqp.Dial(url)

    if err != nil {
        panic("could not establish connection with RabbitMQ:" + err.Error())
    }

    // Create a channel from the connection. We'll use channels to access the data in the queue rather than the connection itself.
    channel, err := connection.Channel()

    if err != nil {
        panic("could not open RabbitMQ channel:" + err.Error())
    }

    // We create an exchange that will bind to the queue to send and receive messages
    err = channel.ExchangeDeclare("my-topic", "topic", true, false, false, false, nil)

    if err != nil {
        panic(err)
    }

    // We create a message to be sent to the queue. 
    // It has to be an instance of the aqmp publishing struct
    message := amqp.Publishing{
        Body: []byte("Hello World"),
    }

    // We publish the message to the exahange we created earlier
    err = channel.Publish("my-topic", "bitcoin-rial", false, false, message)

    if err != nil {
        panic("error publishing a message to the queue:" + err.Error())
    }

    // We create a queue named Test
    _, err = channel.QueueDeclare("tubeBitCRial", true, false, false, false, nil)

    if err != nil {
        panic("error declaring the queue: " + err.Error())
    }

    // We bind the queue to the exchange to send and receive data from the queue
    err = channel.QueueBind("tubeBitCRial", "#", "my-topic", false, nil)

    if err != nil {
        panic("error binding to the queue: " + err.Error())
    }

    // We consume data in the queue named test using the channel we created in go.
    msgs, err := channel.Consume("tubeBitCRial", "", false, false, false, false, nil)

    if err != nil {
        panic("error consuming the queue: " + err.Error())
    }

    // We loop through the messages in the queue and print them to the console.
    // The msgs will be a go channel, not an amqp channel
    for msg := range msgs {
    //print the message to the console
        fmt.Println("message received: " + string(msg.Body))
    // Acknowledge that we have received the message so it can be removed from the queue
        msg.Ack(false)
    }

    // We close the connection after the operation has completed.
    defer connection.Close()
}
