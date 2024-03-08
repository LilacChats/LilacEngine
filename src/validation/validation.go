package validation

import (
	"context"
	"objs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func VerifyUserExists(keyType string, key string, client *mongo.Client) bool {
	collection := client.Database(objs.UserData_DB.Database).Collection(objs.UserData_DB.Collection)
	var bsonData bson.D
	var err error
	filter := bson.D{{keyType, key}}
	err = collection.FindOne(context.TODO(), filter).Decode(&bsonData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
	} else {
		return true
	}
	return false
}
