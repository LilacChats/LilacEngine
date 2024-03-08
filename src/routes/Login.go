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

func LoginUser(req *http.Request, client *mongo.Client) objs.LoginResponse {
	var requestObj = objs.LoginRequest{}
	var responseObj = objs.LoginResponse{}
	var bsonData bson.D
	json.NewDecoder(req.Body).Decode(&requestObj)
	var dataObject = struct {
		ID          string `bson:"_id"`
		Name        string `bson:"name"`
		Email       string `bson:"email`
		Password    string `bson:"password"`
		PictureData string `bson:"picturedata"`
	}{}
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
				filter := bson.M{"email": requestObj.Email}
				err := collection.FindOne(context.TODO(), filter).Decode(&bsonData)
				if err != nil {
					responseObj.Error = err
					responseObj.Status = false
					responseObj.Message = "Error while Fetching Data"
				} else {
					byteData, _ := bson.Marshal(bsonData)
					bson.Unmarshal(byteData, &dataObject)
					if dataObject.Password == requestObj.Password {
						responseObj.Data = objs.SecureUserData{
							Name:        dataObject.Name,
							ID:          dataObject.ID,
							PictureData: dataObject.PictureData,
							Email:       dataObject.Email,
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
