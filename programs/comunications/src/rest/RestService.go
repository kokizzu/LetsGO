package rest

import (
	"encoding/json"
	. "kafka"
	"log"
	"net/http"
	"strings"
	"time"
)

/*
HttpServer method of our Go application to run the server in a port and configure the router to redirect endpoints
invocations to the specific handlers.

In Go in order to control and route request into handle, we use [http] package and function [HandleFunc]
where we pass the endpoint to bind, and the function handle to call once we receive the request.

*/
func HttpServer() {
	println("Running Rest server on port 4000")
	mux := http.NewServeMux()
	mux.HandleFunc("/communication/restKafkaGRPC/", processRequest)
	log.Fatal(http.ListenAndServe(":4000", mux))
}

/*
Function to receive  a rest call, publish to Kafka the message received and subscribe into Kafka to receive the response.
Once the event has been processed by another two services(Kafka, gRPC) we subscribe into a new [Kafka] topic where
we receive the final message.
*/
func processRequest(response http.ResponseWriter, request *http.Request) {
	broker := Broker{Value: "localhost:9092"}
	publishTopic := Topic{Value: "CommunicationTopic"}
	consumeTopic := Topic{Value: "CommunicationRestTopic"}

	channel := make(chan string)
	go SubscribeConsumer(broker, consumeTopic, func(str string) {
		channel <- strings.ToUpper(str)
	})
	time.Sleep(1 * time.Second) //Time enough to subscribe
	PublishEvents(
		broker,
		publishTopic,
		"myKey", "hello world from rest")

	messageResponse := <-channel
	log.Printf("#####################################")
	log.Printf("End of transaction with Message:")
	log.Printf("%s", messageResponse)
	log.Printf("#####################################")
	writeResponse(response, messageResponse)
}

/*
Function to marshal the generic type [interface{}] into json response, then if everything is fine
we return a 200 status code with the response, otherwise a 500 error response
*/
func writeResponse(response http.ResponseWriter, t interface{}) {
	jsonResponse, err := json.Marshal(t)
	if err != nil {
		writeErrorResponse(response, err)
	} else {
		writeSuccessResponse(response, jsonResponse)
	}
}

func writeSuccessResponse(response http.ResponseWriter, jsonResponse []byte) {
	response.Header().Set("Content-Type", "application/jsonResponse")
	response.WriteHeader(http.StatusOK)
	_, _ = response.Write(jsonResponse)
}

func writeErrorResponse(response http.ResponseWriter, err error) {
	response.Header().Set("Content-Type", "application/jsonResponse")
	response.WriteHeader(http.StatusServiceUnavailable)
	errorResponse, _ := json.Marshal("Error in request since " + err.Error())
	_, _ = response.Write(errorResponse)
}
