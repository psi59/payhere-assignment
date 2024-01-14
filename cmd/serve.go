package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/psi59/payhere-assignment/usecase/authtoken"

	"github.com/psi59/payhere-assignment/repository/mysql"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/psi59/payhere-assignment/handler"
	"github.com/psi59/payhere-assignment/internal/db"
	"github.com/psi59/payhere-assignment/internal/ginhelper"
	"github.com/psi59/payhere-assignment/internal/valid"
	"github.com/psi59/payhere-assignment/middleware"
	"github.com/psi59/payhere-assignment/repository"
	"github.com/psi59/payhere-assignment/usecase/user"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

const flagConfigPath = "config-path"

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run API server",
	Run:   runServeCommand,
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringP(flagConfigPath, "c", "config/server.yaml", "config file path")
}

func runServeCommand(cmd *cobra.Command, _ []string) {
	configPath, err := cmd.Flags().GetString(flagConfigPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get config-path flag")
	}

	apiServer, err := NewAPIServer(configPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create server instance")
	}

	if err := apiServer.Start(); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}

type APIServer struct {
	engine *gin.Engine
	config APIServerConfig

	// Middleware
	AuthMiddleware *middleware.AuthMiddleware

	// Handlers
	UserHandler *handler.UserHandler

	// Usecases
	UserUsecase      user.Usecase
	AuthTokenUsecase authtoken.Usecase

	// Repositories
	UserRepository           repository.UserRepository
	TokenBlacklistRepository repository.TokenBlacklistRepository

	// ETC
	dbConn *gorm.DB
}

func NewAPIServer(configPath string) (*APIServer, error) {
	config, err := loadAPIServerConfig(configPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load config")
	}

	engine := gin.New()
	s := &APIServer{
		engine: engine,
		config: config,
	}

	if err := s.initDB(); err != nil {
		return nil, errors.WithStack(err)
	}
	s.initRepositories()
	if err := s.initUsecase(); err != nil {
		return nil, errors.WithStack(err)
	}
	if err := s.initHandler(); err != nil {
		return nil, errors.WithStack(err)
	}
	if err := s.initMiddleware(); err != nil {
		return nil, errors.WithStack(err)
	}
	s.initRoutes()

	return s, nil
}
func (s *APIServer) Start() error {
	srv := &http.Server{
		Addr:    ":1202",
		Handler: s.engine,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error().Err(err).Msg("server closed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(
		quit,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	sig := <-quit

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	log.Info().Stringer("sig", sig).Msg("Shutting down the server")
	if err := srv.Shutdown(ctx); err != nil {
		return errors.WithStack(err)
	}
	log.Info().Msg("Server exiting")

	return nil
}

func (s *APIServer) initRoutes() {
	engine := s.engine
	engine.GET("/healthcheck", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	engine.GET("/docs", func(c *gin.Context) {
		c.File(s.config.APIDoc)
	})

	v1 := engine.Group("/v1")
	v1.Use(
		requestid.New(),
		middleware.SetContext(),
		middleware.Logger(),
		func(c *gin.Context) {
			ctx := ginhelper.GetContext(c)
			ctx = db.ContextWithConn(ctx, s.dbConn)
			ginhelper.SetContext(c, ctx)
			c.Next()
		},
	)

	{
		v1User := v1.Group("/users")
		v1User.POST("/signUp", s.UserHandler.SignUp)
		v1User.POST("/signIn", s.UserHandler.SignIn)
		v1User.POST("/signOut", s.AuthMiddleware.Auth(), s.UserHandler.SignOut)
	}

}

func (s *APIServer) initMiddleware() error {
	authMiddleware, err := middleware.NewAuthMiddleware(s.UserUsecase, s.AuthTokenUsecase)
	if err != nil {
		return errors.WithStack(err)
	}

	s.AuthMiddleware = authMiddleware

	return nil
}

func (s *APIServer) initHandler() error {
	userHandler, err := handler.NewUserHandler(s.UserUsecase, s.AuthTokenUsecase)
	if err != nil {
		return errors.WithStack(err)
	}

	s.UserHandler = userHandler

	return nil
}

func (s *APIServer) initUsecase() error {
	userService, err := user.NewService(s.UserRepository)
	if err != nil {
		return errors.WithStack(err)
	}
	authTokenService, err := authtoken.NewService(s.config.JWTSecret, s.TokenBlacklistRepository)
	if err != nil {
		return errors.WithStack(err)
	}

	s.UserUsecase = userService
	s.AuthTokenUsecase = authTokenService

	return nil
}

func (s *APIServer) initRepositories() {
	userRepository := mysql.NewUserRepository()
	tokenBlacklistRepository := mysql.NewTokenBlacklistRepository()

	s.UserRepository = userRepository
	s.TokenBlacklistRepository = tokenBlacklistRepository
}

func (s *APIServer) initDB() error {
	dbConn, err := db.Connect(s.config.DB)
	if err != nil {
		return errors.WithStack(err)
	}
	s.dbConn = dbConn

	return nil
}

type APIServerConfig struct {
	APIDoc    string    `yaml:"apiDoc"`
	JWTSecret string    `yaml:"jwtSecret"`
	DB        db.Config `yaml:"db"`
}

func loadAPIServerConfig(configPath string) (config APIServerConfig, err error) {
	f, openErr := os.Open(configPath)
	if openErr != nil {
		err = errors.WithStack(openErr)
		return
	}
	defer func() {
		_ = f.Close()
	}()

	if decodeErr := yaml.NewDecoder(f).Decode(&config); decodeErr != nil {
		err = errors.WithStack(decodeErr)
		return
	}
	if validateErr := config.Validate(); validateErr != nil {
		err = errors.WithStack(validateErr)
		return
	}

	return
}

func (c APIServerConfig) Validate() error {
	return errors.WithStack(valid.ValidateStruct(c))
}
