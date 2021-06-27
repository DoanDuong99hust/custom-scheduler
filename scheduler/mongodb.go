package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FreeDisk struct {
	
}

type MachineStatus struct {
	ID          primitive.ObjectID 	`bson:"_id,omitempty"`
	Machine     string				`bson:"machine,omitempty"`
	Cpu         []float64			`bson:"cpu,omitempty"`
	Memory      []float64			`bson:"memory,omitempty"`
	ReceiveNet  []float64			`bson:"receive-net,omitempty"`
	TransmitNet []float64			`bson:"transmit-net,omitempty"`
	FreeDisk	[]float64			`bson:"free_disk,omitempty"`
}

func connectToMongoDb(domain string, databaseName string, collectionName string) (*mongo.Cursor, context.Context, *mongo.Database, *mongo.Collection) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://"+domain+":27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)
	fmt.Println("Connected to MongoDB!")

	database := client.Database(databaseName)
	collection := database.Collection(collectionName)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		panic(err)
	}

	return cursor, ctx, database, collection
}

func getMachineData(
	machineStatus MachineStatus,
	machine string,
	domain string,
	databaseName string,
	collectionName string,
	) MachineStatus {

	cursor, ctx, _, _ := connectToMongoDb(domain,databaseName,collectionName)
	var status []MachineStatus
	if err := cursor.All(ctx, &status); err != nil {
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

func saveData(
	domain string,
	databaseName string,
	collectionName string,
	query string,
	)  {
	var metric MetricResponse
	resp, err4 := http.Get("http://localhost:9090/api/v1/query?query="+query)
	if err4 != nil {
		fmt.Println(err4)
	}
	// Decode the JSON body of the HTTP response into a struct
	decodeJsonDataToStruct(&metric, resp)

	// Connect to MongoDb
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://"+domain+":27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)
	fmt.Println("Connected to MongoDB!")

	database := client.Database(databaseName)
	collection := database.Collection(collectionName)

	for _, data := range metric.Data.Results{
		fmt.Printf("Node name: %s\n", data.MetricInfo["instance"])
		fmt.Printf(collectionName)
		fmt.Printf(": %s\n\n", data.MetricValue[1])

		dataConvert := fmt.Sprintf("%v", data.MetricValue[1])
		value, err := strconv.ParseFloat(dataConvert, 64)
		if err != nil {
			fmt.Println(err)
		}

		result, insertErr := collection.InsertOne(ctx, bson.D{
			{"machine", data.MetricInfo["instance"]},
			{"time", data.MetricValue[0]},
			{"value", value},
		})
		if insertErr != nil {
			fmt.Println("InsertOne ERROR:", insertErr)
		} else {
			fmt.Println("InsertOne() API result:", result)
		}
	}


}
func main()  {
	var data MachineStatus
	datam := getMachineData(data, "shisui", "localhost","test-prometheus","status").Memory[1]

	fmt.Println(datam)

	saveData("localhost", "test-prometheus", "free_disk", "(node_filesystem_avail_bytes{mountpoint=\"/\",job=\"node-exporter\"})/(1024*1024*1024)")
	saveData("localhost", "test-prometheus", "cpu_stat", "100-irate(node_cpu_seconds_total{mode=\"idle\",job=\"node-exporter\"}[10m])*100")
	saveData("localhost", "test-prometheus", "memory_stat", "((node_memory_MemTotal_bytes{job=\"node-exporter\"}-node_memory_MemAvailable_bytes{job=\"node-exporter\"})/(node_memory_MemTotal_bytes{job=\"node-exporter\"}))*100")
	saveData("localhost", "test-prometheus", "network_receive_stat", "irate(node_network_receive_bytes_total{device=\"eth0\"}[5m])/1024")
	saveData("localhost", "test-prometheus", "network_transmit_stat", "irate(node_network_transmit_bytes_total{device=\"eth0\"}[5m])/1024")
}