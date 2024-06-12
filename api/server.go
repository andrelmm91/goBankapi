package api

import (
	db "simplebank/db/sqlc"

	"github.com/gin-gonic/gin"
)

// Server http requests
type Server struct {
	store *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// routes
	router.POST("/account", server.CreateAccount)

	server.router = router

	return server
}

// Start runs the http on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}


// handling error
func errorRespose(err error) gin.H {
	return gin.H{"error": err.Error()}
}