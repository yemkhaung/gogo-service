package dronesevents

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type eventRepository interface {
	UpdateLastTelemetryEvent(TelemetryUpdatedEvent) error
	UpdateLastAlertEvent(AlertSignalledEvent) error
	UpdateLastPositionEvent(PositionChangedEvent) error
}

type mongoEventRepository struct {
	URL string
}

func (repo *mongoEventRepository) UpdateLastTelemetryEvent(event TelemetryUpdatedEvent) error {
	conn, err := repo.connectDb("drones", "telemetry")
	defer conn.disconnect()
	if err != nil {
		return err
	}
	result, err := conn.collection.InsertOne(conn.ctx, event)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Inserted last telemetry event: %s", result.InsertedID)
	return nil
}

func (repo *mongoEventRepository) UpdateLastAlertEvent(event AlertSignalledEvent) error {
	conn, err := repo.connectDb("drones", "alert")
	defer conn.disconnect()
	if err != nil {
		return err
	}
	result, err := conn.collection.InsertOne(conn.ctx, event)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Inserted last alert event: %s", result.InsertedID)
	return nil
}

func (repo *mongoEventRepository) UpdateLastPositionEvent(event PositionChangedEvent) error {
	conn, err := repo.connectDb("drones", "position")
	defer conn.disconnect()
	if err != nil {
		return err
	}
	result, err := conn.collection.InsertOne(conn.ctx, event)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Inserted last position event: %s", result.InsertedID)
	return nil
}

type dbConnection struct {
	client     *mongo.Client
	ctx        context.Context
	cancelCtx  context.CancelFunc
	collection *mongo.Collection
}

func (conn dbConnection) disconnect() {
	if conn.ctx != nil {
		_ = conn.client.Disconnect(conn.ctx)
	}
}

func (repo *mongoEventRepository) connectDb(dbName string, colName string) (conn dbConnection, err error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(repo.URL))
	if err != nil {
		log.Println(err)
		return conn, err
	}
	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Println(err)
		cancelCtx()
		return conn, err
	}
	collection := client.Database(dbName).Collection(colName)
	return dbConnection{client, ctx, cancelCtx, collection}, nil
}
