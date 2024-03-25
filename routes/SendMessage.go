package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"objs"
	"validation"

	"go.mongodb.org/mongo-driver/mongo"
)

type SendMessageHandlers struct{}

func (SendMessageHandlers) Mongo(data objs.SendMessageRequest, client *mongo.Client) error {
	channelID, err := FindDMChannelID(data.SenderID, data.ReceiverID, client)
	if err != nil {
		return err
	} else {
		collection := client.Database(objs.DM_CHANNELS_DATABASE).Collection(channelID)
		_, err := collection.InsertOne(context.TODO(), objs.Message{
			Message:    data.Message,
			SenderID:   data.SenderID,
			ReceiverID: data.ReceiverID,
			Timestamp:  data.Timestamp,
		})
		return err
	}
}

func SendMessageHandler(req *http.Request, client *mongo.Client) objs.SendMessageResponse {
	requestObj := objs.SendMessageRequest{}
	responseObj := objs.SendMessageResponse{}
	json.NewDecoder(req.Body).Decode(&requestObj)
	var sendMessageHandlers objs.SendMessageMethods = SendMessageHandlers{}
	switch req.Method {
	case "GET":
		responseObj.Message = "Bad route config"
		responseObj.Status = false
	case "POST":
		switch objs.DB_CHOICE {
		case "Mongo":
			var mongoValidationHandlers objs.MongoValidationMethods = validation.MongoValidationHandlers{}
			senderValidationObj := mongoValidationHandlers.ValidateUserID(requestObj.SenderID, client)
			receiverValidationObj := mongoValidationHandlers.ValidateUserID(requestObj.ReceiverID, client)
			if senderValidationObj.Status && receiverValidationObj.Status {
				err := sendMessageHandlers.Mongo(requestObj, client)
				if err != nil {
					responseObj.Status = false
					responseObj.Message = err.Error()
				} else {
					responseObj.Status = true
					responseObj.Message = "Message Sent Successfully"
				}
			} else {
				responseObj.Status = false
				responseObj.Message = senderValidationObj.Message + "\n" + receiverValidationObj.Message
			}
		}
	}
	return responseObj
}
