package main

import (
	"fmt"
	"time"
)

func main() {
	// viper.SetEnvPrefix("grpc")
	// viper.SetTypeByDefaultValue(true)
	// viper.SetDefault("name", "hello")
	// viper.BindEnv("name")
	// viper.BindEnv("languages")

	// type Config struct {
	// 	Name      string
	// 	Languages []string
	// }

	// var c Config
	// err := viper.Unmarshal(&c)
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println(c)

	// log.Println(viper.Get("name"))
	// log.Printf("%#v", viper.GetStringSlice("languages"))

	fmt.Println(time.Now().UTC().Format(time.RFC1123Z))
	fmt.Println(time.Now().UTC().Format(time.RFC3339))
	fmt.Println(time.Now().UTC().Format(time.RFC822))
	fmt.Println(time.Now().UTC().Format(time.RFC850))
	fmt.Print(time.Now().Format("2006-01-02 15:04:05 -0700 MST"))
}
