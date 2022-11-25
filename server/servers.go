package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/go-mods/qlog"
	"github.com/go-mods/qlog/console/colored"
	"github.com/go-suite/mailer/config"
	"github.com/go-suite/mailer/controller"
	"github.com/go-suite/mailer/middleware"
	"net/http"
	"time"
)

type Server struct {
	Router *gin.Engine
}

var MailerServer Server

func init() {
	log.Logger = colored.DateTimeMessage
}

func (s *Server) initialize() {

	//
	if !config.C.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	s.Router = gin.New()
	s.Router.StaticFile("/favicon.ico", "./assets/mailer.ico")
	s.Router.Static("/assets", "./assets")
	s.Router.Use(gin.Logger(), gin.Recovery())
	s.Router.SetHTMLTemplate(html_index)

	// Add gin middleware to enable CORS support
	s.Router.Use(cors.Default())

	// initialize routes
	s.initializeRoutes()
}

func (s *Server) initializeRoutes() {

	// Add a homepage
	s.Router.GET("/", controller.Home)

	// Add ping handler to test if the s in online
	s.Router.GET("/check", controller.Check)

	// Add info handler to info's about mailer
	s.Router.GET("/info", controller.Info)

	// If a list of users is defined, the user need to authenticate
	if len(config.C.Users) > 0 {
		// Add token handler
		s.Router.POST("/token", controller.Token)

		// Following routes require to be authenticated
		authorized := s.Router.Group("/")
		authorized.Use(middleware.TokenAuthMiddleware())
		{
			// Add send handler
			authorized.POST("/send", controller.Send)
		}
	} else {
		// Add send handler
		s.Router.POST("/send", controller.Send)
	}
}

func (s *Server) run(addr string) {
	// Server configuration
	server := &http.Server{
		Addr:              addr,
		Handler:           s.Router,
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	log.Info().
		Str("port", addr).
		Msgf("Starting Mailer Server on port '%s'", addr)

	log.Fatal().
		Err(server.ListenAndServe()).
		Msg("Mailer Server Closed")
}

func Run() {
	// Create an instance of the Mailer server and run it
	MailerServer = Server{}
	MailerServer.initialize()
	MailerServer.run(":8080")
}
