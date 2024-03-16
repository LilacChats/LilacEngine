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

func FetchGroupData(req *http.Request, client *mongo.Client) objs.FetchGroupDataResponse {
	requestObj := objs.FetchGroupDataRequest{}
	responseObj := objs.FetchGroupDataResponse{}
	json.NewDecoder(req.Body).Decode(&requestObj)
	bsonData := bson.M{}
	var doc objs.GroupData
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
				groupID, _ := primitive.ObjectIDFromHex(requestObj.GroupID)
				collection.FindOne(context.TODO(), bson.M{"_id": groupID}).Decode(&bsonData)
				byteData, _ := bson.Marshal(bsonData)
				bson.Unmarshal(byteData, &doc)
				responseObj.Data = doc
				responseObj.Message = "Successfully Fetched Group Data"
				responseObj.Status = true
			} else {
				responseObj.Status = false
				responseObj.Message = validationObj.Message
			}
		}
	}
	return responseObj
}
