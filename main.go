package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"
	"sync"
	"techschool/api"
	db "techschool/db/sqlc"
	"techschool/gapi"
	"techschool/pb"
	"techschool/util"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {

	config, err := util.LoadConfig(".env")
	if err != nil {
		panic(err)
	}

	pg, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		panic("panik bang")
	}

	err = pg.Ping()
	if err != nil {
		log.Fatal("cannot connect database")

	}

	store := db.NewStore(pg)
	go runGatewayServer(config, store)
	runGrpcServer(config, store)

}

func runGinServer(config *util.Config, store db.Store) {
	server := api.NewServer(store, config)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {

		defer wg.Done()
		server.Start(config.HTTPServerAddress)
	}()

	wg.Wait()
}

func runGrpcServer(config *util.Config, store db.Store) {

	server, err := gapi.NewServer(store, config)
	if err != nil {
		log.Fatal("cannot create server")
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener")
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot serve grpc server")
	}
}

func runGatewayServer(config *util.Config, store db.Store) {

	server, err := gapi.NewServer(store, config)
	if err != nil {
		log.Fatal("cannot create server")
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register handler server")

	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot create listener")
	}

	log.Printf("start HTTP server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot serve HTTP gateway server")
	}
}
