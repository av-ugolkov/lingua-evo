package log

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"time"

	cfgLog "github.com/av-ugolkov/lingua-evo/internal/config"
)

type Logger struct {
	file         *os.File
	Log          *slog.Logger
	ServerLogger *log.Logger
}

func CustomLogger(cfgLog *cfgLog.Logger) *Logger {
	writers := make([]io.Writer, 0, 2)

	var err error
	var file *os.File
	for _, writer := range cfgLog.Output {
		switch writer {
		case "file":
			err = os.Mkdir("logs", 0750)
			if err != nil {
				slog.Warn(fmt.Sprintf("pgk.log.CustomLogger: %v", err))
			}
			file, err = os.OpenFile(fmt.Sprintf("logs/log_file_%s.log", time.Now().Format("2006_01_02_15_04_05")), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
			if err != nil {
				slog.Error(fmt.Sprintf("pgk.log.CustomLogger: %v", err))
				return nil
			}
			writers = append(writers, file)
		case "console":
			writers = append(writers, os.Stdout)
		default:
			slog.Warn("pgk.log.CustomLogger: unknown writer")
		}
	}
	multiWriter := io.MultiWriter(writers...)

	logHandler := NewCustomHandler(multiWriter, &slog.HandlerOptions{
		Level: getLevel(cfgLog.Level),
	})

	return &Logger{
		file:         file,
		Log:          slog.New(logHandler),
		ServerLogger: slog.NewLogLogger(logHandler, getLevel(cfgLog.Level)),
	}
}

func (l *Logger) Close() {
	err := l.file.Close()
	if err != nil {
		slog.Error(fmt.Sprintf("pgk.Logger.Close: %v", err))
	}
}

func getLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		slog.Warn("pgk.log.getLevel: unknonw level")
		return slog.LevelInfo
	}
}
