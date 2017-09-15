package database

import (
	mgo "gopkg.in/mgo.v2"
)

type Options struct {
	host string
	name string
}

type Option func(*Options)

func Host(host string) Option {
	return func(o *Options) {
		o.host = host
	}
}

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
