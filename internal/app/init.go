package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	jsoniter "github.com/json-iterator/go"

	"github.com/av-ugolkov/lingua-evo/internal/closer"
	"github.com/av-ugolkov/lingua-evo/internal/config"
	pg "github.com/av-ugolkov/lingua-evo/internal/db/postgres"
	"github.com/av-ugolkov/lingua-evo/internal/db/redis"
	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/google"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/kafka"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/analytic"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/log"
	accessHandler "github.com/av-ugolkov/lingua-evo/internal/services/access/handler"
	accessRepository "github.com/av-ugolkov/lingua-evo/internal/services/access/repository"
	accessService "github.com/av-ugolkov/lingua-evo/internal/services/access/service"
	authHandler "github.com/av-ugolkov/lingua-evo/internal/services/auth/handler"
	authRepository "github.com/av-ugolkov/lingua-evo/internal/services/auth/repository"
	authService "github.com/av-ugolkov/lingua-evo/internal/services/auth/service"
	dictService "github.com/av-ugolkov/lingua-evo/internal/services/dictionary"
	dictHandler "github.com/av-ugolkov/lingua-evo/internal/services/dictionary/handler"
	dictRepository "github.com/av-ugolkov/lingua-evo/internal/services/dictionary/repository"
	emailService "github.com/av-ugolkov/lingua-evo/internal/services/email"
	eventsHandler "github.com/av-ugolkov/lingua-evo/internal/services/events/handler"
	eventRepository "github.com/av-ugolkov/lingua-evo/internal/services/events/repository"
	eventService "github.com/av-ugolkov/lingua-evo/internal/services/events/service"
	exampleService "github.com/av-ugolkov/lingua-evo/internal/services/example"
	exampleRepository "github.com/av-ugolkov/lingua-evo/internal/services/example/repository"
	langService "github.com/av-ugolkov/lingua-evo/internal/services/language"
	languageHandler "github.com/av-ugolkov/lingua-evo/internal/services/language/handler"
	langRepository "github.com/av-ugolkov/lingua-evo/internal/services/language/repository"
	notificationService "github.com/av-ugolkov/lingua-evo/internal/services/notifications"
	notificationHandler "github.com/av-ugolkov/lingua-evo/internal/services/notifications/handler"
	notificationRepository "github.com/av-ugolkov/lingua-evo/internal/services/notifications/repository"
	subscribersService "github.com/av-ugolkov/lingua-evo/internal/services/subscribers"
	subscribersHandler "github.com/av-ugolkov/lingua-evo/internal/services/subscribers/handler"
	subscribersRepository "github.com/av-ugolkov/lingua-evo/internal/services/subscribers/repository"
	supportService "github.com/av-ugolkov/lingua-evo/internal/services/support"
	supportHandler "github.com/av-ugolkov/lingua-evo/internal/services/support/handler"
	tagService "github.com/av-ugolkov/lingua-evo/internal/services/tag"
	tagHandler "github.com/av-ugolkov/lingua-evo/internal/services/tag/handler"
	tagRepository "github.com/av-ugolkov/lingua-evo/internal/services/tag/repository"
	userHandler "github.com/av-ugolkov/lingua-evo/internal/services/user/handler"
	userRepository "github.com/av-ugolkov/lingua-evo/internal/services/user/repository"
	userService "github.com/av-ugolkov/lingua-evo/internal/services/user/service"
	vocabHandler "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary/handler"
	vocabRepository "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary/repository"
	vocabService "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary/service"
)

func ServerStart(cfg *config.Config) {
	logger := log.CustomLogger(&cfg.Logger)
	if logger == nil {
		return
	}
	slog.SetDefault(logger.Log)
	closer.Add(func(ctx context.Context) error {
		logger.Close()
		return nil
	})
	if cfg.PprofDebug.Enable {
		initPprof(&cfg.PprofDebug)
	}

	google.InitClient(&cfg.Google)

	pgxPool, err := pg.NewDB(cfg.DbSQL.PgxPoolConfig())
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

	router := fiber.New(fiber.Config{
		AppName:      "Lingua Evo",
		JSONEncoder:  jsoniter.Marshal,
		JSONDecoder:  jsoniter.Unmarshal,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if e, ok := err.(*fiber.Error); ok {
				return c.Status(e.Code).JSON(fiber.Map{
					"error": e.Message,
				})
			}

			// Для других ошибок — возвращаем 500
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error",
			})
		},
	})
	router.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(cfg.Service.AllowedOrigins, ","),
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowCredentials: true,
		AllowHeaders:     "Authorization,Content-Type,Fingerprint",
	}))
	initServer(cfg, router, pgxPool, redisDB)

	address := fmt.Sprintf(":%d", cfg.Service.Port)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	slog.Info("start server")
	go func() {
		if cfg.SSL.Enable {
			err = router.ListenTLS(address, cfg.SSL.GetPublic(), cfg.SSL.GetPrivate())
		} else {
			err = router.Listen(address)
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

	shutdownPprof(ctx)

	if err := router.ShutdownWithContext(ctx); err != nil {
		slog.Error(fmt.Sprintf("server shutdown returned an err: %v\n", err))
	}

	err = closer.Close(ctx)
	if err != nil {
		slog.Error(fmt.Sprintf("closer: %v\n", err.Error()))
	}

	slog.Info("final")
}

func initServer(cfg *config.Config, r *fiber.App, pgxPool *pgxpool.Pool, redis *redis.Redis) {
	tr := transactor.NewTransactor(pgxPool)
	slog.Info("create services")
	emailSvc := emailService.NewService(cfg.Email)
	notificationRepo := notificationRepository.NewRepo(tr)
	notificationSvc := notificationService.NewService(notificationRepo)
	eventsRepo := eventRepository.NewRepo(tr)
	eventsSvc := eventService.NewService(tr, eventsRepo, notificationSvc)
	userRepo := userRepository.NewRepo(tr)
	userSvc := userService.NewService(tr, userRepo, redis, emailSvc)
	accessRepo := accessRepository.NewRepo(tr)
	accessSvc := accessService.NewService(accessRepo)
	langRepo := langRepository.NewRepo(tr)
	langSvc := langService.NewService(langRepo)
	dictRepo := dictRepository.NewRepo(tr)
	dictSvc := dictService.NewService(dictRepo, langSvc)
	exampleRepo := exampleRepository.NewRepo(tr)
	exampleSvc := exampleService.NewService(exampleRepo)
	subscribersRepo := subscribersRepository.NewRepo(tr)
	subscribersSvc := subscribersService.NewService(subscribersRepo)
	tagRepo := tagRepository.NewRepo(tr)
	tagSvc := tagService.NewService(tagRepo)
	vocabRepo := vocabRepository.NewRepo(tr)
	vocabSvc := vocabService.NewService(tr, vocabRepo, exampleSvc, dictSvc, tagSvc, subscribersSvc, eventsSvc)
	authRepo := authRepository.NewRepo(redis)
	authSvc := authService.NewService(authRepo, userSvc, emailSvc)
	supportSvc := supportService.NewService(emailSvc)

	slog.Info("create handlers")
	userHandler.Create(r, userSvc)
	languageHandler.Create(r, langSvc)
	dictHandler.Create(r, dictSvc)
	vocabHandler.Create(r, vocabSvc)
	tagHandler.Create(r, tagSvc)
	authHandler.Create(r, authSvc)
	accessHandler.Create(r, accessSvc)
	subscribersHandler.Create(r, subscribersSvc)
	notificationHandler.Create(r, notificationSvc)
	supportHandler.Create(r, supportSvc)
	eventsHandler.Create(r, eventsSvc)

	slog.Info("end init services")
}
