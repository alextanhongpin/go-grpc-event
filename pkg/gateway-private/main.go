package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	auth0 "github.com/auth0-community/go-auth0"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	jose "gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"

	"github.com/alextanhongpin/go-grpc-event/internal/cors"
	gw "github.com/alextanhongpin/go-grpc-event/proto/event-private"
)

var (
	gauth0APIAudience string
	gjwksURI          string
	gauth0APIIssuer   string
)

func run() error {
	var (
		addr             = flag.String("addr", "localhost:8081", "Address of the private event GRPC service")
		port             = flag.String("port", ":9081", "TCP address to listen on")
		jwksURI          = flag.String("jwks_uri", "", "Auth0 jwks uri available at auth0 dashboard")
		auth0APIIssuer   = flag.String("auth0_iss", "", "Auth0 api issuer available at auth0 dashboard")
		auth0APIAudience = flag.String("auth0_aud", "", "Auth0 api audience available at auth0 dashboard")
	)
	flag.Parse()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(GetUnaryClientInterceptor())),
	}

	mux := runtime.NewServeMux()
	if err := gw.RegisterEventServiceHandlerFromEndpoint(ctx, mux, *addr, opts); err != nil {
		return err
	}
	gjwksURI = *jwksURI
	gauth0APIIssuer = *auth0APIIssuer
	gauth0APIAudience = *auth0APIAudience

	log.Printf("listening to service=private_event at endpoint=%s\n", *addr)
	log.Printf("listening to port *%s\n", *port)
	return http.ListenAndServe(*port, cors.New(checkJWT(mux)))
}

func main() {
	defer glog.Flush()
	if err := run(); err != nil {
		glog.Fatal(err)
	}
}

type Response struct {
	Message string
}

func GetUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		methodName := fmt.Sprintf("client:%s", method)
		log.Println(methodName)
		err := invoker(ctx, method, req, reply, cc, opts...)
		return err
	}
}

func checkJWT(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/public/v1/events" {
			h.ServeHTTP(w, r)
			return
		}
		client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: gjwksURI})
		audience := []string{gauth0APIAudience}

		configuration := auth0.NewConfiguration(client, audience, gauth0APIIssuer, jose.RS256)
		validator := auth0.NewValidator(configuration)

		token, err := validator.ValidateRequest(r)

		if err != nil {
			fmt.Println("Token is not valid or missing token", err)

			response := Response{
				Message: "Missing or invalid token.",
			}

			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)

		} else {
			// Ensure the token has the correct scope
			result := checkScope(r, validator, token)
			if result == true {
				h.ServeHTTP(w, r)
			} else {
				response := Response{
					Message: "You do not have the read:events scope.",
				}
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(response)
			}
		}
	})
}

func checkScope(r *http.Request, validator *auth0.JWTValidator, token *jwt.JSONWebToken) bool {

	log.Printf("Checking token scope: %#v", token)
	claims := map[string]interface{}{}
	err := validator.Claims(r, token, &claims)

	if err != nil {
		fmt.Println(err)
		return false
	}
	// if strings.Contains(claims["scope"].(string), "read:events") {
	// 	return true
	// }
	return true
}
