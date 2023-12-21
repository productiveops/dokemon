package server

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/productiveops/dokemon/pkg/common"
	"github.com/productiveops/dokemon/pkg/crypto/ske"
	"github.com/productiveops/dokemon/pkg/server/handler"
	"github.com/productiveops/dokemon/pkg/server/model"
	"github.com/productiveops/dokemon/pkg/server/requestutil"
	"github.com/productiveops/dokemon/pkg/server/router"
	"github.com/productiveops/dokemon/pkg/server/store"

	"github.com/glebarez/sqlite"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Server struct {
	Echo *echo.Echo
	handler *handler.Handler
}

func (s *Server) Init(dbConnectionString string, dataPath string, logLevel string) {
	setLogLevel(logLevel)
	log.Info().Msg("Starting Dokemon v" + common.Version)

	if dataPath == "" {
		dataPath = "/data"
	}

	if dbConnectionString == "" {
		dbConnectionString = dataPath + "/db"
	}

	// Init compose projects directory
	composeProjectsPath := dataPath + "/compose"
	os.MkdirAll(composeProjectsPath, os.ModePerm)

	// Initialize database
	db, err := gorm.Open(sqlite.Open(dbConnectionString), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(
		&model.ComposeLibraryItem{},
		&model.Credential{},
		&model.Environment{},
		&model.Node{},
		&model.NodeComposeProject{},
		&model.Setting{},
		&model.User{},
		&model.Variable{},
		&model.VariableValue{},
	)

	db.FirstOrCreate(&model.Setting{Id: "SERVER_URL", Value: ""})
	db.FirstOrCreate(&model.Node{Id: 1, Name: "[Dokemon Server]", TokenHash: nil, LastPing: nil})

	sqlNodeComposeProjectStore := store.NewSqlNodeComposeProjectStore(db, composeProjectsPath)
	h := handler.NewHandler(
		composeProjectsPath,
		store.NewSqlComposeLibraryStore(db),
		store.NewSqlCredentialStore(db),
		store.NewSqlEnvironmentStore(db),
		store.NewSqlUserStore(db),
		store.NewSqlNodeStore(db),
		sqlNodeComposeProjectStore,
		store.NewSqlSettingStore(db),
		store.NewSqlVariableStore(db),
		store.NewSqlVariableValueStore(db),
		store.NewLocalFileSystemComposeLibraryStore(db, composeProjectsPath),
		)

	// Init encryption key
	keyFile := dataPath + "/key"
	if _, err := os.Stat(keyFile); errors.Is(err, os.ErrNotExist) {
		log.Info().Msg("key file does not exist. Generating new key.")
 		f, err := os.Create(keyFile)
		if(err != nil){
			log.Fatal().Err(err).Msg("Error while creating key file")
		}
		key, err := ske.GenerateRandomKey()
		if err != nil {
			log.Fatal().Err(err).Msg("Error while generating random key")
		}
		f.WriteString(key)
 	}

	keyBytes, err := os.ReadFile(keyFile)
	if err != nil {
		log.Fatal().Err(err).Msg("Error while reading key file")
	}
	ske.Init(string(keyBytes))

	err = sqlNodeComposeProjectStore.UpdateOldVersionRecords()
	if err != nil {
		log.Error().Err(err).Msg("Error while updating old version data")
	}

	// Web Server
	s.handler = h
	s.Echo = router.New()
	s.Echo.HideBanner = true
	s.Echo.Use(s.authMiddleware)
	s.Echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
		AllowMethods: []string{"*"}, 
	}))
	h.Register(s.Echo)
}

func (s *Server) Run(addr string) {
	if addr == "" {
		addr = ":9090"
	}

	err := s.Echo.Start(addr)
	if err != nil {
		panic(err)
	}
}

func (s *Server) authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !strings.HasPrefix(c.Request().URL.Path, "/api/") || strings.HasPrefix(c.Request().URL.Path, "/api/v1/users") {
			return next(c)
		}

		cc, err := requestutil.GetAuthCookie(c)
		if err != nil {
			return c.NoContent(http.StatusUnauthorized)
		}

		if time.Now().After(cc.Expiry) {
			log.Info().Str("userName", cc.UserName).Msg("Login session expired for user")
			return c.NoContent(http.StatusUnauthorized)
		}

		c.Set("userName", cc.UserName)

		return next(c)
	}
}

func setLogLevel(logLevel string) {
	log.Info().Str("level", logLevel).Msg("Setting log level")
	switch logLevel {
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "ERROR":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "FATAL":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "INFO":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "WARN":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "PANIC":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "TRACE":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
