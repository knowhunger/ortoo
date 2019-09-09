package mongodb

import (
	integration "github.com/knowhunger/ortoo/integration_test"
	"log"
	"testing"
)

func TestMongo(t *testing.T) {
	conf := integration.NewTestMongoDBConfig()
	mongo, err := New(conf)
	if err != nil {
		log.Fatalf("fail to create mongoDB instance:%v", err)
	}
	if err := mongo.InitializeCollections(); err != nil {
		log.Fatalf("fail to initialize collection:%v", err)
	}

	//collection := mongo.db.Collection("hello")
	//ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//res, err := collection.InsertOne(ctx, bson.M{"name": "pi", "value": 3.14159})
	//id := res.InsertedID
	//fmt.Println(id)
}
