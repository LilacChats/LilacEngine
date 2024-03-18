package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"objs"
	"validation"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DeleteGroupHandlers struct{}

func (DeleteGroupHandlers) Mongo(data objs.DeleteGroupRequest, client *mongo.Client) error {
	var err error
	groupCollection := client.Database(objs.DATABASE).Collection(data.GroupID)
	groupCollection.Drop(context.TODO())
	groupListCollection := client.Database(objs.DATABASE).Collection(objs.GroupList_DB.Collection)
	objectID, err := primitive.ObjectIDFromHex(data.GroupID)
	if err != nil {
		return err
	} else {
		_, err = groupListCollection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
		if err != nil {
			return err
		} else {
			return nil
		}
	}
}

func DeleteGroupHandler(req *http.Request, client *mongo.Client) objs.DeleteGroupResponse {
	requestObj := objs.DeleteGroupRequest{}
	responseObj := objs.DeleteGroupResponse{}
	json.NewDecoder(req.Body).Decode(&requestObj)
	var deleteGroupHandlers objs.DeleteGroupMethods = DeleteGroupHandlers{}
	switch req.Method {
	case "GET":
		responseObj.Message = "Bad route config"
		responseObj.Status = false
	case "POST":
		switch objs.DB_CHOICE {
		case "Mongo":
			var mongoValidationHandlers objs.MongoValidationMethods = validation.MongoValidationHandlers{}
			validationObj := mongoValidationHandlers.ValidateGroup(requestObj.UserID, requestObj.GroupID, client)
			if validationObj.Status {
				err := deleteGroupHandlers.Mongo(requestObj, client)
				if err != nil {
					responseObj.Status = false
					responseObj.Message = err.Error()
				} else {
					responseObj.Status = true
					responseObj.Message = "Successfully Deleted Group"
				}
			} else {
				responseObj.Status = false
				responseObj.Message = validationObj.Message
			}
		}
	}
	return responseObj
}
