package main

import (
	"context"
	"net"
	stdhttp "net/http"
	"os"
	"time"

	"github.com/Korbinian0/GameServer-panel/backend/internal/adapters/auth"
	"github.com/Korbinian0/GameServer-panel/backend/internal/adapters/db"
	"github.com/Korbinian0/GameServer-panel/backend/internal/adapters/grpc"
	httprest "github.com/Korbinian0/GameServer-panel/backend/internal/adapters/http"
	"github.com/Korbinian0/GameServer-panel/backend/internal/adapters/websocket"
	"github.com/Korbinian0/GameServer-panel/backend/internal/app"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := log.With().Str("service", "gateway").Logger()

	dbURL := getenv("DATABASE_URL", "postgres://gateway:gatewaypass@localhost:5432/gateway?sslmode=disable")
	jwtSecret := getenv("JWT_SECRET", "changeme")
	grpcCert := getenv("GRPC_TLS_CERT", "/certs/server.crt")
	grpcKey := getenv("GRPC_TLS_KEY", "/certs/server.key")

	repo, err := db.NewPostgresRepository(context.Background(), dbURL)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to database")
	}

	authService := auth.NewJWTService([]byte(jwtSecret))
	registry := grpcclient.NewRegistry(grpcCert)
	appService := app.NewGatewayService(repo, authService, registry)

	grpcServer, err := newGRPCServer(grpcCert, grpcKey)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to configure grpc server")
	}

	go func() {
		listenGRPC(logger, grpcServer)
	}()

	hub := websocket.NewHub()
	go hub.Run()

	r := chi.NewRouter()
	httprest.RegisterRoutes(r, appService, authService)
	httprest.RegisterWebsocketRoutes(r, hub, authService)

	api := chi.NewRouter()
	api.Use(authService.Middleware)
	httprest.RegisterProtectedRoutes(api, appService)
	httprest.RegisterNodeRoutes(api, appService, hub)
	r.Mount("/api", api)

	addr := ":8080"
	logger.Info().Msgf("starting HTTP REST gateway on %s", addr)
	server := &stdhttp.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && err != stdhttp.ErrServerClosed {
		logger.Fatal().Err(err).Msg("HTTP server stopped unexpectedly")
	}
}

func newGRPCServer(certPath, keyPath string) (*grpc.Server, error) {
	creds, err := credentials.NewServerTLSFromFile(certPath, keyPath)
	if err != nil {
		return nil, err
	}
	return grpc.NewServer(grpc.Creds(creds)), nil
}

func listenGRPC(logger zerolog.Logger, server *grpc.Server) {
	addr := ":50051"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to bind grpc listener")
	}
	logger.Info().Msgf("starting gRPC server on %s", addr)
	if err := server.Serve(listener); err != nil {
		logger.Fatal().Err(err).Msg("grpc server failed")
	}
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
