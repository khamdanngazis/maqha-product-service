// cmd/main.go

package main

import (
	"flag"
	"fmt"
	"log"
	"maqhaa/library/logging"
	"maqhaa/library/middleware"
	exRepo "maqhaa/product_service/external/repository"
	"maqhaa/product_service/internal/app/repository"
	"maqhaa/product_service/internal/app/service"
	"maqhaa/product_service/internal/config"
	"maqhaa/product_service/internal/database"
	grpcHandler "maqhaa/product_service/internal/interface/grpc/handler"
	pb "maqhaa/product_service/internal/interface/grpc/model"
	httpHandler "maqhaa/product_service/internal/interface/http/handler"
	"maqhaa/product_service/internal/interface/http/router"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
)

func main() {
	// Define a command line flag for the config file path
	configFilePath := flag.String("config", "config/config.yaml", "path to the config file")
	logFile := flag.String("log.file", "../logs", "Logging file")

	flag.Parse()

	initLogging(*logFile)

	// Load the configuration
	cfg, err := config.LoadConfig(*configFilePath)
	if err != nil {
		logging.Log.Fatalf("Error loading configuration: %v", err)
	}
	logging.Log.Infof("Load configuration from %v", *configFilePath)
	// Access configuration values
	dbConfig := cfg.Database

	db, err := database.NewDB(&dbConfig)
	if err != nil {
		logging.Log.Fatalf("Error loading configuration: %v", err)
	}

	// Close the database connection when done
	sqlDB, err := db.DB()
	if err != nil {
		logging.Log.Fatalf("Error getting DB connection: %v", err)
	}
	defer sqlDB.Close()

	// Initialize handlers
	httpRouter := router.NewMuxRouter()

	pingHandler := httpHandler.NewPingHandler()
	httpRouter.GET("/ping", pingHandler.Ping)

	// Initialize product service
	userRepository := exRepo.NewUserRepository(cfg.ExternalConnection.AuthService.Host)
	imageRepository := repository.NewImagesRepository(cfg.ImagePath)
	productRepository := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepository, userRepository, imageRepository)
	productHandler := httpHandler.NewProductHandler(productService)

	httpRouter.GET("/product", productHandler.GetProductGroupsByCategoryHandler)
	httpRouter.POST("/product", productHandler.AddProductHandler)
	httpRouter.PUT("/product", productHandler.EditProductHandler)
	httpRouter.DELETE("/product", productHandler.DeactiveProductHandler)
	httpRouter.POST("/category", productHandler.AddCategoryHandler)
	httpRouter.PUT("/category", productHandler.EditCategoryHandler)
	httpRouter.DELETE("/category", productHandler.DeactiveCategoryHandler)

	// Start HTTP server
	go func() {
		httpRouter.SERVE(cfg.AppPort)
	}()

	productHandlerGrpc := grpcHandler.NewProductGRPCHandler(productService)
	// Initialize gRPC server
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(middleware.LoggingInterceptor))

	// Register gRPC service implementation
	pb.RegisterProductServer(grpcServer, productHandlerGrpc)

	// Start gRPC server
	listen, err := net.Listen("tcp", cfg.GrpcPort)
	if err != nil {
		logging.Log.Fatalf("Error starting gRPC server: %s", err)
	}

	defer listen.Close()

	logging.Log.Infof("gRPC server listening on %s", cfg.GrpcPort)
	if err := grpcServer.Serve(listen); err != nil {
		logging.Log.Fatalf("Error serving gRPC: %v", err)
	}
}

func initLogging(logFolder string) {
	logging.InitLogger()
	currentDate := time.Now().Format("2006-01-02")

	// Specify the log file with the current date
	logFilePath := fmt.Sprintf("%s/app_%s.log", logFolder, currentDate)

	// Create the log file if it doesn't exist
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal("Error creating log file:", err)
	}

	// Set the logrus output to the log file
	logging.Log.SetOutput(logFile)

	go func() {
		for {
			time.Sleep(time.Hour) // Adjust the sleep duration as needed
			newDate := time.Now().Format("2006-01-02")
			if newDate != currentDate {
				currentDate = newDate
				logFilePath = fmt.Sprintf("%s/app_%s.log", logFolder, currentDate)
				newLogFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
				if err != nil {
					logging.Log.Fatal("Error creating log file:", err)
				}
				logFile = newLogFile
				logging.Log.SetOutput(logFile)
			}
		}
	}()
}
