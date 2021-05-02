package gogoservice

import (
	"context"
	"log"
	"time"

	"github.com/cloudnativego/gogo-engine"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type persistMatchRepository struct {
	dbURL string
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

func newPersistRepository(dbURL string) *persistMatchRepository {
	return &persistMatchRepository{dbURL}
}

func (pr *persistMatchRepository) connectDb() (conn dbConnection, err error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(pr.dbURL))
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
	collection := client.Database("test").Collection("matches")
	return dbConnection{client, ctx, cancelCtx, collection}, nil
}

func (pr *persistMatchRepository) addMatch(match gogo.Match) (err error) {
	// save input Match to DB
	conn, err := pr.connectDb()
	defer conn.disconnect()
	if err != nil {
		return err
	}
	result, err := conn.collection.InsertOne(conn.ctx, match)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Inserted match: ", result.InsertedID)
	return nil
}

func (pr *persistMatchRepository) getMatch(id string) (match gogo.Match, err error) {
	conn, err := pr.connectDb()
	defer conn.disconnect()
	if err != nil {
		return match, err
	}
	err = conn.collection.FindOne(conn.ctx, bson.D{{Key: "id", Value: id}}).Decode(&match)
	if err != nil {
		log.Println(err)
		return match, err
	}
	return match, nil
}

func (pr *persistMatchRepository) getMatches() (mlist []gogo.Match, err error) {
	conn, err := pr.connectDb()
	defer conn.disconnect()
	if err != nil {
		return mlist, err
	}
	cursor, err := conn.collection.Find(conn.ctx, bson.M{})
	if err != nil {
		log.Println(err)
		return mlist, err
	}
	err = cursor.All(conn.ctx, &mlist)
	if err != nil {
		log.Println(err)
		return mlist, err
	}
	return mlist, nil
}
