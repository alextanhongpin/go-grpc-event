package database

import (
	mgo "gopkg.in/mgo.v2"
)

type Options struct {
	Host string
	Name string
}

type Option func(*Options)

func Host(host string) Option {
	return func(o *Options) {
		o.Host = host
	}
}

func Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

// New returns a new database
func New(opts ...Option) (*Database, error) {
	options := Options{
		Host: "localhost",
		Name: "go-engineersmy-event",
	}
	for _, o := range opts {
		o(&options)
	}

	sess, err := mgo.Dial(options.Host)
	if err != nil {
		return nil, err
	}

	return &Database{
		Ref:  sess,
		Name: options.Name,
	}, nil
}

type Database struct {
	Ref  *mgo.Session
	Name string
}

func (db Database) Close() {
	db.Ref.Close()
}

func (db Database) Copy() *mgo.Session {
	return db.Ref.Copy()
}

func (db Database) Collection(sess *mgo.Session, name string) *mgo.Collection {
	return sess.DB(db.Name).C(name)
}
