package objs

type UserData struct {
	Name        string
	Email       string
	Password    string
	PictureData string
	ID          string
}

type SecureUserData struct {
	Name        string
	Email       string
	PictureData string
	ID          string
}

type StandardResponse struct {
	Status  bool
	Message string
	Error   error
}

type SignupRequest struct{ UserData }
type LoginRequest struct {
	Email    string
	Password string
}

type SignupResponse struct {
	Data struct{ ID string }
	StandardResponse
}
type LoginResponse struct {
	Data SecureUserData
	StandardResponse
}

type MongoDBObj struct {
	Collection string
	Database   string
}

var DATABASE string = "Lilac"

var UserData_DB = MongoDBObj{
	Collection: "UserData",
	Database:   DATABASE,
}

var DB_CHOICE = "Mongo"
