package main

import "github.com/gin-gonic/gin"

// Server ....
type Server struct {
	// store  *db.Store
	router *gin.Engine
}

// NewServer .....
//func NewServer(store *db.Store) *Server {
func NewServer() *Server {
	//server := &Server{store: store}
	server := &Server{}
	router := gin.Default()

	// TODO: add routes to router
	router.GET("/users", server.GetUsers)

	server.router = router
	return server
}

// Start ....
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
