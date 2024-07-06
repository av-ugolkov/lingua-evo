package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/av-ugolkov/lingua-evo/internal/config"
	pg "github.com/av-ugolkov/lingua-evo/internal/db/postgres"
	"github.com/av-ugolkov/lingua-evo/internal/db/redis"
	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/kafka"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/analytic"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/log"
	accessService "github.com/av-ugolkov/lingua-evo/internal/services/access"
	accessHandler "github.com/av-ugolkov/lingua-evo/internal/services/access/delivery/handler"
	accessRepository "github.com/av-ugolkov/lingua-evo/internal/services/access/delivery/repository"
	authService "github.com/av-ugolkov/lingua-evo/internal/services/auth"
	authHandler "github.com/av-ugolkov/lingua-evo/internal/services/auth/delivery/handler"
	authRepository "github.com/av-ugolkov/lingua-evo/internal/services/auth/delivery/repository"
	dictService "github.com/av-ugolkov/lingua-evo/internal/services/dictionary"
	dictHandler "github.com/av-ugolkov/lingua-evo/internal/services/dictionary/delivery/handler"
	dictRepository "github.com/av-ugolkov/lingua-evo/internal/services/dictionary/delivery/repository"
	exampleService "github.com/av-ugolkov/lingua-evo/internal/services/example"
	exampleRepository "github.com/av-ugolkov/lingua-evo/internal/services/example/delivery/repository"
	langService "github.com/av-ugolkov/lingua-evo/internal/services/language"
	languageHandler "github.com/av-ugolkov/lingua-evo/internal/services/language/delivery/handler"
	langRepository "github.com/av-ugolkov/lingua-evo/internal/services/language/delivery/repository"
	tagService "github.com/av-ugolkov/lingua-evo/internal/services/tag"
	tagHandler "github.com/av-ugolkov/lingua-evo/internal/services/tag/delivery/handler"
	tagRepository "github.com/av-ugolkov/lingua-evo/internal/services/tag/delivery/repository"
	userService "github.com/av-ugolkov/lingua-evo/internal/services/user"
	userHandler "github.com/av-ugolkov/lingua-evo/internal/services/user/delivery/handler"
	userRepository "github.com/av-ugolkov/lingua-evo/internal/services/user/delivery/repository"
	vocabService "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"
	vocabHandler "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary/delivery/handler"
	vocabRepository "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary/delivery/repository"
	wordService "github.com/av-ugolkov/lingua-evo/internal/services/word"
	wordHandler "github.com/av-ugolkov/lingua-evo/internal/services/word/delivery/handler"
	wordRepository "github.com/av-ugolkov/lingua-evo/internal/services/word/delivery/repository"
)

func ServerStart(cfg *config.Config) {
	logger := log.CustomLogger(&cfg.Logger)
	if logger == nil {
		return
	}
	slog.SetDefault(logger.Log)

	if cfg.PprofDebug.Enable {
		go func() {
			slog.Error(http.ListenAndServe(cfg.PprofDebug.Addr(), nil).Error())
		}()
	}

	pgxPool, err := pg.NewDB(cfg.DbSQL)
	if err != nil {
		slog.Error(fmt.Sprintf("can't create pg pool: %v", err))
		return
	}

	if cfg.Kafka.Enable {
		kafkaUserAction := kafka.NewWriter(cfg.Kafka.Addr(), cfg.Kafka.Topics[0])
		analytics.SetKafka(kafkaUserAction)
	}

	redisDB := redis.New(cfg)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Service.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowCredentials: true,
		AllowHeaders:     []string{"Authorization", "Content-Type", "Fingerprint"},
	}))
	initServer(cfg, router, pgxPool, redisDB)

	address := fmt.Sprintf(":%d", cfg.Service.Port)

	listener, err := net.Listen(cfg.Service.Type, address)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	slog.Info(fmt.Sprintf("web address: %s", listener.Addr()))
	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		ErrorLog:     logger.ServerLoger,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	slog.Info("start server")
	go func() {
		if cfg.SSL.Enable {
			err = server.ServeTLS(listener, cfg.SSL.GetPublic(), cfg.SSL.GetPrivate())
		} else {
			err = server.Serve(listener)
		}
		if err != nil {
			switch {
			case errors.Is(err, http.ErrServerClosed):
				slog.Warn("server shutdown")
			default:
				slog.Error(err.Error())
				return
			}
		}
	}()

	<-ctx.Done()
	if err := server.Shutdown(context.Background()); err != nil {
		slog.Info(fmt.Sprintf("server shutdown returned an err: %v\n", err))
		logger.Close()
	}

	slog.Info("final")
}

func initServer(cfg *config.Config, r *gin.Engine, pgxPool *pgxpool.Pool, redis *redis.Redis) {
	tr := transactor.NewTransactor(pgxPool)
	slog.Info("create services")
	userRepo := userRepository.NewRepo(pgxPool)
	userSvc := userService.NewService(userRepo, redis)
	langRepo := langRepository.NewRepo(pgxPool)
	langSvc := langService.NewService(langRepo)
	dictRepo := dictRepository.NewRepo(pgxPool)
	dictSvc := dictService.NewService(dictRepo, langSvc)
	exampleRepo := exampleRepository.NewRepo(pgxPool)
	exampleSvc := exampleService.NewService(exampleRepo)
	tagRepo := tagRepository.NewRepo(pgxPool)
	tagSvc := tagService.NewService(tagRepo)
	vocabRepo := vocabRepository.NewRepo(pgxPool)
	vocabSvc := vocabService.NewService(tr, vocabRepo, langSvc, tagSvc)
	wordRepo := wordRepository.NewRepo(pgxPool)
	wordSvc := wordService.NewService(tr, wordRepo, userSvc, vocabSvc, dictSvc, exampleSvc)
	authRepo := authRepository.NewRepo(redis)
	authSvc := authService.NewService(cfg.Email, authRepo, userSvc)
	accessRepo := accessRepository.NewRepo(pgxPool)
	accessSvc := accessService.NewService(accessRepo)

	slog.Info("create handlers")
	userHandler.Create(r, userSvc)
	languageHandler.Create(r, langSvc)
	dictHandler.Create(r, dictSvc)
	wordHandler.Create(r, wordSvc)
	vocabHandler.Create(r, vocabSvc)
	tagHandler.Create(r, tagSvc)
	authHandler.Create(r, authSvc)
	accessHandler.Create(r, accessSvc)

	slog.Info("end init services")
}
