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

To enable opentracing for the jaeger library:

```
$ go get -u github.com/uber/jaeger-client-go/
$ cd $GOPATH/src/github.com/uber/jaeger-client-go/
$ git submodule update --init --recursive
$ make install
```

And then delete the vendored opentracing library from the jaeger vendor folder.