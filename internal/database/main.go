package database

import (
	"log"
	"time"

	mgo "gopkg.in/mgo.v2"
)

// New returns a new database
func New(opts ...Option) (*DB, error) {
	options := Options{
		host:     "localhost",
		db:       "engineersmy",
		username: "",
		password: "",
		timeout:  60,
	}
	for _, o := range opts {
		o(&options)
	}

	log.Println("connecting to db")
	sess, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{options.host}, // options.host},
		Timeout:  time.Duration(options.timeout) * time.Second,
		Database: options.db,
		Username: options.username,
		Password: options.password,
	})

	if err != nil {
		return nil, err
	}
	log.Println("connected to db")

	return &DB{
		Ref:  sess,
		Name: options.db,
	}, nil
}
