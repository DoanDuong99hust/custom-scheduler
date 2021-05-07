package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type MachineStatus struct {
	ID          primitive.ObjectID 	`bson:"_id,omitempty"`
	Machine     string				`bson:"machine,omitempty"`
	Cpu         []float64			`bson:"cpu,omitempty"`
	Memory      []float64			`bson:"memory,omitempty"`
	ReceiveNet  []float64			`bson:"receive-net,omitempty"`
	TransmitNet []float64		`bson:"transmit-net,omitempty"`

}

func getMongoDbData(machineStatus MachineStatus, machine string) MachineStatus {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://192.168.100.137:27017")

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

	for i := len(status)-1; i >= 0; i-- {
		if status[i].Machine == machine {
			machineStatus = status[i]
			break
		}
	}
	return machineStatus
}

// func main()  {
// 	var data MachineStatus
// 	datam := getMongoDbData(data, "worker02").Memory[1]

// 	fmt.Println(datam)
// }