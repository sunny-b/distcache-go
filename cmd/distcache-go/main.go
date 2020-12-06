package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"distcache-go/lrucache"

	"github.com/gorilla/mux"
)

var (
	db    *sql.DB
	cache = lrucache.New(10)

	dbUser = os.Getenv("DB_USER")
	dbPass = os.Getenv("DB_PASS")
	dbHost = os.Getenv("DB_HOST")
	dbPort = os.Getenv("DB_PORT")
	dbName = os.Getenv("DB_NAME")
)

// Message ...
type Message struct {
	ID   int
	Text string
}

// ListResponse ...
type ListResponse struct {
	Messages []*Message `json:"messages"`
}

func main() {
	dbPort, err := strconv.Atoi(dbPort)
	if err != nil {
		panic(err)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY,
		message TEXT NOT NULL,
	)`)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/messages", ListMessagesHandler).Methods(http.MethodGet)
	r.HandleFunc("/messages", CreateMessageHandler).Methods(http.MethodPost)
	r.HandleFunc("/messages/{id:[0-9]+}", GetMessageHandler).Methods(http.MethodGet)
	http.Handle("/", r)
}

// ListMessagesHandler ...
func ListMessagesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	val := cache.Get(r.URL.Path)
	if val != nil {
		w.Write([]byte(val.(string)))
		w.WriteHeader(http.StatusOK)
		return
	}

	rows, err := db.Query(`SELECT id, message FROM messages;`)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	messages := []*Message{}
	for rows.Next() {
		message := &Message{}
		err := rows.Scan(
			&message.ID,
			&message.Text,
		)
		if err != nil {
			log.Printf("Error: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}

		messages = append(messages, message)
	}

	JSON, err := json.Marshal(&ListResponse{
		Messages: messages,
	})
	if err != nil {
		log.Printf("Error: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(JSON)
	w.WriteHeader(http.StatusOK)
}

// CreateMessageHandler ...
func CreateMessageHandler(w http.ResponseWriter, r *http.Request) {
	return
}

// GetMessageHandler ...
func GetMessageHandler(w http.ResponseWriter, r *http.Request) {
	return
}
