package database

import mgo "gopkg.in/mgo.v2"

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
