package app

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"lingua-evo/internal/config"
	pg "lingua-evo/internal/db/postgres"
	"lingua-evo/internal/db/redis"
	authService "lingua-evo/internal/services/auth"
	authHandler "lingua-evo/internal/services/auth/delivery/handler"
	authRepository "lingua-evo/internal/services/auth/delivery/repository"
	dictService "lingua-evo/internal/services/dictionary"
	dictHandler "lingua-evo/internal/services/dictionary/delivery/handler"
	dictRepository "lingua-evo/internal/services/dictionary/delivery/repository"
	exampleService "lingua-evo/internal/services/example"
	exampleRepository "lingua-evo/internal/services/example/delivery/repository"
	langService "lingua-evo/internal/services/language"
	languageHandler "lingua-evo/internal/services/language/delivery/handler"
	langRepository "lingua-evo/internal/services/language/delivery/repository"
	tagService "lingua-evo/internal/services/tag"
	tagRepository "lingua-evo/internal/services/tag/delivery/repository"
	userService "lingua-evo/internal/services/user"
	userHandler "lingua-evo/internal/services/user/delivery/handler"
	userRepository "lingua-evo/internal/services/user/delivery/repository"
	vocabularyService "lingua-evo/internal/services/vocabulary"
	vocabularyHandler "lingua-evo/internal/services/vocabulary/delivery/handler"
	vocabularyRepository "lingua-evo/internal/services/vocabulary/delivery/repository"
	wordService "lingua-evo/internal/services/word"
	wordHandler "lingua-evo/internal/services/word/delivery/handler"
	wordRepository "lingua-evo/internal/services/word/delivery/repository"
)

func ServerStart(cfg *config.Config) {
	if cfg.PprofDebug.Enable {
		go func() {
			slog.Error("%v", http.ListenAndServe("localhost:6060", nil))
		}()
	}

	db, err := pg.NewDB(cfg.DbSQL.GetConnStr())
	if err != nil {
		slog.Error(fmt.Errorf("can't create pg pool: %v", err).Error())
		return
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
		AllowedOrigins:   []string{"http://localhost:5173"},
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
	slog.Info("start sertver")
	if err := server.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			slog.Warn("server shutdown")
		default:
			slog.Error(err.Error())
			return
		}
	}
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
	vocabularySvc := vocabularyService.NewService(vocabularyRepo, wordSvc, exampleSvc, tagSvc)
	dictRepo := dictRepository.NewRepo(db)
	dictSvc := dictService.NewService(dictRepo, vocabularyRepo)
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
