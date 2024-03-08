package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"objs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FetchUsers(req *http.Request, client *mongo.Client) objs.FetchUsersResponse {
	requestObj := objs.FetchUsersRequest{}
	responseObj := objs.FetchUsersResponse{}
	json.NewDecoder(req.Body).Decode(&requestObj)
	var bsonData []bson.M
	doc := struct {
		ID          string `bson:"_id"`
		Name        string `bson:"name"`
		PictureData string `bson:"picturedata"`
	}{}
	switch req.Method {
	case "GET":
		responseObj.Status = false
		responseObj.Message = "Bad route config"
	case "POST":
		switch objs.DB_CHOICE {
		case "Mongo":
			collection := client.Database(objs.UserData_DB.Database).Collection(objs.UserData_DB.Collection)
			cursor, _ := collection.Find(context.TODO(), bson.D{}, options.Find().SetProjection(bson.M{
				"_id":         1,
				"name":        1,
				"picturedata": 1,
			}))
			cursor.All(context.TODO(), &bsonData)
			for _, document := range bsonData {
				byteData, _ := bson.Marshal(document)
				bson.Unmarshal(byteData, &doc)
				responseObj.Data = append(responseObj.Data, doc)
			}
			cursor.Close(context.TODO())
			responseObj.Status = true
			responseObj.Message = "Successfully Fetched User Data"
		}
	}
	return responseObj
}
