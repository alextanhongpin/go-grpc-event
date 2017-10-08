package main

import "github.com/spf13/viper"

func init() {
	viper.BindEnv("port")
	viper.BindEnv("mgo_host")
	viper.BindEnv("mgo_db")
	viper.BindEnv("mgo_usr")
	viper.BindEnv("mgo_pwd")
	viper.BindEnv("tracer")
	viper.BindEnv("tracer_sampler_url")
	viper.BindEnv("tracer_reporter_url")

	viper.SetDefault("port", ":5000")               // The port to listen to
	viper.SetDefault("mgo_host", "localhost:27017") // The mongo host to listen to
	viper.SetDefault("mgo_db", "engineersmy")       // The mongo database name
	viper.SetDefault("mgo_usr", "user")             // The default mongo username
	viper.SetDefault("mgo_pwd", "password")         // The default mongo password
	viper.SetDefault("tracer", "photo_service")
	viper.SetDefault("tracer_sampler_url", "localhost:5775")
	viper.SetDefault("tracer_reporter_url", "localhost:6831")
}
