package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"objs"
	"validation"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FetchGroupsHandlers struct{}

func (FetchGroupsHandlers) Mongo(data objs.FetchGroupsRequest, client *mongo.Client) ([]objs.GroupDataJSON, error) {
	doc := objs.GroupDataBSON{}
	groupData := []objs.GroupDataJSON{}
	var bsonData []bson.M
	var isUserPresentInGroup = false
	var err error
	collection := client.Database(objs.GroupList_DB.Database).Collection(objs.GroupList_DB.Collection)
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return groupData, err
	} else {
		cursor.All(context.TODO(), &bsonData)
		for _, document := range bsonData {
			byteData, err := bson.Marshal(document)
			if err != nil {
				return groupData, err
			} else {
				bson.Unmarshal(byteData, &doc)
				isUserPresentInGroup = false
				for _, item := range doc.Members {
					if item == data.UserID {
						isUserPresentInGroup = true
						break
					}
				}
				if isUserPresentInGroup {
					groupData = append(groupData, objs.GroupDataJSON{
						ID:      doc.ID,
						Name:    doc.Name,
						Members: doc.Members,
					})
				}
			}
		}
		cursor.Close(context.TODO())
		return groupData, nil
	}
}

func FetchGroupsHandler(req *http.Request, client *mongo.Client) objs.FetchGroupsResponse {
	requestObj := objs.FetchGroupsRequest{}
	responseObj := objs.FetchGroupsResponse{}
	json.NewDecoder(req.Body).Decode(&requestObj)
	var fetchGroupsHandlers objs.FetchGroupsMethods = FetchGroupsHandlers{}
	responseObj.Error = nil
	switch req.Method {
	case "GET":
		responseObj.Status = false
		responseObj.Message = "Bad route config"
	case "POST":
		switch objs.DB_CHOICE {
		case "Mongo":
			var mongoValidationHandlers objs.MongoValidationMethods = validation.MongoValidationHandlers{}
			validationObj := mongoValidationHandlers.ValidateUserID(requestObj.UserID, client)
			if validationObj.Status {
				groupData, err := fetchGroupsHandlers.Mongo(requestObj, client)
				if err != nil {
					responseObj.Status = false
					responseObj.Message = err.Error()
				} else {
					responseObj.Data = groupData
					responseObj.Status = true
					responseObj.Message = "Successfully Fetched Group Data"
				}
			} else {
				responseObj.Status = false
				responseObj.Message = validationObj.Message
			}
		}
	}
	return responseObj
}
