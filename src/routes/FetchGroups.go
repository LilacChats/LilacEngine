package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"objs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func FetchGroups(req *http.Request, client *mongo.Client) objs.FetchGroupsResponse {
	requestObj := objs.FetchGroupsRequest{}
	responseObj := objs.FetchGroupsResponse{}
	json.NewDecoder(req.Body).Decode(&requestObj)
	doc := struct {
		ID   string `bson:"_id"`
		Name string `bson:"name"`
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
			collection := client.Database(objs.GroupList_DB.Database).Collection(objs.GroupList_DB.Collection)
			cursor, _ := collection.Find(context.TODO(), bson.D{})
			cursor.All(context.TODO(), &bsonData)
			for _, document := range bsonData {
				byteData, _ := bson.Marshal(document)
				bson.Unmarshal(byteData, &doc)
				responseObj.Data = append(responseObj.Data, doc)
			}
			cursor.Close(context.TODO())
			responseObj.Status = true
			responseObj.Message = "Successfully Fetched Group Data"
		}
	}
	return responseObj
}
