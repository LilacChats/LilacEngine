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

func SignupUser(req *http.Request, client *mongo.Client) objs.SignupResponse {
	responseObj := objs.SignupResponse{}
	requestObj := objs.SignupRequest{}
	json.NewDecoder(req.Body).Decode(&requestObj)
	// requestObj.LoggedIn = true
	// requestObj.Offline = false
	responseObj.Error = nil
	switch req.Method {
	case "GET":
		responseObj.Status = false
		responseObj.Message = "Bad route config"
	case "POST":
		switch objs.DB_CHOICE {
		case "Mongo":
			if !validation.VerifyUserExists("email", requestObj.Email, client) {
				collection := client.Database(objs.UserData_DB.Database).Collection(objs.UserData_DB.Collection)
				collectionResponse, _ := collection.InsertOne(context.Background(), requestObj)
				responseObj.Status = true
				responseObj.Data.ID = collectionResponse.InsertedID.(primitive.ObjectID).Hex()
				responseObj.Message = "Successfully Created User"
			} else {
				responseObj.Message = "User Already Exists"
				responseObj.Status = false
			}
		}
	default:
		responseObj.Message = "Error"
		responseObj.Status = false
	}
	return responseObj
}
