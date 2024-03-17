package objs

type UserData struct {
	ID          string `bson:"_id"`
	Name        string `bson:"name"`
	Email       string `bson:"email"`
	Password    string `bson:"password"`
	PictureData string `bson:"picturedata"`
}

type GroupData struct {
	ID      string   `bson:"_id"`
	Name    string   `bson:"name"`
	Members []string `bson:"members"`
}

type SecureUserData struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PictureData string `json:"pictureData"`
	ID          string `json:"id"`
}

type StandardResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Error   error  `json:"error"`
}

type SignupRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	PictureData string `json:"pictureData"`
}

type UpdateGroupRequest struct {
	Name    string   `json:"name"`
	UserID  string   `json:"userID"`
	GroupID string   `json:"groupID"`
	Members []string `json:"members"`
}

type UpdateGroupResponse struct {
	Data struct{} `json:"data"`
	StandardResponse
}

type FetchGroupDataRequest struct {
	UserID  string `bson:"userid"`
	GroupID string `bson:"groupid"`
}

type FetchGroupDataResponse struct {
	Data struct {
		ID      string   `json:"id"`
		Name    string   `json:"name"`
		Members []string `json:"members"`
	} `json:"data"`
	StandardResponse
}

type DeleteGroupRequest struct {
	UserID  string `json:"userID"`
	GroupID string `json:"groupID"`
}

type DeleteGroupResponse struct {
	Data struct{} `json:"data"`
	StandardResponse
}

type CreateGroupRequest struct {
	GroupName string   `json:"groupName"`
	UserID    string   `json:"userID"`
	Members   []string `json:"members"`
}

type CreateGroupResponse struct {
	Data struct {
		GroupID string `json:"groupID"`
	} `json:"data"`
	StandardResponse
}

type FetchGroupsRequest struct {
	UserID string `json:"userID"`
}

type FetchGroupsResponse struct {
	Data []struct {
		ID      string   `json:"id"`
		Name    string   `json:"name"`
		Members []string `json:"members"`
	} `json:"data"`
	StandardResponse
}

type FetchUsersRequest struct {
	UserID string `json:"userID"`
}

type FetchUsersResponse struct {
	Data []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		PictureData string `json:"pictureData"`
	} `json:"data"`
	StandardResponse
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupResponse struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
	StandardResponse
}
type LoginResponse struct {
	Data SecureUserData `json:"data"`
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
