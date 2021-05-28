package dronesevents

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type eventRepository interface {
	UpdateLastTelemetryEvent(TelemetryRecord) error
	UpdateLastAlertEvent(AlertRecord) error
	UpdateLastPositionEvent(PositionRecord) error
}

type mongoEventRepository struct {
	URL string
}

func getFilterQuery(droneID string) interface{} {
	return bson.D{{Key: "drone_id", Value: droneID}}
}

func (repo *mongoEventRepository) UpdateLastTelemetryEvent(record TelemetryRecord) error {
	options := options.Replace().SetUpsert(true)

	conn, err := repo.connectDb("drones", "telemetry")
	defer conn.disconnect()
	if err != nil {
		return err
	}
	result, err := conn.collection.ReplaceOne(conn.ctx, getFilterQuery(record.DroneID), record, options)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Inserted last telemetry event: %s", result.UpsertedID)
	return nil
}

func (repo *mongoEventRepository) UpdateLastAlertEvent(record AlertRecord) error {
	options := options.Replace().SetUpsert(true)

	conn, err := repo.connectDb("drones", "alert")
	defer conn.disconnect()
	if err != nil {
		return err
	}
	result, err := conn.collection.ReplaceOne(conn.ctx, getFilterQuery(record.DroneID), record, options)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Inserted last telemetry event: %s", result.UpsertedID)
	return nil
}

func (repo *mongoEventRepository) UpdateLastPositionEvent(record PositionRecord) error {
	options := options.Replace().SetUpsert(true)

	conn, err := repo.connectDb("drones", "position")
	defer conn.disconnect()
	if err != nil {
		return err
	}
	result, err := conn.collection.ReplaceOne(conn.ctx, getFilterQuery(record.DroneID), record, options)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Inserted last telemetry event: %s", result.UpsertedID)
	return nil
}

func (repo *mongoEventRepository) GetLastTelemetry(droneID string) (TelemetryRecord, error) {
	conn, err := repo.connectDb("drones", "telemetry")
	defer conn.disconnect()
	if err != nil {
		return TelemetryRecord{}, err
	}
	var telemetry TelemetryRecord
	err = conn.collection.FindOne(conn.ctx, bson.D{{Key: "drone_id", Value: droneID}}).Decode(&telemetry)
	if err != nil {
		log.Println(err)
		return TelemetryRecord{}, err
	}
	return telemetry, nil
}

func (repo *mongoEventRepository) GetLastAlert(droneID string) (AlertRecord, error) {
	conn, err := repo.connectDb("drones", "alert")
	defer conn.disconnect()
	if err != nil {
		return AlertRecord{}, err
	}
	var alert AlertRecord
	err = conn.collection.FindOne(conn.ctx, bson.D{{Key: "drone_id", Value: droneID}}).Decode(&alert)
	if err != nil {
		log.Println(err)
		return AlertRecord{}, err
	}
	return alert, nil
}

func (repo *mongoEventRepository) GetLastPosition(droneID string) (PositionRecord, error) {
	conn, err := repo.connectDb("drones", "position")
	defer conn.disconnect()
	if err != nil {
		return PositionRecord{}, err
	}
	var position PositionRecord
	err = conn.collection.FindOne(conn.ctx, bson.D{{Key: "drone_id", Value: droneID}}).Decode(&position)
	if err != nil {
		log.Println(err)
		return PositionRecord{}, err
	}
	return position, nil
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
