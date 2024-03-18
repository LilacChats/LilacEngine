package routes

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"objs"
	"validation"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LoginHandlers struct{}

func (LoginHandlers) Mongo(data objs.LoginRequest, client *mongo.Client) (objs.SecureUserData, error) {
	var bsonData bson.M
	dataObject := objs.UserData{}
	userData := objs.SecureUserData{}
	collection := client.Database(objs.UserData_DB.Database).Collection(objs.UserData_DB.Collection)
	err := collection.FindOne(context.TODO(), bson.M{"email": data.Email}).Decode(&bsonData)
	if err != nil {
		return userData, err
	} else {
		byteData, err := bson.Marshal(bsonData)
		if err != nil {
			return userData, err
		} else {
			bson.Unmarshal(byteData, &dataObject)
			if dataObject.Password == data.Password {
				return objs.SecureUserData{
					ID:          dataObject.ID,
					Email:       dataObject.Email,
					Name:        dataObject.Name,
					PictureData: dataObject.PictureData,
				}, nil
			} else {
				return userData, errors.New("password mismatch")
			}
		}
	}
}

func LoginHandler(req *http.Request, client *mongo.Client) objs.LoginResponse {
	var requestObj = objs.LoginRequest{}
	var responseObj = objs.LoginResponse{}
	json.NewDecoder(req.Body).Decode(&requestObj)
	var loginHandlers objs.LoginMethods = LoginHandlers{}
	responseObj.Error = nil
	switch req.Method {
	case "GET":
		responseObj.Message = "Bad route config"
		responseObj.Status = false
	case "POST":
		switch objs.DB_CHOICE {
		case "Mongo":
			var mongoValidationHandlers objs.MongoValidationMethods = validation.MongoValidationHandlers{}
			validationObj := mongoValidationHandlers.ValidateUserEmail(requestObj.Email, client)
			if validationObj.Status {
				loginData, err := loginHandlers.Mongo(requestObj, client)
				if err != nil {
					responseObj.Status = false
					responseObj.Message = err.Error()
				} else {
					responseObj.Data = loginData
					responseObj.Status = true
					responseObj.Message = "Successfully Fetched User Data"
				}
			} else {
				responseObj.Status = false
				responseObj.Message = validationObj.Message
			}
		}

	}
	return responseObj
}
