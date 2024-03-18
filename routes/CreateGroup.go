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

type CreateGroupHandlers struct{}

func (CreateGroupHandlers) Mongo(data objs.CreateGroupRequest, client *mongo.Client) (string, error) {
	var newCollectionID string
	var err error = nil
	collection := client.Database(objs.GroupList_DB.Database).Collection(objs.GroupList_DB.Collection)
	collectionResponse, err := collection.InsertOne(context.TODO(), struct {
		Name    string
		Members []string
	}{Name: data.GroupName,
		Members: data.Members,
	})
	if err != nil {
		return "", err
	} else {
		newCollectionID = collectionResponse.InsertedID.(primitive.ObjectID).Hex()
		err := client.Database(objs.GroupList_DB.Database).CreateCollection(context.TODO(), newCollectionID)
		if err != nil {
			return "", err
		} else {
			return newCollectionID, nil
		}
	}

}

func CreateGroupHandler(req *http.Request, client *mongo.Client) objs.CreateGroupResponse {
	var requestObj = objs.CreateGroupRequest{}
	var responseObj = objs.CreateGroupResponse{}
	json.NewDecoder(req.Body).Decode(&requestObj)
	var createGroupHandlers objs.CreateGroupMethods = CreateGroupHandlers{}
	responseObj.Error = nil
	switch req.Method {
	case "GET":
		responseObj.Message = "Bad route config"
		responseObj.Status = false
	case "POST":
		switch objs.DB_CHOICE {
		case "Mongo":
			var mongoValidationHandlers objs.MongoValidationMethods = validation.MongoValidationHandlers{}
			validationObj := mongoValidationHandlers.ValidateUserID(requestObj.UserID, client)
			if validationObj.Status {
				newCollectionID, err := createGroupHandlers.Mongo(requestObj, client)
				if err != nil {
					responseObj.Status = false
					responseObj.Message = err.Error()
				} else {
					responseObj.Data.GroupID = newCollectionID
					responseObj.Status = true
					responseObj.Message = "Successfully Created Group"
				}
			} else {
				responseObj.Status = false
				responseObj.Message = validationObj.Message
			}
		}
	}
	return responseObj
}
