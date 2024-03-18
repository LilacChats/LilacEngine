package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"objs"
	"routes"
)

func main() {
	client, err := InitializeConnection(objs.MONGO_URL)
	if err != nil {
		fmt.Println(err)
	} else {
		mux := http.NewServeMux()
		mux.HandleFunc("/usersignup", func(res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(routes.SignupHandler(req, client))
			res.Write(response)
		})
		mux.HandleFunc("/login", func(res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(routes.LoginHandler(req, client))
			res.Write(response)
		})
		mux.HandleFunc("/creategroup", func(res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(routes.CreateGroupHandler(req, client))
			res.Write(response)
		})
		mux.HandleFunc("/fetchgroups", func(res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(routes.FetchGroupsHandler(req, client))
			res.Write(response)
		})
		mux.HandleFunc("/fetchusers", func(res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(routes.FetchUsersHandler(req, client))
			res.Write(response)
		})
		mux.HandleFunc("/deletegroup", func(res http.ResponseWriter, req *http.Request) {
			res.Header().Set("content-Type", "application/json")
			response, _ := json.Marshal(routes.DeleteGroupHandler(req, client))
			res.Write(response)
		})
		mux.HandleFunc("/fetchgroupdata", func(res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(routes.FetchGroupDataHandler(req, client))
			res.Write(response)
		})
		mux.HandleFunc("/updategroupdata", func(res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(routes.UpdateGroupHandler(req, client))
			res.Write(response)
		})
		err := http.ListenAndServe(":3000", mux)
		if err != nil {
			fmt.Println("Failed to started server!\n", err)
		}
	}
}
