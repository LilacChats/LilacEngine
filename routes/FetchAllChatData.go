package routes

import (
	"context"
	"objs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FetchChatDataHandlers struct{}

func (FetchChatDataHandlers) Mongo(channelID string, client *mongo.Client) ([]objs.Message, error) {
	var bsonData []bson.M
	messages := []objs.Message{}
	var message objs.MessageBSON
	collection := client.Database(objs.DM_CHANNELS_DATABASE).Collection(channelID)
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return []objs.Message{}, err
	} else {
		cursor.All(context.TODO(), &bsonData)
		for _, document := range bsonData {
			byteData, _ := bson.Marshal(document)
			bson.Unmarshal(byteData, &message)
			messages = append(messages, objs.Message{
				MessageID:  message.MessageID,
				SenderID:   message.SenderID,
				ReceiverID: message.ReceiverID,
				Timestamp:  message.Timestamp,
				Message:    message.Message,
			})
		}
		return messages, nil
	}
}

func FetchAllChatsHandler(id1 string, id2 string, client *mongo.Client) []objs.Message {
	var fetchChatDataHandlers objs.FetchAllChatDataMethods = FetchChatDataHandlers{}
	channelID, err := FindDMChannelID(id1, id2, client)
	if err != nil {
		return []objs.Message{}
	} else {
		messages, _ := fetchChatDataHandlers.Mongo(channelID, client)
		return messages
	}
}
