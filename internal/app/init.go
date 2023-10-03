package app

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/gorilla/mux"

	"lingua-evo/internal/config"
	repository "lingua-evo/internal/db"
	accountHandler "lingua-evo/internal/services/account/delivery"
	indexHandler "lingua-evo/internal/services/index/delivery"
	langRepository "lingua-evo/internal/services/language/repository"
	langService "lingua-evo/internal/services/language/service"
	signInHandler "lingua-evo/internal/services/sign_in/delivery"
	signUpHandler "lingua-evo/internal/services/sign_up/delivery"
	userRepository "lingua-evo/internal/services/user/repository"
	userService "lingua-evo/internal/services/user/service"
	wordHandler "lingua-evo/internal/services/word/delivery"
	wordRepository "lingua-evo/internal/services/word/repository"
	wordService "lingua-evo/internal/services/word/service"
)

const (
	filePath = "/static/"
	rootPath = "./../static"
)

func ServerStart(cfg *config.Config) {
	if cfg.PprofDebug.Enable {
		go func() {
			slog.Error("%v", http.ListenAndServe("localhost:6060", nil))
		}()
	}

	db, err := repository.NewDB(cfg.Database.GetConnStr())
	if err != nil {
		slog.Error("can't create pg pool: %v", err)
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
	slog.Info("web address: %v", listener.Addr())

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

	slog.Info("<----- create handlers ----->")
	slog.Info("index")
	indexHandler.Create(r)

	slog.Info("sign_in")
	signInHandler.Create(r, userSvc)

	slog.Info("sign_up")
	signUpHandler.Create(r, userSvc)

	slog.Info("account")
	accountHandler.Create(r)

	slog.Info("word")
	wordHandler.Create(r, wordSvc, langSvc)

	slog.Info("<----- end init services ----->")
}
