package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/tatekicct/simplebank/db/sqlc"
	"github.com/tatekicct/simplebank/token"
	"github.com/tatekicct/simplebank/util"
)

// Server serrves HTTP requests for our banking service
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and sets up routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	// tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannnot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	// Set up custom validation for currency
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	server.setupRouter()

	return server, nil
}

// Start runs the HTTP server on the specified address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	autoRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	autoRoutes.POST("/accounts", server.createAccount)
	autoRoutes.GET("/accounts/:id", server.getAccount)
	autoRoutes.GET("/accounts", server.listAccounts)

	autoRoutes.POST("/transfers", server.createTransfer)

	server.router = router
}
