package main

import (
	protobuf "awesomeProject6/delivery-service/proto/delivery"
	tr "awesomeProject6/transport-service/proto/transport"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

const connstring = "testuser:testpassword@tcp(database_mysql_1:3306)/Deliveries"

type repository interface {
	Create(delivery *protobuf.Delivery) (*protobuf.Delivery, int64, error)
	GetAll() []*protobuf.Delivery
}

type Repository struct {
	mtx        sync.RWMutex
	deliveries []*protobuf.Delivery
}

func (repo *Repository) Create(delivery *protobuf.Delivery) (*protobuf.Delivery, int64, error) {
	repo.mtx.Lock()
	db, err := sql.Open("mysql", connstring)
	if err != nil {
		log.Fatalf("Cannot connect to the database: %v", err)
	}
	defer db.Close()
	container := delivery.Containers[0]
	insert, err := db.Exec(
		"insert into Deliveries.Deliveries (Container_Count,Weight,Description, Customer, User, Origin) values (?,?,?,?,?,?)",
		delivery.ContainerCount, delivery.Weight, delivery.Description, container.CId, container.UId, container.Origin)
	entityId, err := insert.LastInsertId()
	if err != nil {
		log.Fatalf("failed to add delivery: %v", err)
	}
	log.Printf("Added delivery with id: %v", entityId)
	//updated := append(repo.deliveries, delivery)
	//repo.deliveries = updated
	repo.mtx.Unlock()
	return delivery, entityId, nil
}

func (repo *Repository) GetAll() []*protobuf.Delivery {
	db, err := sql.Open("mysql", connstring)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows, err := db.Query("select * from Deliveries.Deliveries")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var deliveries []*protobuf.Delivery

	for rows.Next() {
		trId := sql.NullString{}
		cont := protobuf.Container{}
		d := protobuf.Delivery{}

		err := rows.Scan(&d.Id, &d.ContainerCount, &d.Weight, &d.Description, &trId, &cont.CId, &cont.UId, &cont.Origin)
		d.VId = ""
		if trId.Valid {
			d.VId = trId.String
		}
		d.Containers = []*protobuf.Container{&cont}

		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(d)
		deliveries = append(deliveries, &d)
	}

	return deliveries
}

type service struct {
	repo repository
}

func createDescription(delivery *protobuf.Delivery) (*tr.Description, error) {
	var description *tr.Description
	str, _ := json.Marshal(delivery)
	json.Unmarshal(str, &description)
	return description, nil
}

func (s *service) CreateDelivery(ctx context.Context, req *protobuf.Delivery) (*protobuf.Response, error) {

	log.Printf("Incoming request: Create delivery")

	delivery, entityId, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}
	const address = "transport-service:9090"
	// const address = "localhost:9090"
	conn2, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Connection error: %v", err)
	}
	log.Printf("Connected to: %v", address)
	defer conn2.Close()
	client2 := tr.NewTransportServiceClient(conn2)

	dscr, _ := createDescription(delivery)

	av, err := client2.Available(context.Background(), dscr)


	if err != nil {
		log.Fatalf("Available function error: %v", err)
	}

	log.Printf("Transport available: %v", av)

	db, err := sql.Open("mysql", connstring)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	upd, err := db.Exec("update Deliveries.Deliveries set Transport = ? where Id = ?", av.Transport.Id, entityId)
	_ = upd
	//delivery.VId = av.Transport.Id
	if err!=nil {
		log.Fatalf("Failed to update: %v", err)
	}
	log.Printf("Updated: %v", entityId)
	delivery.Id = string(entityId)

	return &protobuf.Response{Created: true, Delivery: delivery}, nil
}

func (s *service) GetDeliveries(ctx context.Context, req *protobuf.GetRequest) (*protobuf.Response, error) {
	deliveries := s.repo.GetAll()
	return &protobuf.Response{Deliveries: deliveries}, nil
}

func main() {

	repo := &Repository{}

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Listen error: %v", err)
	}
	srv := grpc.NewServer()

	protobuf.RegisterDeliveryServiceServer(srv, &service{repo})

	reflection.Register(srv)

	log.Println("Server", ":8080")
	if err := srv.Serve(l); err != nil {
		log.Fatalf("Serve error: %v", err)
	}

}
