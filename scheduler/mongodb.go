package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type MachineStatus struct {
	ID          primitive.ObjectID 	`bson:"_id,omitempty"`
	Machine     string				`bson:"machine,omitempty"`
	Cpu         []interface{}		`bson:"cpu,omitempty"`
	Memory      []interface{}		`bson:"memory,omitempty"`
	ReceiveNet  []interface{}		`bson:"receive-net,omitempty"`
	TransmitNet []interface{}		`bson:"transmit-net,omitempty"`

}

func getMongoDbData(machineStatus MachineStatus, machine string) MachineStatus {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)
	fmt.Println("Connected to MongoDB!")

	database := client.Database("test-prometheus")
	statusCollection := database.Collection("status")

	cursor, err := statusCollection.Find(ctx, bson.M{})
	if err != nil {
		panic(err)
	}

	var status []MachineStatus
	if err = cursor.All(ctx, &status); err != nil {
		panic(err)
	}

	for i := 0; i < len(status); i++ {
		if status[i].Machine == machine {
			fmt.Println(status[i].Machine)
			machineStatus = status[i]
			break
		}
	}
	return machineStatus
}

//func main()  {
//	var data MachineStatus
//	fmt.Println(getMongoDbData(data, "shisui"))
//}