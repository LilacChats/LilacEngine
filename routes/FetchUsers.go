package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"objs"
	"validation"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FetchUsersHandlers struct{}

func (FetchUsersHandlers) Mongo(data objs.FetchUsersRequest, client *mongo.Client) ([]struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	PictureData string `json:"pictureData"`
}, error) {
	var err error
	var bsonData []bson.M
	doc := struct {
		ID          string `bson:"_id"`
		Name        string `bson:"name"`
		PictureData string `bson:"pictureData"`
	}{}
	usersData := []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		PictureData string `json:"pictureData"`
	}{}
	collection := client.Database(objs.UserData_DB.Database).Collection(objs.UserData_DB.Collection)
	cursor, err := collection.Find(context.TODO(), bson.D{}, options.Find().SetProjection(bson.M{
		"_id":         1,
		"name":        1,
		"picturedata": 1,
	}))
	if err != nil {
		return []struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			PictureData string `json:"pictureData"`
		}{}, err
	} else {
		cursor.All(context.TODO(), &bsonData)
		for _, document := range bsonData {
			byteData, err := bson.Marshal(document)
			if err != nil {
				return []struct {
					ID          string `json:"id"`
					Name        string `json:"name"`
					PictureData string `json:"pictureData"`
				}{}, err
			} else {
				bson.Unmarshal(byteData, &doc)
				usersData = append(usersData, struct {
					ID          string `json:"id"`
					Name        string `json:"name"`
					PictureData string `json:"pictureData"`
				}{
					ID:          doc.ID,
					Name:        doc.Name,
					PictureData: doc.PictureData,
				})
			}
		}
		cursor.Close(context.TODO())
		return usersData, nil
	}
}

func FetchUsersHandler(req *http.Request, client *mongo.Client) objs.FetchUsersResponse {
	requestObj := objs.FetchUsersRequest{}
	responseObj := objs.FetchUsersResponse{}
	json.NewDecoder(req.Body).Decode(&requestObj)
	var fetchUsersHandlers objs.FetchUsersMethods = FetchUsersHandlers{}
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
				usersData, err := fetchUsersHandlers.Mongo(requestObj, client)
				if err != nil {
					responseObj.Status = false
					responseObj.Message = err.Error()
				} else {
					responseObj.Data = usersData
					responseObj.Status = true
					responseObj.Message = "Successfully Fetched User Data"
				}
			}

		}
	}
	return responseObj
}
