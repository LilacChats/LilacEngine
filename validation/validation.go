package validation

import (
	"context"
	"fmt"
	"objs"

	"github.com/charmbracelet/lipgloss"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoValidationHandlers struct{}

func DisplayError(message string) {
	fmt.Print(lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1).Background(lipgloss.Color("9")).Render("Error"))
	fmt.Println(lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1).Background(lipgloss.Color("5")).Render(message))
}

func (MongoValidationHandlers) VerifyUserExists(keyType string, key string, client *mongo.Client) bool {
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

func (MongoValidationHandlers) VerifyGroupExists(groupID string, client *mongo.Client) bool {
	collection := client.Database(objs.GroupList_DB.Database).Collection(objs.GroupList_DB.Collection)
	var bsonData objs.GroupDataBSON
	objectID, _ := primitive.ObjectIDFromHex(groupID)
	err := collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&bsonData)
	if err != nil {
		return false
	} else {
		return true
	}
}

func (MongoValidationHandlers) ValidateUserEmail(email string, client *mongo.Client) objs.StandardResponse {
	var mongoValidationHandlers objs.MongoValidationMethods = MongoValidationHandlers{}
	if mongoValidationHandlers.VerifyUserExists("email", email, client) {
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

func (MongoValidationHandlers) ValidateUserID(id string, client *mongo.Client) objs.StandardResponse {
	var mongoValidationHandlers objs.MongoValidationMethods = MongoValidationHandlers{}
	if mongoValidationHandlers.VerifyUserExists("id", id, client) {
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

func (MongoValidationHandlers) ValidateGroup(id string, groupID string, client *mongo.Client) objs.StandardResponse {
	var mongoValidationHandlers objs.MongoValidationMethods = MongoValidationHandlers{}
	isUserIDValid := mongoValidationHandlers.VerifyUserExists("id", id, client)
	isGroupIDValid := mongoValidationHandlers.VerifyGroupExists(groupID, client)
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
