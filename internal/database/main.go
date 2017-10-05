package database

import (
	"log"
	"time"

	mgo "gopkg.in/mgo.v2"
)

// New returns a new database
func New(opts ...Option) (*Database, error) {
	options := Options{
		host:     "localhost",
		db:       "engineersmy",
		username: "",
		password: "",
	}
	for _, o := range opts {
		o(&options)
	}

	log.Println("connecting to db")
	sess, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{options.host}, // options.host},
		Timeout:  10 * time.Second,
		Database: options.db,
		Username: options.username,
		Password: options.password,
	})

	if err != nil {
		return nil, err
	}
	log.Println("connected to db")

	return &Database{
		Ref:  sess,
		Name: options.db,
	}, nil
}
