package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"validation"

	"objs"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateGroup(req *http.Request, client *mongo.Client) objs.CreateGroupResponse {
	var requestObj = objs.CreateGroupRequest{}
	var responseObj = objs.CreateGroupResponse{}
	var newCollectionID string = ""
	json.NewDecoder(req.Body).Decode(&requestObj)
	responseObj.Error = nil
	switch req.Method {
	case "GET":
		responseObj.Message = "Bad route config"
		responseObj.Status = false
	case "POST":
		switch objs.DB_CHOICE {
		case "Mongo":
			validationObj := validation.ValidateUserID(requestObj.UserID, client)
			if validationObj.Status {
				collection := client.Database(objs.GroupList_DB.Database).Collection(objs.GroupList_DB.Collection)
				collectionResponse, _ := collection.InsertOne(context.TODO(), struct {
					Name    string
					Members []string
				}{Name: requestObj.GroupName,
					Members: requestObj.Members,
				})
				newCollectionID = collectionResponse.InsertedID.(primitive.ObjectID).Hex()
				client.Database(objs.GroupList_DB.Database).CreateCollection(context.TODO(), newCollectionID)
				responseObj.Data.GroupID = newCollectionID
				responseObj.Status = true
				responseObj.Message = "Successfully Created Group"
			} else {
				responseObj.Status = false
				responseObj.Message = validationObj.Message
			}
		}
	}
	return responseObj
}
