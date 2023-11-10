package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"
	"sync"
	"techschool/api"
	db "techschool/db/sqlc"
	"techschool/gapi"
	"techschool/pb"
	"techschool/util"
	"techschool/worker"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	if config.Environment == "development" {

		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	}

	pg, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		panic("panik bang")
	}

	err = pg.Ping()
	if err != nil {
		log.Fatal().Msg("cannot connect database")

	}

	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(pg)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	worker := worker.NewRedisTaskDistributor(redisOpt)
	go runTaskProcessor(redisOpt, store)
	go runGatewayServer(config, store, worker)
	runGrpcServer(config, store, worker)

}

func runDBMigration(migrationURL, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Msg("cannot create new migration instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Msg("failed to run migrate up")
	}

	log.Print("db migration successfully")
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

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) {

	proccessor := worker.NewRedisTaskProcessor(redisOpt, store)
	log.Info().Msg("start task processor")
	err := proccessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}

}

func runGrpcServer(config *util.Config, store db.Store, worker worker.TaskDistributor) {

	server, err := gapi.NewServer(store, config, worker)
	if err != nil {
		log.Fatal().Msg("cannot create server")
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot create listener")
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Msg("cannot serve grpc server")
	}
}

func runGatewayServer(config *util.Config, store db.Store, worker worker.TaskDistributor) {

	server, err := gapi.NewServer(store, config, worker)
	if err != nil {
		log.Fatal().Msg("cannot create server")
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
		log.Fatal().Msg("cannot register handler server")

	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot create listener")
	}

	log.Printf("start HTTP server at %s", listener.Addr().String())
	handler := gapi.HttpLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Msg("cannot serve HTTP gateway server")
	}
}
