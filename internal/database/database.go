package database

import (
	mgo "gopkg.in/mgo.v2"
)

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
