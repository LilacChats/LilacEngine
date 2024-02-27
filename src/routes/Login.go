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

func bsonToStruct(bsonData bson.D) objs.UserData {
	var dataObject = struct {
		ID       string `bson:"_id"`
		UserData objs.UserData
	}{}
	byteData, _ := bson.Marshal(bsonData)
	bson.Unmarshal(byteData, &dataObject)
	dataObject.UserData.ID = dataObject.ID
	return dataObject.UserData
}

func LoginUser(req *http.Request, client *mongo.Client) objs.LoginResponse {
	var requestObj = objs.LoginRequest{}
	var responseObj = objs.LoginResponse{}
	var bsonData bson.D
	json.NewDecoder(req.Body).Decode(&requestObj)
	responseObj.Error = nil
	switch req.Method {
	case "GET":
		responseObj.Message = "Bad route config"
		responseObj.Status = false
	case "POST":
		switch objs.DB_CHOICE {
		case "Mongo":
			if validation.VerifyUserExists("email", requestObj.Email, client) {
				collection := client.Database(objs.UserData_DB.Database).Collection(objs.UserData_DB.Collection)
				filter := bson.M{"userdata.email": requestObj.Email}
				err := collection.FindOne(context.TODO(), filter).Decode(&bsonData)
				if err != nil {
					responseObj.Error = err
					responseObj.Status = false
					responseObj.Message = "Error while Fetching Data"
				} else {
					if bsonToStruct(bsonData).Password == requestObj.Password {
						data := bsonToStruct(bsonData)
						responseObj.Data = objs.SecureUserData{
							Name:        data.Name,
							ID:          data.ID,
							PictureData: data.PictureData,
							Email:       data.Email,
						}
						responseObj.Message = "Successfully Fetched User Data"
						responseObj.Status = true
					} else {
						responseObj.Status = false
						responseObj.Message = "Wrong Password"
					}
				}
			} else {
				responseObj.Status = false
				responseObj.Message = "User Does Not Exist"
			}
		}

	}
	return responseObj
}
