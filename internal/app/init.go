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

	"lingua-evo/internal/config"
	pg "lingua-evo/internal/db/postgres"
	"lingua-evo/internal/db/redis"
	authHandler "lingua-evo/internal/services/auth/delivery/handler"
	authRepository "lingua-evo/internal/services/auth/delivery/repository"
	authService "lingua-evo/internal/services/auth/service"
	dictHandler "lingua-evo/internal/services/lingua/dictionary/delivery/handler"
	dictRepository "lingua-evo/internal/services/lingua/dictionary/delivery/repository"
	dictService "lingua-evo/internal/services/lingua/dictionary/service"
	exampleRepository "lingua-evo/internal/services/lingua/example/delivery/repository"
	exampleService "lingua-evo/internal/services/lingua/example/service"
	languageHandler "lingua-evo/internal/services/lingua/language/delivery/handler"
	langRepository "lingua-evo/internal/services/lingua/language/delivery/repository"
	langService "lingua-evo/internal/services/lingua/language/service"
	tagRepository "lingua-evo/internal/services/lingua/tag/delivery/repository"
	tagService "lingua-evo/internal/services/lingua/tag/service"
	vocabularyHandler "lingua-evo/internal/services/lingua/vocabulary/delivery/handler"
	vocabularyRepository "lingua-evo/internal/services/lingua/vocabulary/delivery/repository"
	vocabularyService "lingua-evo/internal/services/lingua/vocabulary/service"
	wordHandler "lingua-evo/internal/services/lingua/word/delivery/handler"
	wordRepository "lingua-evo/internal/services/lingua/word/delivery/repository"
	wordService "lingua-evo/internal/services/lingua/word/service"
	sessionService "lingua-evo/internal/services/session/service"
	accountHandler "lingua-evo/internal/services/site/account/delivery/http"
	userHandler "lingua-evo/internal/services/user/delivery/handler"
	userRepository "lingua-evo/internal/services/user/delivery/repository"
	userService "lingua-evo/internal/services/user/service"

	signInHandler "lingua-evo/internal/services/site/auth/sign_in/delivery"
	signUpHandler "lingua-evo/internal/services/site/auth/sign_up/delivery"
	indexHandler "lingua-evo/internal/services/site/index/delivery/handler"
)

const (
	filePath = "/website/"
)

func ServerStart(cfg *config.Config, webPath string) {
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
	initServer(router, db, redisDB, webPath)

	address := fmt.Sprintf(":%s", cfg.Service.Port)

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

func initServer(r *mux.Router, db *sql.DB, redis *redis.Redis, webPath string) {
	fs := http.FileServer(http.Dir(webPath))
	r.PathPrefix(filePath).Handler(http.StripPrefix(filePath, fs))

	slog.Info("create services")
	userRepo := userRepository.NewRepo(db)
	userSvc := userService.NewService(userRepo, redis)
	wordRepo := wordRepository.NewRepo(db)
	wordSvc := wordService.NewService(wordRepo)
	langRepo := langRepository.NewRepo(db)
	langSvc := langService.NewService(langRepo)
	dictRepo := dictRepository.NewRepo(db)
	dictSvc := dictService.NewService(dictRepo)
	exampleRepo := exampleRepository.NewRepo(db)
	exampleSvc := exampleService.NewService(exampleRepo)
	tagRepo := tagRepository.NewRepo(db)
	tagSvc := tagService.NewService(tagRepo)
	vocabularyRepo := vocabularyRepository.NewRepo(db)
	vocabularySvc := vocabularyService.NewService(vocabularyRepo, wordSvc, exampleSvc, tagSvc)
	authRepo := authRepository.NewRepo(redis)
	authSvc := authService.NewService(authRepo, userSvc)
	sessionSvc := sessionService.NewService(redis)

	slog.Info("create handlers")
	indexHandler.Create(r, sessionSvc, userSvc, wordSvc)
	userHandler.Create(r, userSvc)
	signInHandler.Create(r, userSvc)
	signUpHandler.Create(r)
	accountHandler.Create(r)
	languageHandler.Create(r, langSvc)
	wordHandler.Create(r, wordSvc, langSvc)
	dictHandler.Create(r, dictSvc)
	vocabularyHandler.Create(r, vocabularySvc)
	authHandler.Create(r, authSvc)

	slog.Info("end init services")
}
