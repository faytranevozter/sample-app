package main

import (
	"app/helpers/connection"
	gormrepo "app/internal/repository/gorm"
	"app/internal/rest"
	"app/internal/rest/middleware"
	"app/user"

	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func init() {
	_ = godotenv.Load()
}

func main() {
	// default log level
	logLevel := logger.Info

	if os.Getenv("GO_ENV") == "production" || os.Getenv("GO_ENV") == "prod" {
		// gin mode realease when go env is production
		gin.SetMode(gin.ReleaseMode)

		logLevel = logger.Silent
	}

	timeoutStr := os.Getenv("TIMEOUT")
	if timeoutStr == "" {
		timeoutStr = "5"
	}
	timeout, _ := strconv.Atoi(timeoutStr)
	timeoutContext := time.Duration(timeout) * time.Second

	// logger
	writers := make([]io.Writer, 0)
	if logSTDOUT, _ := strconv.ParseBool(os.Getenv("LOG_TO_STDOUT")); logSTDOUT {
		writers = append(writers, os.Stdout)
	}

	if logFILE, _ := strconv.ParseBool(os.Getenv("LOG_TO_FILE")); logFILE {
		logMaxSize, _ := strconv.Atoi(os.Getenv("LOG_MAX_SIZE"))
		if logMaxSize == 0 {
			logMaxSize = 50 //default 50 megabytes
		}

		logFilename := os.Getenv("LOG_FILENAME")
		if logFilename == "" {
			logFilename = "server.log"
		}

		lg := &lumberjack.Logger{
			Filename:   logFilename,
			MaxSize:    logMaxSize,
			MaxBackups: 1,
			LocalTime:  true,
		}

		writers = append(writers, lg)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(io.MultiWriter(writers...))

	// set gin writer to logrus
	gin.DefaultWriter = logrus.StandardLogger().Writer()

	// init redis database
	var redisClient *redis.Client
	if useRedis, err := strconv.ParseBool(os.Getenv("USE_REDIS")); err == nil && useRedis {
		redisClient = connection.NewRedis(timeoutContext, os.Getenv("REDIS_URL"))
	}

	// Initialize GORM repository with custom logger
	gormDB, err := gorm.Open(
		connection.NewMysqlGORM(timeoutContext, os.Getenv("DB_URL")), // for mysql
		// connection.NewPostgresGORM(timeoutContext, os.Getenv("DB_URL")), // for postgres
		// connection.NewSQLiteGORM(timeoutContext, os.Getenv("DB_URL")), // for sqlite
		// connection.NewSQLServerGORM(timeoutContext, os.Getenv("DB_URL")), // for sql server
		&gorm.Config{
			Logger: logger.New(
				// logrus.StandardLogger(),
				logrus.New(),
				logger.Config{
					SlowThreshold:             200 * time.Millisecond,
					LogLevel:                  logLevel,
					IgnoreRecordNotFoundError: true,
					Colorful:                  false,
				},
			),
		},
	)
	if err != nil {
		panic("error connecting to gorm: " + err.Error())
	}

	// init repo
	userRepo := gormrepo.NewUserRepository(gormDB)

	// Build service Layer
	userService := user.NewService(userRepo)

	// init middleware
	mdl := middleware.NewMiddleware(redisClient)

	// init gin
	ginEngine := gin.New()

	// add exception handler
	// ginEngine.Use(mdl.Recovery())

	// add logger
	ginEngine.Use(mdl.Logger(io.MultiWriter(writers...)))

	// cors
	ginEngine.Use(mdl.Cors())

	// default route
	ginEngine.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]any{
			"message": "It works",
		})
	})

	// init route
	rest.NewUserHandler(ginEngine.Group(""), userService, mdl)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5050"
	}

	logrus.Infof("Service running on port %s", port)
	ginEngine.Run(":" + port)
}
