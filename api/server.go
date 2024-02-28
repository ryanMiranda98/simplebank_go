package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/ryanMiranda98/simplebank/db/sqlc"
	val"github.com/ryanMiranda98/simplebank/api/validator"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and sets up routing.
func NewServer(store db.Store) *Server {
	router := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", val.ValidCurrency)
	}

	server := &Server{
		store:  store,
		router: router,
	}

	router.GET("/", server.index)
	router.GET("/accounts", server.getAllAccounts)
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.POST("/transfers", server.createTransfer)

	// Catch-all route (NOT FOUND - 404)
	router.NoRoute(func(ctx *gin.Context) {
		err := errors.New("route not found")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
	})

	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) index(ctx *gin.Context) {
	json := struct {
		Route   string `json:"route"`
		Message string `json:"message"`
	}{
		Route:   "/",
		Message: "Welcome to SimpleBank API",
	}
	ctx.JSON(http.StatusOK, json)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
