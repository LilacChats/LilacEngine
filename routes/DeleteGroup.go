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

func DeleteGroup(req *http.Request, client *mongo.Client) objs.DeleteGroupResponse {
	requestObj := objs.DeleteGroupRequest{}
	responseObj := objs.DeleteGroupResponse{}
	json.NewDecoder(req.Body).Decode(&requestObj)
	switch req.Method {
	case "GET":
		responseObj.Message = "Bad route config"
		responseObj.Status = false
	case "POST":
		switch objs.DB_CHOICE {
		case "Mongo":
			validationObj := validation.ValidateGroup(requestObj.UserID, requestObj.GroupID, client)
			if validationObj.Status {
				groupCollection := client.Database(objs.DATABASE).Collection(requestObj.GroupID)
				groupCollection.Drop(context.TODO())
				groupListCollection := client.Database(objs.DATABASE).Collection(objs.GroupList_DB.Collection)
				objectID, _ := primitive.ObjectIDFromHex(requestObj.GroupID)
				groupListCollection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
				responseObj.Status = true
				responseObj.Message = "Successfully Deleted Group"
			} else {
				responseObj.Status = false
				responseObj.Message = validationObj.Message
			}
		}
	}
	return responseObj
}
