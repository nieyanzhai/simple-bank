package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/nieyanzhai/simple-bank/db/sqlc"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	s := Server{
		store: store,
	}
	e := gin.Default()
	e.GET("/accounts/:id", s.GetAccount)
	e.GET("/accounts", s.GetAccounts)
	e.POST("/accounts", s.CreateAccount)
	s.router = e
	return &s
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
