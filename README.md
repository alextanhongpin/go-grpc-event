## Example of running separate makefile

```
$ touch Build.mk
$ make -f Build.mk scream
```

## Example of adding cors

```go
// allowCORS allows Cross Origin Resoruce Sharing from any origin.
// Don't do this without consideration in production systems.
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	glog.Infof("preflight request for %s", r.URL.Path)
	return
}
```

## Example of using Cors library (did not work though)

```go
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowedHeaders:   []string{"Authorization", "Access-Control-Allow-Headers", "Origin", "Accept", "X-Requested-With", "Content-Type", "Access-Control-Request-Method", "Access-Control-Request-Headers"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	}).Handler(mux)
```

## Example of passing value through http context

```go
// ...
	ctx := context.WithValue(r.Context(), "authorized", false)
	h.ServeHTTP(w, r.WithContext(ctx))
```

##  Example of getting metadata from context

```go
		md, ok := metadata.FromContext(ctx)
		log.Printf("md: %#v, ok: %v", md, ok)
```

## Example of setting metadata from context

```go
// create new context with metadata
md := metadata.Pairs("authorization", "Bearer XXXX")
ctx := metadata.NewContext(context.Background(), md)

something, err := client.SomeRPCCall(ctx, req)
```


## Example of creating a unary client interceptor

```go
func main () {
	// Not shown due to brevity
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(GetUnaryClientInterceptor())),
	}
	// Not shown due to brevity
}
func GetUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		methodName := fmt.Sprintf("client:%s", method)
		log.Println(methodName)
		err := invoker(ctx, method, req, reply, cc, opts...)
		return err
	}
}
```

## Passing metadata from middleware to client interceptor at Grpc Gateway

When running the code, it will go through the middleware first before going through the client interceptor. Since the client interceptor cannot access the request context (authorization, jwt token), you want to get the context from the middleware first, and then pass it to the client interceptor.

```go
func run () error {
	// Not shown here
	return http.ListenAndServe(*port, authMiddleware(handler))
}

func authMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Attempt 1: Passing it through context. Unfortunately it does not pass through...
		var key Key = "Grpc-Metadata-new"
		var key2 Key = "john"
		ctx := context.WithValue(r.Context(), key, "Setting!")
		ctx = context.WithValue(ctx, key2, "key2")

		// Attempt 2: Trying to read the metadata from the context. Unfortunately it doesn't work too...
		md, ok := metadata.FromIncomingContext(ctx)
		log.Println("ok", ok)
		if ok {
			log.Printf("Got auth metadata: %#v\n", md)
		}

		// Attempt 3: How about sending them through the header? Nope, no effect
		header := metadata.Pairs("header-key", "val")
		grpc.SendHeader(ctx, header)

		// Attempt 4: Trailer doesn't work either
		trailer := metadata.Pairs("trailer-key", "val")
		grpc.SetTrailer(ctx, trailer)

		// Attempt 5: How about setting it through the metadata? False hope...
		md = metadata.Pairs(
			"auth-middleware", "sending metadata from auth middleware",
			"Grpc-Metadata-cc", "cc lo",
		)
		ctx = metadata.NewOutgoingContext(ctx, md)

		// Attempt 6: Final attempts! Set the header with the Grpc-Metadata-<fieldname>: value. It works!
		r.Header.Set("Grpc-Metadata-testing", "it works!")

		h.ServeHTTP(w, r.WithContext(ctx))
		return
	})
}

// From the unary client interceptor, you can pass metadata to the server side. But how about reading them?

func GetUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		// Attempt 1: Try reading the metadata from the incoming context. Doesn't work...
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			log.Printf("got metadata: %#v\n", md)
		}

		// Attempt 2: How about reading them from the context itself? Nah, fail
		md2, ok2 := metadata.FromContext(ctx)
		log.Println("ok2", ok2)
		if ok2 {
			log.Println(md2)
		}

		// Attempt 3: Turns out you can actually read it from outgoing context...just log the ctx
		// Note that this metadata also receives the `Grpc-Metadata-<field>` set from the headers in
		// a curl request
		// curl -H "Grpc-Metadata-Example: test" http://localhost:3100/v1/events
		md, ok := metadata.FromOutgoingContext(ctx)
		if ok {
			log.Println("got value!", md)
		}

		// Attempt 4: We can also pass metadata to the server! Note that we don't need to prefix the field with `Grpc-Metadata-<fieldname>` as it the metadata will automatically do it for us
		md = metadata.Pairs(
			"hahahaha", "awesome! this passes~",
			"Grpc-Metadata-anotherone", "will this pass?", // Prefix not required
		)
		ctx = metadata.NewOutgoingContext(ctx, md)
		err := invoker(ctx, method, req, reply, cc, opts...)
		return err
	}
}
```

## Handling auth at Grpc Server

You can implement grpc server interceptor and also auth
```go

import "github.com/grpc-ecosystem/go-grpc-middleware/auth"
func main () {
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_opentracing.StreamServerInterceptor(tracerOpts...),
			grpc_auth.StreamServerInterceptor(AuthFunc),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_opentracing.UnaryServerInterceptor(tracerOpts...),
			SomeInterceptor(),
			grpc_auth.UnaryServerInterceptor(AuthFunc),
		)),
	)
}

// Example of interceptor on the server side with grpc
func SomeInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Println("unary auth interceptor", req, info)
		// log.Println("in UnaryServerInterceptor")
		// log.Println(ctx)
		// // Note that this metadata also receives the `Grpc-Metadata-<field>` set from the headers in
		// // a curl request
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			log.Println("Got metadata", md)
		}

		return handler(ctx, req)
	}
}

// Example of auth. This however requires the bearer token to be present - in our use case where
// we have both public and private users, this can be limiting
func AuthFunc(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}
	log.Println("token", token)
	// Parse token
	// tokenInfo, err := parseToken(token)
	grpc_ctxtags.Extract(ctx).Set("auth.sub", "something")
	newCtx := context.WithValue(ctx, "tokenInfo", "new token")
	return newCtx, nil
}

```


## Interceptor for api gateway

```go
func run () error {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(GetUnaryClientInterceptor())),
	}
}

// Your unary interceptor
func GetUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		methodName := fmt.Sprintf("client:%s", method)
		log.Println(methodName)
		err := invoker(ctx, method, req, reply, cc, opts...)
		return err
	}
}

```


## Error handling for grpc

```go
	if !usr.IsAdmin() {
		span.SetTag("error", "User is not authorized to perform this action")
		return nil, grpc.Errorf(codes.Unauthenticated, "User is not authorized to perform this action")
	}
```

## To enable opentracing for the jaeger library:

```
$ go get -u github.com/uber/jaeger-client-go/
$ cd $GOPATH/src/github.com/uber/jaeger-client-go/
$ git submodule update --init --recursive
$ make install
```


And then delete the vendored opentracing library from the jaeger vendor folder.

### Securing mongodb

```bash
$ mongo admin
$ db.auth("admin", "password")
```

### Add GPG keys

```bash
$ git commit -S -am "eat: add gpg key"

# To refresh the key,  gpg --refresh-keys
# some random message
```