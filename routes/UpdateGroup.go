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

func UpdateGroup(req *http.Request, client *mongo.Client) objs.UpdateGroupResponse {
	requestObj := objs.UpdateGroupRequest{}
	responseObj := objs.UpdateGroupResponse{}
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
				collection := client.Database(objs.GroupList_DB.Database).Collection(objs.GroupList_DB.Collection)
				objectID, _ := primitive.ObjectIDFromHex(requestObj.GroupID)
				collection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, bson.M{"$set": bson.M{"name": requestObj.Name, "members": requestObj.Members}})
				responseObj.Status = true
				responseObj.Message = "Successfully Updated Group Data"
			} else {
				responseObj.Status = false
				responseObj.Message = validationObj.Message
			}
		}

	}
	return responseObj
}
