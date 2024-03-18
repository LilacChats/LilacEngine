package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"objs"
	"validation"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SignupHandlers struct{}

func (SignupHandlers) Mongo(data objs.SignupRequest, client *mongo.Client) (string, error) {
	collection := client.Database(objs.UserData_DB.Database).Collection(objs.UserData_DB.Collection)
	collectionResponse, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return "", err
	} else {
		return collectionResponse.InsertedID.(primitive.ObjectID).Hex(), nil
	}
}

func SignupHandler(req *http.Request, client *mongo.Client) objs.SignupResponse {
	responseObj := objs.SignupResponse{}
	requestObj := objs.SignupRequest{}
	var signupHandlers objs.SignupMethods = SignupHandlers{}
	json.NewDecoder(req.Body).Decode(&requestObj)
	responseObj.Error = nil
	switch req.Method {
	case "GET":
		responseObj.Status = false
		responseObj.Message = "Bad route config"
	case "POST":
		switch objs.DB_CHOICE {
		case "Mongo":
			var mongoValidationHandlers objs.MongoValidationMethods = validation.MongoValidationHandlers{}
			validationObj := mongoValidationHandlers.ValidateUserEmail(requestObj.Email, client)
			if !validationObj.Status {
				signupResponse, err := signupHandlers.Mongo(requestObj, client)
				if err != nil {
					responseObj.Status = false
					responseObj.Message = err.Error()
				} else {
					responseObj.Status = true
					responseObj.Data.ID = signupResponse
					responseObj.Message = "Successfully Created User"
				}
			} else {
				responseObj.Message = validationObj.Message
				responseObj.Status = false
			}
		}
	default:
		responseObj.Message = "Error"
		responseObj.Status = false
	}
	return responseObj
}
