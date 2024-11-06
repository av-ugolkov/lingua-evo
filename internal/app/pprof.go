package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/config"
)

var pprofSrv *http.Server

func initPprof(pprofDebug *config.PprofDebug) {
	go func() {
		pprofSrv = &http.Server{
			Addr:              pprofDebug.Addr(),
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 10 * time.Second, //TODO created into config
		}
		slog.Error(pprofSrv.ListenAndServe().Error())
	}()
}

func shutdownPprof(ctx context.Context) {
	if pprofSrv != nil {
		if err := pprofSrv.Shutdown(ctx); err != nil {
			slog.Error(fmt.Sprintf("server pprof shutdown returned an err: %v\n", err))
		}
	}
}
