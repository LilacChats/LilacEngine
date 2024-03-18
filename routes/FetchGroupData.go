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

type FetchGroupDataHandlers struct{}

func (FetchGroupDataHandlers) Mongo(data objs.FetchGroupDataRequest, client *mongo.Client) (objs.GroupDataJSON, error) {
	var err error = nil
	bsonData := bson.M{}
	var doc objs.GroupDataBSON
	collection := client.Database(objs.GroupList_DB.Database).Collection(objs.GroupList_DB.Collection)
	groupID, err := primitive.ObjectIDFromHex(data.GroupID)
	if err != nil {
		return objs.GroupDataJSON{}, err
	} else {
		collection.FindOne(context.TODO(), bson.M{"_id": groupID}).Decode(&bsonData)
		byteData, err := bson.Marshal(bsonData)
		if err != nil {
			return objs.GroupDataJSON{}, err
		} else {
			bson.Unmarshal(byteData, &doc)
			return objs.GroupDataJSON{
				ID:      doc.ID,
				Name:    doc.Name,
				Members: doc.Members,
			}, err
		}
	}

}

func FetchGroupDataHandler(req *http.Request, client *mongo.Client) objs.FetchGroupDataResponse {
	requestObj := objs.FetchGroupDataRequest{}
	responseObj := objs.FetchGroupDataResponse{}
	json.NewDecoder(req.Body).Decode(&requestObj)
	var fetchGroupDataHandlers objs.FetchGroupDataMethods = FetchGroupDataHandlers{}
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
				groupData, err := fetchGroupDataHandlers.Mongo(requestObj, client)
				if err != nil {
					responseObj.Status = false
					responseObj.Message = err.Error()
				} else {
					responseObj.Data = groupData
					responseObj.Message = "Successfully Fetched Group Data"
					responseObj.Status = true
				}
			} else {
				responseObj.Status = false
				responseObj.Message = validationObj.Message
			}
		}
	}
	return responseObj
}
