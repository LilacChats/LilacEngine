package routes

import (
	"context"
	"objs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindDMChannelID(id1 string, id2 string, client *mongo.Client) (string, error) {
	var bsonData bson.M
	var channelData objs.DMChannel
	collection := client.Database(objs.ChannelList_DB.Database).Collection(objs.ChannelList_DB.Collection)
	err := collection.FindOne(context.TODO(), bson.M{"$or": bson.A{
		bson.M{"id1": id1, "id2": id2},
		bson.M{"id2": id1, "id1": id2},
	}}).Decode(&bsonData)
	if err != nil {
		collectionData, err := collection.InsertOne(context.TODO(), bson.M{"id1": id1, "id2": id2})
		if err != nil {
			return "", err
		} else {
			return collectionData.InsertedID.(primitive.ObjectID).Hex(), nil
		}
	} else {
		byteData, err := bson.Marshal(bsonData)
		if err != nil {
			return "", err
		} else {
			err := bson.Unmarshal(byteData, &channelData)
			if err != nil {
				return "", err
			} else {
				return channelData.ChannelID, nil
			}
		}
	}
}
