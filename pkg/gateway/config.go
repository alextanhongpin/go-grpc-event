package main

import "github.com/spf13/viper"

func init() {
	// Gateway
	viper.BindEnv("port") // TCP address to listen to

	// gRPC Server
	viper.BindEnv("photo_addr") // Address of the photo gRPC server
	viper.BindEnv("event_addr") // Address of the event gRPC server

	// Auth0
	viper.BindEnv("auth0_jwk")       // auth0 jwks uri
	viper.BindEnv("auth0_iss")       // auth0 jwks issuer
	viper.BindEnv("auth0_aud")       // auth0 jwks audience
	viper.BindEnv("auth0_whitelist") // auth0 whitelisted admin emails

	// Tracer
	viper.BindEnv("tracer") // opentracing tracer namespace
	viper.BindEnv("tracer_sampler_url")
	viper.BindEnv("tracer_reporter_url")

	// Defaults
	viper.SetDefault("port", ":3000")
	viper.SetDefault("tracer", "gateway")
	viper.SetDefault("tracer_sampler_url", "localhost:5775")
	viper.SetDefault("tracer_reporter_url", "localhost:6831")
}
