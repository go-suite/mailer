package server

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginLog "github.com/go-mods/zerolog-gin"
	"github.com/go-mods/zerolog-quick/console/colored"
	"github.com/go-suite/mailer/config"
	"github.com/go-suite/mailer/controller"
	"github.com/go-suite/mailer/middleware"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/http"
	"os"
	"time"
)

type Server struct {
	Router *gin.Engine
}

var MailerServer Server
var MainLogger zerolog.Logger
var FileLogger zerolog.Logger

func init() {
	c := config.C.LogFile

	// console writer
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    true,
		TimeFormat: "2006-01-02 15:04:05",
		PartsOrder: []string{
			zerolog.TimestampFieldName,
			zerolog.MessageFieldName,
		},
		FormatExtra: colored.Colorize,
	}

	// console logger
	MainLogger = zerolog.New(consoleWriter).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Logger()

	// file writer
	fileWriter := &lumberjack.Logger{
		Filename:   c.Filename,
		MaxSize:    c.MaxSize,
		MaxBackups: c.MaxBackups,
		MaxAge:     c.MaxAge,
		Compress:   c.Compress,
	}

	// file logger
	FileLogger = zerolog.New(fileWriter).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Logger()
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
	s.Router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		var statusColor, methodColor, resetColor string
		if param.IsOutputColor() {
			statusColor = param.StatusCodeColor()
			methodColor = param.MethodColor()
			resetColor = param.ResetColor()
		}
		return fmt.Sprintf("[MAILER] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			statusColor, param.StatusCode, resetColor,
			param.Latency,
			param.ClientIP,
			methodColor, param.Method, resetColor,
			param.Path,
			param.ErrorMessage,
		)
	}))
	s.Router.Use(gin.Recovery())
	s.Router.SetHTMLTemplate(html_index)

	// Add gin middleware to enable CORS support
	s.Router.Use(cors.Default())

	// initialize routes
	s.initializeRoutes()
}

func (s *Server) initializeRoutes() {

	// Not logged routes
	routes := s.Router.Group("/")

	// Logged routes
	loggedRoutes := s.Router.Group("/")
	loggedRoutes.Use(ginLog.LoggerWithOptions(&ginLog.Options{
		Logger:        &FileLogger,
		FieldsExclude: []string{ginLog.BodyFieldName},
	}))

	// Add a homepage
	loggedRoutes.GET("/", controller.Home)

	// Add ping handler to test if the s in online
	loggedRoutes.GET("/check", controller.Check)

	// Add info handler to info's about mailer
	loggedRoutes.GET("/info", controller.Info)

	// If a list of users is defined, the user need to authenticate
	if len(config.C.Users) > 0 {
		// Add token handler
		routes.POST("/token", controller.Token)

		// Following routes require to be authenticated
		authorized := loggedRoutes.Group("")
		authorized.Use(middleware.TokenAuthMiddleware())
		{
			// Add send handler
			authorized.POST("/send", controller.Send)
		}
	} else {
		// Add send handler
		loggedRoutes.POST("/send", controller.Send)
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

	MainLogger.Info().
		Str("port", addr).
		Msg("Starting Mailer Server")

	MainLogger.Fatal().
		Err(server.ListenAndServe()).
		Msg("Mailer Server Closed")
}

func Run() {
	// Create an instance of the Mailer server and run it
	MailerServer = Server{}
	MailerServer.initialize()
	MailerServer.run(":8080")
}
