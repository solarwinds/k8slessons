package main

import (
	"net/http"
	"github.com/go-redis/redis"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"fmt"
	"time"
)

func main()  {
	port := ":7575"
	storage := NewRedisPersister("localhost:6379")

	http.Handle("/store", storeHandler(storage))
	fmt.Println("[-] Booting on ", port)
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

func NewRedisPersister(c string) *RedisPersist {
	return &RedisPersist{
		redis.NewClient(&redis.Options{
			Addr:     c,
			Password: "", // no password set
			DB:       0,  // use default DB
	}),
	}
}

// Put (as written) is useless! :-)
func (rp *RedisPersist) Put(m Message) error {
	return rp.client.Set(fmt.Sprintf("%d", time.Now().UTC().UnixNano()), m.Data, 0).Err()
}


