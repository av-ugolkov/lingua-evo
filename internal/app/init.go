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
	repository "lingua-evo/internal/db"
	authHandler "lingua-evo/internal/services/auth/delivery"
	authRepository "lingua-evo/internal/services/auth/repository"
	authService "lingua-evo/internal/services/auth/service"
	accountHandler "lingua-evo/internal/services/lingua/account/delivery"
	dictHandler "lingua-evo/internal/services/lingua/dictionary/delivery"
	dictRepository "lingua-evo/internal/services/lingua/dictionary/repository"
	dictService "lingua-evo/internal/services/lingua/dictionary/service"
	exampleRepository "lingua-evo/internal/services/lingua/example/repository"
	exampleService "lingua-evo/internal/services/lingua/example/service"
	languageHandler "lingua-evo/internal/services/lingua/language/delivery"
	langRepository "lingua-evo/internal/services/lingua/language/repository"
	langService "lingua-evo/internal/services/lingua/language/service"
	tagRepository "lingua-evo/internal/services/lingua/tag/repository"
	tagService "lingua-evo/internal/services/lingua/tag/service"
	vocabularyHandler "lingua-evo/internal/services/lingua/vocabulary/delivery"
	vocabularyRepository "lingua-evo/internal/services/lingua/vocabulary/repository"
	vocabularyService "lingua-evo/internal/services/lingua/vocabulary/service"
	wordHandler "lingua-evo/internal/services/lingua/word/delivery"
	wordRepository "lingua-evo/internal/services/lingua/word/repository"
	wordService "lingua-evo/internal/services/lingua/word/service"
	userHandler "lingua-evo/internal/services/user/delivery"
	userRepository "lingua-evo/internal/services/user/repository"
	userService "lingua-evo/internal/services/user/service"

	signInHandler "lingua-evo/internal/services/site/auth/sign_in/delivery"
	signUpHandler "lingua-evo/internal/services/site/auth/sign_up/delivery"
	indexHandler "lingua-evo/internal/services/site/index/delivery"
)

const (
	filePath = "/website/"
	rootPath = "./../website"
)

func ServerStart(cfg *config.Config) {
	if cfg.PprofDebug.Enable {
		go func() {
			slog.Error("%v", http.ListenAndServe("localhost:6060", nil))
		}()
	}

	db, err := repository.NewDB(cfg.Database.GetConnStr())
	if err != nil {
		slog.Error(fmt.Errorf("can't create pg pool: %v", err).Error())
		return
	}

	router := mux.NewRouter()
	initServer(router, db)

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
	slog.Info("<----- start sertver ----->")
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

func initServer(r *mux.Router, db *sql.DB) {
	fs := http.FileServer(http.Dir(rootPath))
	r.PathPrefix(filePath).Handler(http.StripPrefix(filePath, fs))

	slog.Info("<----- create services ----->")
	slog.Info("user service")
	userRepo := userRepository.NewRepo(db)
	userSvc := userService.NewService(userRepo)

	slog.Info("word service")
	wordRepo := wordRepository.NewRepo(db)
	wordSvc := wordService.NewService(wordRepo)

	slog.Info("user service")
	langRepo := langRepository.NewRepo(db)
	langSvc := langService.NewService(langRepo)

	slog.Info("dictionary service")
	dictRepo := dictRepository.NewRepo(db)
	dictSvc := dictService.NewService(dictRepo)

	slog.Info("example service")
	exampleRepo := exampleRepository.NewRepo(db)
	exampleSvc := exampleService.NewService(exampleRepo)

	slog.Info("tag service")
	tagRepo := tagRepository.NewRepo(db)
	tagSvc := tagService.NewService(tagRepo)

	slog.Info("vocabulary service")
	vocabularyRepo := vocabularyRepository.NewRepo(db)
	vocabularySvc := vocabularyService.NewService(vocabularyRepo, wordSvc, exampleSvc, tagSvc)

	authRepo := authRepository.NewRepo(db)
	authSvc := authService.NewService(authRepo, userSvc)

	slog.Info("<----- create handlers ----->")
	slog.Info("index handler")
	indexHandler.Create(r, wordSvc)

	slog.Info("user handler")
	userHandler.Create(r, userSvc)

	slog.Info("sign_in handler")
	signInHandler.Create(r, userSvc)

	slog.Info("sign_up handler")
	signUpHandler.Create(r)

	slog.Info("account handler")
	accountHandler.Create(r)

	slog.Info("language handler")
	languageHandler.Create(r, langSvc)

	slog.Info("word handler")
	wordHandler.Create(r, wordSvc, langSvc)

	slog.Info("dictionary handler")
	dictHandler.Create(r, dictSvc)

	slog.Info("vocabulary handler")
	vocabularyHandler.Create(r, vocabularySvc)

	slog.Info("auth handler")
	authHandler.Create(r, authSvc)

	slog.Info("<----- end init services ----->")
}
