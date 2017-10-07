package database

import (
	mgo "gopkg.in/mgo.v2"
)

// DB represents the database config
type DB struct {
	Ref  *mgo.Session
	Name string
}

// Close will close the connection to the database
func (db DB) Close() {
	db.Ref.Close()
}

// Copy will create a new database session
func (db DB) Copy() *mgo.Session {
	return db.Ref.Copy()
}

// Collection returns the a new collection by name
func (db DB) Collection(sess *mgo.Session, name string) *mgo.Collection {
	return sess.DB(db.Name).C(name)
}
