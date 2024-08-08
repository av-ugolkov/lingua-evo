package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/av-ugolkov/lingua-evo/internal/closer"
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
	subscribersService "github.com/av-ugolkov/lingua-evo/internal/services/subscribers"
	subscribersHandler "github.com/av-ugolkov/lingua-evo/internal/services/subscribers/delivery/handler"
	subscribersRepository "github.com/av-ugolkov/lingua-evo/internal/services/subscribers/delivery/repository"
	tagService "github.com/av-ugolkov/lingua-evo/internal/services/tag"
	tagHandler "github.com/av-ugolkov/lingua-evo/internal/services/tag/delivery/handler"
	tagRepository "github.com/av-ugolkov/lingua-evo/internal/services/tag/delivery/repository"
	userService "github.com/av-ugolkov/lingua-evo/internal/services/user"
	userHandler "github.com/av-ugolkov/lingua-evo/internal/services/user/delivery/handler"
	userRepository "github.com/av-ugolkov/lingua-evo/internal/services/user/delivery/repository"
	vocabHandler "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary/delivery/handler"
	vocabRepository "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary/delivery/repository"
	vocabService "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary/service"
	vocabAccessService "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary_access"
	vocabAccessRepository "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary_access/delivery/repository"
	vocabWordsService "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary_words"
	vocabWordsHandler "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary_words/delivery/handler"
	wordService "github.com/av-ugolkov/lingua-evo/internal/services/word"
	wordHandler "github.com/av-ugolkov/lingua-evo/internal/services/word/delivery/handler"
	wordRepository "github.com/av-ugolkov/lingua-evo/internal/services/word/delivery/repository"
)

func ServerStart(cfg *config.Config) {
	closer := closer.NewCloser()

	logger := log.CustomLogger(&cfg.Logger)
	if logger == nil {
		return
	}
	slog.SetDefault(logger.Log)
	closer.Add(func(ctx context.Context) error {
		logger.Close()
		return nil
	})
	var pprofSrv *http.Server
	if cfg.PprofDebug.Enable {
		go func() {
			pprofSrv = &http.Server{
				Addr: cfg.PprofDebug.Addr(),
			}
			slog.Error(pprofSrv.ListenAndServe().Error())
		}()
	}

	pgxPool, err := pg.NewDB(cfg.DbSQL)
	if err != nil {
		slog.Error(fmt.Sprintf("can't create pg pool: %v", err))
		return
	}
	closer.Add(func(ctx context.Context) error {
		pgxPool.Close()
		return nil
	})

	if cfg.Kafka.Enable {
		kafkaUserAction := kafka.NewWriter(cfg.Kafka.Addr(), cfg.Kafka.Topics[0])
		analytics.SetKafka(kafkaUserAction)
	}

	redisDB := redis.New(cfg)
	closer.Add(func(ctx context.Context) error {
		err := redisDB.Close()
		if err != nil {
			return err
		}
		return nil
	})

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

	server := http.Server{
		Addr:         address,
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		ErrorLog:     logger.ServerLoger,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	slog.Info("start server")
	go func() {
		if cfg.SSL.Enable {
			err = server.ListenAndServeTLS(cfg.SSL.GetPublic(), cfg.SSL.GetPrivate())
		} else {
			err = server.ListenAndServe()
		}

		if err != nil {
			switch {
			case errors.Is(err, http.ErrServerClosed):
				slog.Warn("server shutdown")
			default:
				slog.Error(fmt.Sprintf("server returned an err: %v\n", err.Error()))
				return
			}
		}
	}()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if pprofSrv != nil {
		if err := pprofSrv.Shutdown(ctx); err != nil {
			slog.Error(fmt.Sprintf("server pprof shutdown returned an err: %v\n", err))
		}
	}
	if err := server.Shutdown(ctx); err != nil {
		slog.Error(fmt.Sprintf("server shutdown returned an err: %v\n", err))
	}

	err = closer.Close(ctx)
	if err != nil {
		slog.Error(fmt.Sprintf("closer: %v\n", err.Error()))
	}

	slog.Info("final")
}

func initServer(cfg *config.Config, r *gin.Engine, pgxPool *pgxpool.Pool, redis *redis.Redis) {
	tr := transactor.NewTransactor(pgxPool)
	slog.Info("create services")
	accessRepo := accessRepository.NewRepo(pgxPool)
	accessSvc := accessService.NewService(accessRepo)
	userRepo := userRepository.NewRepo(pgxPool)
	userSvc := userService.NewService(userRepo, redis)
	langRepo := langRepository.NewRepo(pgxPool)
	langSvc := langService.NewService(langRepo)
	dictRepo := dictRepository.NewRepo(pgxPool)
	dictSvc := dictService.NewService(dictRepo, langSvc)
	exampleRepo := exampleRepository.NewRepo(pgxPool)
	exampleSvc := exampleService.NewService(exampleRepo)
	subscribersRepo := subscribersRepository.NewRepo(pgxPool)
	subscribersSvc := subscribersService.NewService(subscribersRepo)
	tagRepo := tagRepository.NewRepo(pgxPool)
	tagSvc := tagService.NewService(tagRepo)
	vocabularyAccessRepo := vocabAccessRepository.NewRepo(pgxPool)
	vocabularyAccessSvc := vocabAccessService.NewService(vocabularyAccessRepo)
	vocabRepo := vocabRepository.NewRepo(pgxPool)
	vocabSvc := vocabService.NewService(tr, vocabRepo, tagSvc, subscribersSvc, vocabularyAccessSvc)
	wordRepo := wordRepository.NewRepo(pgxPool)
	wordSvc := wordService.NewService(tr, wordRepo, userSvc, vocabSvc, vocabularyAccessSvc, dictSvc, exampleSvc)
	vocabWordsSvc := vocabWordsService.NewService(vocabSvc, wordSvc)
	authRepo := authRepository.NewRepo(redis)
	authSvc := authService.NewService(cfg.Email, authRepo, userSvc)

	slog.Info("create handlers")
	userHandler.Create(r, userSvc)
	languageHandler.Create(r, langSvc)
	dictHandler.Create(r, dictSvc)
	wordHandler.Create(r, wordSvc)
	vocabHandler.Create(r, vocabSvc)
	tagHandler.Create(r, tagSvc)
	authHandler.Create(r, authSvc)
	accessHandler.Create(r, accessSvc)
	subscribersHandler.Create(r, subscribersSvc)
	vocabWordsHandler.Create(r, vocabWordsSvc)

	slog.Info("end init services")
}
