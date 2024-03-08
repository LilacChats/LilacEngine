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

type SignupRequest struct {
	Name        string
	Email       string
	Password    string
	PictureData string
}

type CreateGroupRequest struct {
	GroupName string
	UserID    string
}

type CreateGroupResponse struct {
	Data struct{ GroupID string }
	StandardResponse
}

type FetchGroupsRequest struct {
	UserID string
}

type FetchGroupsResponse struct {
	Data []struct {
		ID   string `bson:"_id"`
		Name string `bson:"name"`
	}
	StandardResponse
}

type FetchUsersRequest struct {
	UserID string
}

type FetchUsersResponse struct {
	Data []struct {
		ID          string `bson:"_id"`
		Name        string `bson:"name"`
		PictureData string `bson:"picturedata"`
	}
	StandardResponse
}

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

var GroupList_DB = MongoDBObj{
	Collection: "Groups",
	Database:   DATABASE,
}

var DB_CHOICE = "Mongo"
