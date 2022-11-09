package servers

import (
	"fmt"
	"github.com/gennesseaux/mailer/config"
	"github.com/gennesseaux/mailer/controllers"
	"github.com/gennesseaux/mailer/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Server struct {
	Router *gin.Engine
}

var MailerServer Server

func (s *Server) Initialize() {

	//
	if !config.C.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	s.Router = gin.New()
	s.Router.Use(gin.Logger(), gin.Recovery())

	// Add gin middleware to enable CORS support
	s.Router.Use(cors.Default())

	// Initialize routes
	s.InitializeRoutes()
}

func (s *Server) InitializeRoutes() {

	// Add ping handler to test if the s in online
	s.Router.GET("/check", controllers.Check)

	// If a list of users is defined, the user need to authenticate
	if len(config.C.Users) > 0 {
		// Add token handler
		s.Router.POST("/token", controllers.Token)

		// Following routes require to be authenticated
		authorized := s.Router.Group("/")
		authorized.Use(middleware.TokenAuthMiddleware())
		{
			// Add send handler
			authorized.POST("/send", controllers.Send)
		}
	} else {
		// Add send handler
		s.Router.POST("/send", controllers.Send)
	}
}

func (s *Server) Run(addr string) {
	log.Println(fmt.Sprintf("Listen on port %s \n", addr))
	log.Fatal(http.ListenAndServe(addr, s.Router))
}

func Run() {
	log.Println("Running Mailer server ...")

	// Create an instance of the Mailer server and run it
	MailerServer = Server{}
	MailerServer.Initialize()
	MailerServer.Run(":8080")

	log.Println("Mailer server exiting")
}
