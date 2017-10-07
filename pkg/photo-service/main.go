package main

import (
	"log"
	"net"

	"github.com/spf13/viper"

	"google.golang.org/grpc"

	"github.com/alextanhongpin/go-grpc-event/internal/database"
	pb "github.com/alextanhongpin/go-grpc-event/proto/photo"
)

func main() {
	//
	// TCP
	//
	lis, err := net.Listen("tcp", viper.GetString("port"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//
	// DATABASE
	//
	db, err := database.New(
		database.Host(viper.GetString("mgo_host")),
		database.Name(viper.GetString("mgo_db")),
		database.Username(viper.GetString("mgo_usr")),
		database.Password(viper.GetString("mgo_pwd")),
	)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	defer db.Close()

	//
	// GRPC
	//
	grpcServer := grpc.NewServer()
	pb.RegisterPhotoServiceServer(grpcServer, &photoServer{
		db: db,
	})
	log.Printf("listening to port *%v", viper.GetString("port"))
	grpcServer.Serve(lis)
}
