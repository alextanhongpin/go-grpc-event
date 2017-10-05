package main

import (
	"log"

	"github.com/spf13/viper"
)

func main() {
	viper.SetEnvPrefix("grpc")
	viper.SetTypeByDefaultValue(true)
	viper.SetDefault("name", "hello")
	viper.BindEnv("name")
	viper.BindEnv("languages")

	type Config struct {
		Name      string
		Languages []string
	}

	var c Config
	err := viper.Unmarshal(&c)
	if err != nil {
		log.Println(err)
	}
	log.Println(c)

	log.Println(viper.Get("name"))
	log.Printf("%#v", viper.GetStringSlice("languages"))
}
