package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
	"time"

	"distcache-go/lrucache"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var (
	db    *sql.DB
	cache = lrucache.New(100)

	dbUser             = os.Getenv("DB_USER")
	dbPass             = os.Getenv("DB_PASS")
	dbHost             = os.Getenv("DB_HOST")
	dbPort             = os.Getenv("DB_PORT")
	dbName             = os.Getenv("DB_NAME")
	nodesMap, nodesLen = mustCreateNodeMap(strings.Split(os.Getenv("NODES"), ";"))
	self               = os.Getenv("SELF")
)

// Message ...
type Message struct {
	ID   int64  `json:"id"`
	Text string `json:"text"`
}

// ListResponse ...
type ListResponse struct {
	Messages []*Message `json:"messages"`
}

// CreateMessageRequest ...
type CreateMessageRequest struct {
	Text string `json:"text"`
}

func main() {
	dbPort, err := strconv.Atoi(dbPort)
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		text TEXT NOT NULL
	)`)
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/messages", ListMessagesHandler).Methods(http.MethodGet)
	r.HandleFunc("/messages", CreateMessageHandler).Methods(http.MethodPost)
	r.HandleFunc("/messages/{id:[0-9]+}", GetMessageHandler).Methods(http.MethodGet)
	srv := &http.Server{
		Handler: r,
		Addr:    ":8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("listening on localhost port 8080")

	log.Fatal(srv.ListenAndServe())
}

// ListMessagesHandler ...
func ListMessagesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("request received")

	w.Header().Set("Content-Type", "application/json")
	val := cache.Get(r.URL.Path)
	if val != nil {
		w.Write(val.([]byte))
		w.WriteHeader(http.StatusOK)
		return
	}

	db.Query(`SELECT pg_sleep(5);`)

	rows, err := db.Query(`SELECT id, text FROM messages;`)
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
			return
		}

		messages = append(messages, message)
	}

	JSON, err := json.Marshal(&ListResponse{
		Messages: messages,
	})
	if err != nil {
		log.Printf("Error: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cache.Set(r.URL.Path, JSON)

	w.WriteHeader(http.StatusOK)
	w.Write(JSON)
}

// CreateMessageHandler ...
func CreateMessageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("request received")
	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()

	var createReq *CreateMessageRequest
	err := json.NewDecoder(r.Body).Decode(&createReq)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hashKey := len(createReq.Text) % nodesLen
	addr := nodesMap[hashKey]

	if addr != self {
		message, err := redirectCreateRequest(addr, createReq)
		if err != nil {
			handleErr(w, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&message)
		return
	}

	message := &Message{
		Text: createReq.Text,
	}
	err = db.QueryRow(`INSERT INTO messages (text) VALUES ($1) RETURNING id;`, message.Text).Scan(&message.ID)
	if err != nil {
		handleErr(w, err)
		return
	}

	cache.Invalidate(r.URL.Path)

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(&message)
	if err != nil {
		handleErr(w, err)
		return
	}
}

// GetMessageHandler ...
func GetMessageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	val := cache.Get(r.URL.Path)
	if val != nil {
		w.Write(val.([]byte))
		w.WriteHeader(http.StatusOK)
		return
	}

	db.Query(`SELECT pg_sleep(5);`)

	message := &Message{}
	err := db.
		QueryRow(`SELECT id, text FROM messages WHERE id = $1;`, id).
		Scan(
			&message.ID,
			&message.Text,
		)
	if err != nil {
		handleErr(w, err)
		return
	}

	JSON, err := json.Marshal(&message)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	cache.Set(r.URL.Path, JSON)

	w.WriteHeader(http.StatusOK)
	w.Write(JSON)
}

func handleErr(w http.ResponseWriter, err error) {
	log.Printf("Error: %s", err.Error())
	w.WriteHeader(http.StatusInternalServerError)
}

func mustCreateNodeMap(nodes []string) (map[int]string, int) {
	nodeMap := make(map[int]string)
	numNodes := 0

	for _, nodeAddr := range nodes {
		nodeSlice := strings.Split(nodeAddr, "::")
		mapKey, err := strconv.Atoi(nodeSlice[0])
		if err != nil {
			panic(err)
		}

		nodeMap[mapKey] = nodeSlice[1]
		numNodes++
	}

	return nodeMap, numNodes
}

func redirectCreateRequest(addr string, r *CreateMessageRequest) (*Message, error) {
	JSON, err := json.Marshal(&r)
	if err != nil {
		return nil, err
	}

	reader := strings.NewReader(string(JSON))

	resp, err := http.Post("http://"+addr+"/messages", "application/json", reader)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(dump))

	var message *Message
	err = json.NewDecoder(resp.Body).Decode(&message)
	if err != nil {
		return nil, err
	}

	return message, nil
}
