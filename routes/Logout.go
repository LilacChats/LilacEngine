package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"objs"
	"validation"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LogoutHandlers struct{}

func (LogoutHandlers) Mongo(data objs.LogoutRequest, client *mongo.Client) error {
	collection := client.Database(objs.UserData_DB.Database).Collection(objs.UserData_DB.Collection)
	_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": data.UserID}, bson.M{"$set": bson.M{"online": false}})
	return err
}

func LogoutHandler(req *http.Request, client *mongo.Client) objs.LogoutResponse {
	var requestObj = objs.LogoutRequest{}
	var responseObj = objs.LogoutResponse{}
	json.NewDecoder(req.Body).Decode(&requestObj)
	var logoutHandlers objs.LogoutMethods = LogoutHandlers{}
	responseObj.Error = nil
	switch req.Method {
	case "GET":
		responseObj.Message = "Bad Route Config"
		responseObj.Status = false
	case "POST":
		switch objs.DB_CHOICE {
		case "Mongo":
			var mongoValidationHandlers objs.MongoValidationMethods = validation.MongoValidationHandlers{}
			validationObj := mongoValidationHandlers.ValidateUserID(requestObj.UserID, client)
			if validationObj.Status {
				logoutStatus := logoutHandlers.Mongo(requestObj, client)
				if logoutStatus == nil {
					responseObj.Status = true
					responseObj.Message = "Successfully Logged Out User"
				} else {
					responseObj.Status = false
					responseObj.Message = logoutStatus.Error()
				}
			} else {
				responseObj.Status = false
				responseObj.Message = validationObj.Message
			}
		}
	}
	return responseObj
}
