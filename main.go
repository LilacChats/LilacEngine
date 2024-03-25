package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"objs"
	"os"
	"routes"
	"strconv"
	"time"
	"validation"

	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/mattn/go-tty"
	"go.mongodb.org/mongo-driver/mongo"
)

var upgrader = websocket.Upgrader{}

func GenerateClientConnection() (*mongo.Client, error) {
	client, err := InitializeConnection(objs.MONGO_URL)
	if err != nil {
		return nil, err
	} else {
		return client, nil
	}
}

func ReRenderDBSelection(selectedOption int) {
	fmt.Print("\033[3;0H")
	var selector string
	if selectedOption == 0 {
		selector = "> "
	} else {
		selector = "  "
	}
	if selectedOption == 0 {
		fmt.Println(selector, objs.SelectedItemStyle.Render(objs.DBS[0]))
	} else {
		fmt.Println(selector, objs.DBS[0])
	}
	if selectedOption == 1 {
		selector = "> "
	} else {
		selector = "  "
	}
	if selectedOption == 1 {
		fmt.Println(selector, objs.SelectedItemStyle.Render(objs.DBS[1]))
	} else {
		fmt.Println(selector, objs.DBS[1])
	}
	fmt.Print("---------------\nPress CTRL+C to Exit")
}

func StartServer(client *mongo.Client, URL string) {
	clients := []string{}
	mux := http.NewServeMux()
	mux.HandleFunc("/sendmessage", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(routes.SendMessageHandler(req, client))
		res.Write(response)
	})
	mux.HandleFunc("/connect", func(res http.ResponseWriter, req *http.Request) {
		conn, err := upgrader.Upgrade(res, req, nil)
		packet := struct {
			Type string `json:"type"`
			ID1  string `json:"id1"`
			ID2  string `json:"id2"`
		}{}

		if err == nil {
			defer conn.Close()
			for {
				messageType, message, _ := conn.ReadMessage()
				err := json.Unmarshal(message, &packet)
				if err == nil {
					checker := false
					for i := 0; i < len(clients); i++ {
						if clients[i] == packet.ID1 {
							checker = true
							break
						}
					}
					if checker == false {
						clients = append(clients, packet.ID1)
					}
					messageData := routes.FetchAllChatsHandler(packet.ID1, packet.ID2, client)
					jsonData, err := json.Marshal(messageData)
					if err == nil {
						conn.WriteMessage(messageType, []byte(jsonData))
					}
				} else {
					conn.WriteMessage(messageType, []byte(err.Error()))
				}
			}
		}
	})
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
	mux.HandleFunc("/logout", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(routes.LogoutHandler(req, client))
		res.Write(response)
	})
	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		fmt.Println("Failed to started server!\n", err)
	}
}

func Animate(status string) {
	fmt.Print("\033[H\033[2J")
	fmt.Println("Database:", "\t", lipgloss.NewStyle().Underline(true).Foreground(lipgloss.Color("#4A89FD")).Render(objs.MONGO_URL))
	fmt.Println("Server URL:", "\t", lipgloss.NewStyle().Underline(true).Foreground(lipgloss.Color("#4A89FD")).Render("http://localhost:3000"))
	fmt.Print("---------------------\n ", "Server Status ")
	for {
		fmt.Print("    ")
		time.Sleep(1 * time.Second)
		fmt.Print("\033[4;15H")
		fmt.Print(" 🟢   ")
		time.Sleep(1 * time.Second)
		fmt.Print("\033[4;15H")
	}
}

func FetchDatabaseURL() string {
	ttyObj, _ := tty.Open()
	fmt.Print("\033[H\033[2J")
	fmt.Print("Enter Database URL:")
	input, _ := ttyObj.ReadString()
	ttyObj.Close()
	return input
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
	dbType := os.Getenv("DB")
	if dbType == "Mongo" {
		objs.MONGO_URL = os.Getenv("URL")
	}
	args := os.Args
	args = args[1:]
	switch args[0] {
	case "run":
		client, err := GenerateClientConnection()
		if err == nil {
			go StartServer(client, objs.MONGO_URL)
			Animate("ewfn")
		} else {
			validation.DisplayError("Connection to MongoDB Failed")
		}
	case "connect":
		fmt.Print("\033[H\033[2J")
		ttyObj, err := tty.Open()
		if err != nil {
			log.Fatal(err)
		}
		sequence := ""
		selection := 0
		fmt.Println("Select DB\n---------------")
		ReRenderDBSelection(selection)
		for i := 0; i >= 0; i++ {
			r, err := ttyObj.ReadRune()
			if err != nil {
				log.Fatal(err)
			} else {
				token := strconv.QuoteRuneToASCII(r)
				sequence += token[1 : len(token)-1]
				if sequence == "\\x00\\x1b[B\\x00" || sequence == "\\x1b[B\\x00" {
					// Pressed Down
					if selection == 0 {
						selection = 1
					}
					ReRenderDBSelection(selection)
					sequence = ""
				} else if sequence == "\\x00\\x1b[A\\x00" || sequence == "\\x1b[A\\x00" {
					// Pressed Up
					if selection == 1 {
						selection = 0
					}
					ReRenderDBSelection(selection)
					sequence = ""
				} else if sequence == "\\x00\\r\\x00" || sequence == "\\r\\x00" {
					break
				}
			}
		}
		ttyObj.Close()

		fmt.Println("\033[H\033[2J")

		if selection == 0 {
			objs.MONGO_URL = FetchDatabaseURL()
			fmt.Println("Database URL Set -> ", objs.MONGO_URL)
			dbType := ""
			if selection == 0 {
				dbType = "Mongo"
			} else {
				dbType = "MySQL"
			}
			os.WriteFile(".env", []byte("URL="+objs.MONGO_URL+"\nDB="+dbType), 0666)
		}
	}
}
