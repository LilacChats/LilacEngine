package validation

import (
	"context"
	"objs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func VerifyUserExists(keyType string, key string, client *mongo.Client) bool {
	collection := client.Database(objs.UserData_DB.Database).Collection(objs.UserData_DB.Collection)
	var bsonData bson.D
	var filter bson.M
	switch keyType {
	case "id":
		objectID, _ := primitive.ObjectIDFromHex(key)
		filter = bson.M{"_id": objectID}
	case "email":
		filter = bson.M{"email": key}
	}
	err := collection.FindOne(context.TODO(), filter).Decode(&bsonData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
	} else {
		return true
	}
	return false
}

func VerifyGroupExists(groupID string, client *mongo.Client) bool {
	collection := client.Database(objs.GroupList_DB.Database).Collection(objs.GroupList_DB.Collection)
	var bsonData objs.GroupData
	objectID, _ := primitive.ObjectIDFromHex(groupID)
	err := collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&bsonData)
	if err != nil {
		return false
	} else {
		return true
	}
}

func ValidateUserEmail(email string, client *mongo.Client) objs.StandardResponse {
	if VerifyUserExists("email", email, client) {
		return objs.StandardResponse{
			Status: true,
		}
	} else {
		return objs.StandardResponse{
			Message: "User Does Not Exist",
			Status:  false,
		}
	}
}

func ValidateUserID(id string, client *mongo.Client) objs.StandardResponse {
	if VerifyUserExists("id", id, client) {
		return objs.StandardResponse{
			Status: true,
		}
	} else {
		return objs.StandardResponse{
			Message: "User Does Not Exist",
			Status:  false,
		}
	}
}

func ValidateGroup(id string, groupID string, client *mongo.Client) objs.StandardResponse {
	isUserIDValid := VerifyUserExists("id", id, client)
	isGroupIDValid := VerifyGroupExists(groupID, client)
	if isUserIDValid && isGroupIDValid {
		return objs.StandardResponse{
			Status: true,
		}
	} else {
		var responseMessage string
		if !isUserIDValid {
			responseMessage = "Invalid User ID"
		}
		responseMessage += "\n"
		if !isGroupIDValid {
			responseMessage = "Invalid Group ID\n"
		}
		return objs.StandardResponse{
			Message: responseMessage,
			Status:  false,
		}
	}
}
