package server

import (
	"errors"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/productiveops/dokemon/pkg/common"
	"github.com/productiveops/dokemon/pkg/crypto/ske"
	"github.com/productiveops/dokemon/pkg/crypto/ssl"
	"github.com/productiveops/dokemon/pkg/dockerapi"
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
	dataPath string
	sslEnabled bool
}

func NewServer(dbConnectionString string, dataPath string, logLevel string, sslEnabled string) (*Server) {
	s := Server{}

	setLogLevel(logLevel)
	log.Info().Msg("Starting Dokemon v" + common.Version)

	if dataPath == "" {
		dataPath = "/data"
	}
	s.dataPath = dataPath

	if dbConnectionString == "" {
		dbConnectionString = dataPath + "/db"
	}
	
	s.sslEnabled = sslEnabled == "1"

	composeProjectsPath := path.Join(dataPath, "/compose")
	initCompose(composeProjectsPath)
	initEncryption(dataPath)
	db, err := initDatabase(dbConnectionString)
	if err != nil {
		log.Fatal().Err(err).Msg("Error while initializing database")
	}

	// Setup stores
	sqlNodeComposeProjectStore := store.NewSqlNodeComposeProjectStore(db, composeProjectsPath)
	h := handler.NewHandler(
		composeProjectsPath,
		store.NewSqlComposeLibraryStore(db),
		store.NewSqlCredentialStore(db),
		store.NewSqlEnvironmentStore(db),
		store.NewSqlUserStore(db),
		store.NewSqlNodeStore(db),
		sqlNodeComposeProjectStore,
		store.NewSqlNodeComposeProjectVariableStore(db),
		store.NewSqlSettingStore(db),
		store.NewSqlVariableStore(db),
		store.NewSqlVariableValueStore(db),
		store.NewLocalFileSystemComposeLibraryStore(db, composeProjectsPath),
		)

	err = sqlNodeComposeProjectStore.UpdateOldVersionRecords()
	if err != nil {
		log.Error().Err(err).Msg("Error while updating old version data")
	}

	go dockerapi.ContainerScheduleRefreshStaleStatus()

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

	return &s
}

func initCompose(composeProjectsPath string) {
	os.MkdirAll(composeProjectsPath, os.ModePerm)
}

func initEncryption(dataPath string) {
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
}

func initDatabase(dbConnectionString string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbConnectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&model.ComposeLibraryItem{},
		&model.Credential{},
		&model.Environment{},
		&model.Node{},
		&model.NodeComposeProject{},
		&model.NodeComposeProjectVariable{},
		&model.Setting{},
		&model.User{},
		&model.Variable{},
		&model.VariableValue{},
	)
	if err != nil {
		return nil, err
	}

	err = db.FirstOrCreate(&model.Setting{Id: "SERVER_URL", Value: ""}).Error
	if err != nil {
		return nil, err
	}

	err = db.FirstOrCreate(&model.Node{Id: 1, Name: "[Dokemon Server]", TokenHash: nil, LastPing: nil}).Error
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (s *Server) Run(addr string) {
	if addr == "" {
		addr = ":9090"
	}

	var err error
	if s.sslEnabled {
		certsDirPath := path.Join(s.dataPath, "certs")
		certPath := path.Join(certsDirPath, "server.crt")
		keyPath := path.Join(certsDirPath, "server.key")
		s.generateSelfSignedCerts(certsDirPath, certPath, keyPath)

		err = s.Echo.StartTLS(addr, certPath, keyPath)
	} else {
		err = s.Echo.Start(addr)
	}

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

func (s *Server) generateSelfSignedCerts(certDirPath string, certPath string, keyPath string) {
	if _, err := os.Stat(certDirPath); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(certDirPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	if _, err := os.Stat(certPath); errors.Is(err, os.ErrNotExist) {
		log.Debug().Msg("SSL certificate file does not exist. Generating self-signed certificate...")

		cert, key, err := ssl.GenerateSelfSignedCert()
		if err != nil {
			panic(err)
		}

		certFile, err := os.Create(certPath)
		if err != nil {
			panic(err)
		}
		_, err = certFile.WriteString(cert)
		if err != nil {
			panic(err)
		}

		keyFile, err := os.Create(keyPath)
		if err != nil {
			panic(err)
		}
		_, err = keyFile.WriteString(key)
		if err != nil {
			panic(err)
		}
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
