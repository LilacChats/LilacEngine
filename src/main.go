package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"routes"
)

func main() {
	client, err := InitializeConnection("mongodb://localhost:27017")
	if err != nil {
		fmt.Println(err)
	} else {
		mux := http.NewServeMux()
		mux.HandleFunc("/usersignup", func(res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(routes.SignupUser(req, client))
			res.Write(response)
		})
		mux.HandleFunc("/login", func(res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json")
			response, _ := json.Marshal(routes.LoginUser(req, client))
			res.Write(response)
		})
		err := http.ListenAndServe(":3000", mux)
		if err != nil {
			fmt.Println("Failed to started server!\n", err)
		}
	}
}
