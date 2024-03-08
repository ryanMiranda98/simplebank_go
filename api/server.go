package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	val "github.com/ryanMiranda98/simplebank/api/validator"
	db "github.com/ryanMiranda98/simplebank/db/sqlc"
	"github.com/ryanMiranda98/simplebank/token"
	"github.com/ryanMiranda98/simplebank/util"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and sets up routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", val.ValidCurrency)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.GET("/", server.index)
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/token/renew_access", server.RenewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.GET("/accounts", server.getAllAccounts)
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.POST("/transfers", server.createTransfer)

	// Catch-all route (NOT FOUND - 404)
	router.NoRoute(func(ctx *gin.Context) {
		err := errors.New("route not found")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
	})

	server.router = router
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
