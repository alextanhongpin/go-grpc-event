package main

import (
	"log"

	"github.com/spf13/viper"
)

func init() {
	log.Println("calling init")
	viper.BindEnv("port")
	viper.BindEnv("mgo_host")
	viper.BindEnv("mgo_db")
	viper.BindEnv("mgo_usr")
	viper.BindEnv("mgo_pwd")
	viper.BindEnv("tracer")
	viper.BindEnv("tracer_sampler_url")
	viper.BindEnv("tracer_reporter_url")

	viper.BindEnv("slack_channel")
	viper.BindEnv("slack_username")
	viper.BindEnv("slack_icon")
	viper.BindEnv("slack_webhook")

	viper.SetDefault("port", ":8080")                         // TCP port to listen to
	viper.SetDefault("mgo_host", "mongodb://localhost:27017") // mongoDB URI connection
	viper.SetDefault("mgo_db", "engineersmy")                 // mongoDB database name
	viper.SetDefault("tracer", "event_service")               // The namespace of the opentracing tracer
	viper.SetDefault("tracer_sampler_url", "localhost:5775")
	viper.SetDefault("tracer_reporter_url", "localhost:6831")
}
