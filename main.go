package main

import (
	"net/http"
	"github.com/go-redis/redis"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"fmt"
	"time"
	"os"
)

func emitBootBanner(port string)  {
	template := `
 ____  _  _  _  _  ____    ____  ____  __  ____  ____ 
(    \/ )( \( \/ )(  _ \  / ___)(_  _)/  \(  _ \(  __)
 ) D () \/ (/ \/ \ ) _ (  \___ \  )( (  O ))   / ) _) 
(____/\____/\_)(_/(____/  (____/ (__) \__/(__\_)(____)

[-] Starting on %s

`
	fmt.Printf(template, port)
}

func main()  {
	port := ":7575"

	redisHost := os.Getenv("REDIS_HOST")

	if redisHost == "" {
		log.Fatalln("App requires Redis host in env (REDIS_HOST)")
	}

	redisPass := os.Getenv("REDIS_PASSWORD")


	if redisPass == "" {
		log.Fatalln("App requires Redis password in env (REDIS_PASSWORD)")
	}

	storage := NewRedisPersister(redisHost, redisPass)

	http.Handle("/store", storeHandler(storage))
	emitBootBanner(port)
	http.ListenAndServe(port, nil)
}

// storeHandler decodes whatever it gets and puts it in storage
func storeHandler(storage Persister) http.HandlerFunc  {
	return func(w http.ResponseWriter, r *http.Request){
		var bodyMessage Message
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&bodyMessage)
		if err != nil {
			log.Info(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err = storage.Put(bodyMessage); err != nil {
			log.Info(err)
			w.WriteHeader(http.StatusBadGateway)
			return
		}

		log.Info(bodyMessage.Data)
		w.WriteHeader(http.StatusOK)
	}
}

type Message struct{
	Data string `json:"data"`
}

type Persister interface{
	Put(Message) error
}

type RedisPersist struct{
	client *redis.Client
}

func NewRedisPersister(c, p string) *RedisPersist {
	password := os.Getenv("REDIS_PASSWORD")
	if password == ""{
		log.Fatalln("App requires Redis password in env (REDIS_PASSWORD)")
	}
	return &RedisPersist{
		redis.NewClient(&redis.Options{
			Addr:     c,
			Password: p,
			DB:       0,  // use default DB
	}),
	}
}

// Put (as written) is useless! :-)
func (rp *RedisPersist) Put(m Message) error {
	return rp.client.Set(fmt.Sprintf("%d", time.Now().UTC().UnixNano()), m.Data, 0).Err()
}

