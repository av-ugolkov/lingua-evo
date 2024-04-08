package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/av-ugolkov/lingua-evo/internal/config"
	pg "github.com/av-ugolkov/lingua-evo/internal/db/postgres"
	"github.com/av-ugolkov/lingua-evo/internal/db/redis"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/kafka"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/analytic"
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
	tagRepository "github.com/av-ugolkov/lingua-evo/internal/services/tag/delivery/repository"
	userService "github.com/av-ugolkov/lingua-evo/internal/services/user"
	userHandler "github.com/av-ugolkov/lingua-evo/internal/services/user/delivery/handler"
	userRepository "github.com/av-ugolkov/lingua-evo/internal/services/user/delivery/repository"
	vocabularyService "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"
	vocabularyHandler "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary/delivery/handler"
	vocabularyRepository "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary/delivery/repository"
	wordService "github.com/av-ugolkov/lingua-evo/internal/services/word"
	wordHandler "github.com/av-ugolkov/lingua-evo/internal/services/word/delivery/handler"
	wordRepository "github.com/av-ugolkov/lingua-evo/internal/services/word/delivery/repository"
)

func ServerStart(cfg *config.Config) {
	if cfg.PprofDebug.Enable {
		go func() {
			slog.Error("%v", http.ListenAndServe(cfg.PprofDebug.Addr(), nil))
		}()
	}

	db, err := pg.NewDB(cfg.DbSQL.GetConnStr())
	if err != nil {
		slog.Error(fmt.Errorf("can't create pg pool: %v", err).Error())
		return
	}

	if cfg.Kafka.Enable {
		kafkaUserAction := kafka.NewWriter(cfg.Kafka.Addr(), cfg.Kafka.Topics[0])
		analytics.SetKafka(kafkaUserAction)
	}

	redisDB := redis.New(cfg)

	router := mux.NewRouter()
	initServer(router, db, redisDB)

	address := fmt.Sprintf(":%s", cfg.Service.Port)

	listener, err := net.Listen(cfg.Service.Type, address)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	slog.Info(fmt.Sprintf("web address: %s", listener.Addr()))

	c := cors.New(cors.Options{
		AllowedOrigins:   cfg.Service.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Fingerprint"},
	})
	handler := c.Handler(router)

	server := &http.Server{
		Handler:      handler,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	slog.Info("start sertver")
	go func() {
		if err := server.Serve(listener); err != nil {
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
	if err := server.Shutdown(context.TODO()); err != nil {
		slog.Info(fmt.Sprintf("server shutdown returned an err: %v\n", err))
	}

	slog.Info("final")
}

func initServer(r *mux.Router, db *sql.DB, redis *redis.Redis) {
	slog.Info("create services")
	userRepo := userRepository.NewRepo(db)
	userSvc := userService.NewService(userRepo, redis)
	wordRepo := wordRepository.NewRepo(db)
	wordSvc := wordService.NewService(wordRepo)
	langRepo := langRepository.NewRepo(db)
	langSvc := langService.NewService(langRepo)
	exampleRepo := exampleRepository.NewRepo(db)
	exampleSvc := exampleService.NewService(exampleRepo)
	tagRepo := tagRepository.NewRepo(db)
	tagSvc := tagService.NewService(tagRepo)
	vocabularyRepo := vocabularyRepository.NewRepo(db)
	dictRepo := dictRepository.NewRepo(db)
	dictSvc := dictService.NewService(dictRepo, vocabularyRepo, langSvc)
	vocabularySvc := vocabularyService.NewService(vocabularyRepo, dictSvc, wordSvc, exampleSvc, tagSvc)
	authRepo := authRepository.NewRepo(redis)
	authSvc := authService.NewService(authRepo, userSvc)

	slog.Info("create handlers")
	userHandler.Create(r, userSvc)
	languageHandler.Create(r, langSvc)
	wordHandler.Create(r, wordSvc, langSvc)
	dictHandler.Create(r, dictSvc)
	vocabularyHandler.Create(r, vocabularySvc)
	authHandler.Create(r, authSvc)

	slog.Info("end init services")
}
