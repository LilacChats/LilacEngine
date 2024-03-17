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

func FetchGroups(req *http.Request, client *mongo.Client) objs.FetchGroupsResponse {
	requestObj := objs.FetchGroupsRequest{}
	responseObj := objs.FetchGroupsResponse{}
	json.NewDecoder(req.Body).Decode(&requestObj)
	doc := struct {
		ID      string   `bson:"_id"`
		Name    string   `bson:"name"`
		Members []string `bson:"members"`
	}{}
	var bsonData []bson.M
	responseObj.Error = nil
	switch req.Method {
	case "GET":
		responseObj.Status = false
		responseObj.Message = "Bad route config"
	case "POST":
		switch objs.DB_CHOICE {
		case "Mongo":
			validationObj := validation.ValidateUserID(requestObj.UserID, client)
			if validationObj.Status {
				collection := client.Database(objs.GroupList_DB.Database).Collection(objs.GroupList_DB.Collection)
				cursor, _ := collection.Find(context.TODO(), bson.D{})
				cursor.All(context.TODO(), &bsonData)
				for _, document := range bsonData {
					byteData, _ := bson.Marshal(document)
					bson.Unmarshal(byteData, &doc)
					responseObj.Data = append(responseObj.Data, struct {
						ID      string   "json:\"id\""
						Name    string   "json:\"name\""
						Members []string "json:\"members\""
					}{
						ID:      doc.ID,
						Name:    doc.Name,
						Members: doc.Members,
					})
				}
				cursor.Close(context.TODO())
				responseObj.Status = true
				responseObj.Message = "Successfully Fetched Group Data"
			} else {
				responseObj.Status = false
				responseObj.Message = validationObj.Message
			}
		}
	}
	return responseObj
}
