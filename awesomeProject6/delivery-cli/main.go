package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"context"

	protobuf "awesomeProject6/delivery-service/proto/delivery"
	"google.golang.org/grpc"
)

const defaultFilename = "delivery.json"

func parseFile(file string) (*protobuf.Delivery, error) {
	var delivery *protobuf.Delivery
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &delivery)
	return delivery, err
}



func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Connection error: %v", err)
	}
	defer conn.Close()
	client := protobuf.NewDeliveryServiceClient(conn)



	file := defaultFilename
	if len(os.Args) > 1 {
		file = os.Args[1]
	}

	delivery, err := parseFile(file)

	if err != nil {
		log.Fatalf("Parsing error: %v", err)
	}

	create, err := client.CreateDelivery(context.Background(), delivery)
	if err != nil {
		log.Fatalf("CreateDelivery error: %v", err)
	}
	log.Printf("Delivery created: %t", create.Created)


	getAll, err := client.GetDeliveries(context.Background(), &protobuf.GetRequest{})
	if err != nil {
		log.Fatalf("Get Deliveries error: %v", err)
	}
	for _, v := range getAll.Deliveries {
		log.Println(v)
	}

}
