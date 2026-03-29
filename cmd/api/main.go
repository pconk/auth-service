package main

import (
	"auth-service/internal/config"
	"auth-service/internal/handler"
	"auth-service/internal/middleware"
	"auth-service/internal/pb"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"auth-service/pkg/database"
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
)

func main() {
	// 1. Setup Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// 2. Load Config
	cfg := config.LoadConfig()

	// 3. Database Connection
	db := database.ConnectDB(cfg)

	// Dapatkan instance sql.DB generic untuk menutup koneksi saat exit
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("Failed to get underlying DB connection", "error", err)
		os.Exit(1)
	}
	defer sqlDB.Close()

	// 4. Init Layers
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg)

	httpHandler := handler.NewAuthHTTPHandler(authService)
	grpcHandler := handler.NewAuthGRPCHandler(authService)

	// 5. Setup HTTP Server
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	// Gunakan Logger custom kita
	r.Use(middleware.Logger(logger))
	// Gunakan Recoverer bawaan Chi
	r.Use(chimiddleware.Recoverer)

	r.Post("/auth/register", httpHandler.Register)
	r.Post("/auth/login", httpHandler.Login)

	httpServer := &http.Server{
		Addr:    ":" + cfg.HTTP_PORT,
		Handler: r,
	}

	// 6. Run Servers in Goroutines
	go func() {
		logger.Info("Starting HTTP Server", "port", cfg.HTTP_PORT)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server error", "error", err)
			os.Exit(1)
		}
	}()

	// gRPC Server Setup
	lis, err := net.Listen("tcp", ":"+cfg.GRPC_PORT)
	if err != nil {
		logger.Error("failed to listen grpc", "error", err)
		os.Exit(1)
	}
	grpcLoggerInterceptor := middleware.NewGrpcLoggerInterceptor(logger)
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(grpcLoggerInterceptor.Unary()))
	pb.RegisterAuthServiceServer(grpcServer, grpcHandler)

	go func() {
		logger.Info("Starting gRPC Server", "port", cfg.GRPC_PORT)
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("gRPC server error", "error", err)
			os.Exit(1)
		}
	}()

	// 7. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("HTTP Shutdown error", "error", err)
	}
	grpcServer.GracefulStop()

	logger.Info("Servers exited successfully")
}
