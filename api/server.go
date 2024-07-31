package api

import (
	"fmt"
	db "simplebank/db/sqlc"
	"simplebank/token"
	"simplebank/util"
	"simplebank/worker"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	//swagger
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

// Server configures and holds the HTTP server
// @title Simple Bank API
// @version 1.0
// @description This is the API for the Simple Bank application.
// @contact.name API Support
// @contact.email support@example.com
// @basePath /
type Server struct {
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	router          *gin.Engine
	taskDistributor worker.TaskDistributor
}

// NewServer creates a new Server instance with the given configuration, store, and task distributor.
// @Summary Create a new Server instance
// @Description Initializes a new server with provided configuration and dependencies.
// @Tags server
// @Accept json
// @Produce json
// @Param config body util.Config true "Configuration settings"
// @Param store body db.Store true "Database store"
// @Param taskDistributor body worker.TaskDistributor true "Task distributor"
// @Success 200 {object} Server "Server successfully created"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /server [post]
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	// Custom validator for currency
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	// Custom validator for role
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("role", validRole)
	}

	// setup routes
	server.setupRoutes()

	return server, nil
}

// setupRoutes initializes the routes for the server and Swagger documentation
// @Summary Setup API routes
// @Description Configure API routes and Swagger documentation
// @Tags server
// @Accept json
// @Produce json
// @Router / [get]
func (server *Server) setupRoutes() {
	router := gin.Default()

	// User routes
	router.POST("/users", server.createUser)                     // Creates a new user
	router.POST("/users/login", server.loginUser)                // User login
	router.POST("/tokens/renew_access", server.RenewAccessToken) // Renew access token
	router.GET("/verify_email", server.verifyEmail)              // Verify email

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // Swagger documentation

	// Middleware authentication
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// Authenticated routes
	authRoutes.POST("/accounts", server.createAccount)   // Create a new account
	authRoutes.GET("/accounts/:id", server.getAccount)   // Get account details by ID
	authRoutes.GET("/accounts", server.listAccounts)     // List all accounts
	authRoutes.PATCH("/users/update", server.updateUser) // Update user information
	authRoutes.POST("/transfers", server.createTransfer) // Create a new transfer

	// Set router to the server
	server.router = router
}

// Start runs the HTTP server on the specified address
// @Summary Start the HTTP server
// @Description Starts the HTTP server and listens on the specified address
// @Tags server
// @Accept json
// @Produce json
// @Param address query string true "Address to bind the server"
// @Success 200 {string} string "Server started successfully"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /start [get]
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// errorResponse formats an error message for JSON response
// @Summary Format error response
// @Description Formats error messages for JSON responses
// @Tags server
// @Accept json
// @Produce json
// @Param err query string true "Error message"
// @Success 200 {object} gin.H "Formatted error response"
// @Router /error [get]
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
