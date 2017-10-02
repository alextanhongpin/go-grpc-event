package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strings"

	// auth0 "github.com/auth0-community/go-auth0"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	"github.com/alextanhongpin/go-grpc-event/internal/auth0"
	"github.com/alextanhongpin/go-grpc-event/internal/cors"
	gw "github.com/alextanhongpin/go-grpc-event/proto/event-private"
)

// UserInfo represents the schema from the auth0 userinfo endpoint
// type UserInfo struct {
// 	EmailVerified bool   `json:"email_verified"` // false
// 	Email         string `json:"email"`          // "test.account@userinfo.com"
// 	ClientID      string `json:"clientID"`       // "q2hnj2iu..."
// 	UpdatedAt     string `json:"updated_at"`     // "2016-12-05T15:15:40.545Z"
// 	Name          string `json:"name"`           //  "test.account@userinfo.com"
// 	Picture       string `json:"picture"`        // "https://s.gravatar.com/avatar/dummy.png"
// 	UserID        string `json:"user_id"`        // "auth0|58454..."
// 	Nickname      string `json:"nickname"`       // "test.account"
// 	CreatedAt     string `json:"created_at"`     // "2016-12-05T11:16:59.640Z"
// 	Sub           string `json:"sub"`            // "auth0|58454..."
// }

// Response represents the payload that is returned on error
type Response struct {
	Message string `json:"message"`
}

var auth0Validator *auth0.Auth0
var whitelist []string

func run() error {
	var (
		addr              = flag.String("addr", "localhost:8081", "Address of the private event GRPC service")
		port              = flag.String("port", ":9081", "TCP address to listen on")
		jwksURI           = flag.String("jwks_uri", "", "Auth0 jwks uri available at auth0 dashboard")
		auth0APIIssuer    = flag.String("auth0_iss", "", "Auth0 api issuer available at auth0 dashboard")
		auth0APIAudience  = flag.String("auth0_aud", "", "Auth0 api audience available at auth0 dashboard")
		whitelistedEmails = flag.String("whitelisted_emails", "", "A list of admin emails that are whitelisted")
	)
	flag.Parse()

	if *whitelistedEmails != "" {
		tmp := strings.Split(*whitelistedEmails, ",")
		for _, v := range tmp {
			whitelist = append(whitelist, strings.TrimSpace(v))
		}
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	mux := runtime.NewServeMux()
	if err := gw.RegisterEventServiceHandlerFromEndpoint(ctx, mux, *addr, opts); err != nil {
		return err
	}

	// Create a global validator that is only initialized once
	auth0Validator = auth0.New(
		auth0.Audience(*auth0APIAudience),
		auth0.JWKSURI(*jwksURI),
		auth0.Issuer(*auth0APIIssuer),
	)

	log.Printf("grpc_server = %s\n", *addr)
	log.Printf("listening to port *%s. press ctrl + c to cancel.\n", *port)
	return http.ListenAndServe(*port, cors.New(checkJWT(mux)))
}

func main() {
	defer glog.Flush()
	if err := run(); err != nil {
		glog.Fatal(err)
	}
}

func checkJWT(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Is public user
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			h.ServeHTTP(w, r)
			return
		}

		_, err := auth0Validator.Validate(r)
		// log.Println("got token", token)

		if err != nil {
			log.Printf("Error validating token: %#v\n", err)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(Response{
				Message: "Missing or invalid token.",
			})
			return
		}

		// Fetch the user details and pass it through the grpc-metadata
		if err = fetchUserDetails(r); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{
				Message: "Unable to get users data.",
			})
			return
		}
		h.ServeHTTP(w, r)
	})
}

func fetchUserDetails(r *http.Request) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://engineersmy.auth0.com/userinfo", nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", r.Header.Get("Authorization"))
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	userinfo := make(map[string]string) // UserInfo

	if err = json.NewDecoder(resp.Body).Decode(&userinfo); err != nil {
		return err
	}
	// For each of the users metadata present, write it to the grpc-metadata
	for k, v := range userinfo {
		var buff bytes.Buffer
		buff.WriteString("Grpc-Metadata-")
		buff.WriteString(k)
		r.Header.Set(buff.String(), v)
	}
	if email, ok := userinfo["email"]; ok {
		if len(whitelist) > 0 {
			for _, v := range whitelist {
				if v == email {
					r.Header.Set("Grpc-Metadata-Admin", "true")
				}
			}
		}
	}

	// Do a checking to see the user emails which are whitelisted
	return nil
}
