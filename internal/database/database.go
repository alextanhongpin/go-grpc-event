package database

import (
	mgo "gopkg.in/mgo.v2"
)

// Options represents the available options for the database
type Options struct {
	host string
	name string
}

// Option is a function that returns a new option
type Option func(*Options)

// Host is a the database host
func Host(host string) Option {
	return func(o *Options) {
		o.host = host
	}
}

// Name is the database name
func Name(name string) Option {
	return func(o *Options) {
		o.name = name
	}
}

// New returns a new database
func New(opts ...Option) (*Database, error) {
	options := Options{
		host: "localhost",
		name: "go-engineersmy-event",
	}
	for _, o := range opts {
		o(&options)
	}

	sess, err := mgo.Dial(options.host)
	if err != nil {
		return nil, err
	}

	return &Database{
		Ref:  sess,
		Name: options.name,
	}, nil
}

// Database represents the database config
type Database struct {
	Ref  *mgo.Session
	Name string
}

// Close will close the connection to the database
func (db Database) Close() {
	db.Ref.Close()
}

// Copy will create a new database session
func (db Database) Copy() *mgo.Session {
	return db.Ref.Copy()
}

// Collection returns the a new collection by name
func (db Database) Collection(sess *mgo.Session, name string) *mgo.Collection {
	return sess.DB(db.Name).C(name)
}
