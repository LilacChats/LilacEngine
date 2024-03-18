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

type UpdateGroupHandlers struct{}

func (UpdateGroupHandlers) Mongo(data objs.UpdateGroupRequest, client *mongo.Client) error {
	var err error
	var objectID primitive.ObjectID
	var deleteGroupHandlers objs.DeleteGroupMethods = DeleteGroupHandlers{}
	if len(data.Members) == 0 {
		deleteRequest := objs.DeleteGroupRequest{
			UserID:  data.UserID,
			GroupID: data.GroupID,
		}
		err = deleteGroupHandlers.Mongo(deleteRequest, client)
		if err != nil {
			return err
		} else {
			return nil
		}
	} else {
		collection := client.Database(objs.GroupList_DB.Database).Collection(objs.GroupList_DB.Collection)
		objectID, err = primitive.ObjectIDFromHex(data.GroupID)
		if err != nil {
			return err
		} else {
			_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, bson.M{"$set": bson.M{"name": data.Name, "members": data.Members}})
			if err != nil {
				return err
			} else {
				return nil
			}
		}

	}
}

func UpdateGroupHandler(req *http.Request, client *mongo.Client) objs.UpdateGroupResponse {
	requestObj := objs.UpdateGroupRequest{}
	responseObj := objs.UpdateGroupResponse{}
	json.NewDecoder(req.Body).Decode(&requestObj)
	var updateGroupHandlers objs.UpdateGroupMethods = UpdateGroupHandlers{}
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
				err := updateGroupHandlers.Mongo(requestObj, client)
				if err != nil {
					responseObj.Status = false
					responseObj.Message = err.Error()
				} else {
					responseObj.Status = true
					if len(requestObj.Members) == 0 {
						responseObj.Message = "Group Deleted as empty members list received"
					} else {
						responseObj.Message = "Group Data updated successfully!"
					}
				}
			} else {
				responseObj.Status = false
				responseObj.Message = validationObj.Message
			}
		}

	}
	return responseObj
}
