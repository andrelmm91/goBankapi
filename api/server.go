package api

import (
	"fmt"
	db "simplebank/db/sqlc"
	"simplebank/token"
	"simplebank/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server http requests
type Server struct {
	config util.Config
	store  db.Store
	tokenMaker token.Maker
	router *gin.Engine
}

func NewServer(config util.Config,store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token mnaker: %w", err)
	}

	server := &Server{
		config: config,
		store: store, 
		tokenMaker: tokenMaker,
	}

	// Custom validator for currency
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	// setup routes
	server.setupRoutes()

	return server, nil
}

func (server *Server) setupRoutes() {
	router := gin.Default()
	// routes
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	
	router.POST("/transfers", server.createTransfer)

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	// server
	server.router = router
}

// Start runs the http on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// handling error
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
