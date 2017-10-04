package main

import (
	"log"

	"github.com/spf13/viper"
)

func init() {
	log.Println("calling init")
	viper.BindEnv("port")
	viper.BindEnv("mgo_host")
	viper.BindEnv("tracer")

	viper.SetDefault("port", ":8080")                         // TCP port to listen to
	viper.SetDefault("mgo_host", "mongodb://localhost:27017") // mongoDB URI connection
	viper.SetDefault("tracer", "event_service")               // The namespace of the opentracing tracer
}
