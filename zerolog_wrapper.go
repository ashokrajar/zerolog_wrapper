// Package zerolog_wrapper provides a custom log wrapper based on github.com/rs/zerolog
//
// How to use:
//
//	import (
//	    log "github.com/ashokrajar/zerolog_wrapper"
//	)
//
//	func init() {
//	    log.InitLog(log.TraceLevel, "dev")
//	}
//
//	func main() {
//	    log.Info().Msg("hello world")
//	}
//	// Output: {"time":1494567715,"level":"info","message":"hello world"}
//
// Fields can be added to log messages:
//
//	log.Info().Str("foo", "bar").Msg("hello world")
//	// Output: {"time":1494567715,"level":"info","message":"hello world","foo":"bar"}
//
// # Updating the logger context
//
//	log.UpdateContext(func(c zerolog.Context) zerolog.Context {
//				return c.Str("some_default_key", "some_default_value")
//			})
package zerolog_wrapper

import (
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type LogLevel string

const (
	TraceLevel LogLevel = "trace"
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
	FatalLevel LogLevel = "fatal"
	PanicLevel LogLevel = "panic"
)

type Env string

const (
	Prod  Env = "prod"
	Stage Env = "stage"
	QA    Env = "qa"
	Dev   Env = "dev"
)

var once sync.Once

var Logger zerolog.Logger

// Get local address of the running system
func getLocalIP() net.IP {
	conn, err := net.Dial("udp", "1.1.1.1:53")
	if err != nil {
		Logger.Fatal().Err(err)
	}
	defer conn.Close()

	return conn.LocalAddr().(*net.UDPAddr).IP
}

// InitLog initializes a global logger
func InitLog(logLevelStr LogLevel, appEnv Env) {
	once.Do(func() {
		var logLevel zerolog.Level

		switch logLevelStr {
		case TraceLevel:
			logLevel = zerolog.TraceLevel
		case DebugLevel:
			logLevel = zerolog.DebugLevel
		case InfoLevel:
			logLevel = zerolog.InfoLevel
		case WarnLevel:
			logLevel = zerolog.WarnLevel
		case ErrorLevel:
			logLevel = zerolog.ErrorLevel
		case FatalLevel:
			logLevel = zerolog.FatalLevel
		case PanicLevel:
			logLevel = zerolog.PanicLevel
		default:
			logLevel = zerolog.InfoLevel // default to INFO
		}

		output := zerolog.MultiLevelWriter(os.Stderr)

		// enforce TRACE and console output in development environment
		if appEnv == Dev {
			var consoleOutput io.Writer = zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.RFC3339,
			}
			logLevel = zerolog.TraceLevel
			output = zerolog.MultiLevelWriter(consoleOutput)
		}

		// Shorter file name in caller field
		zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
			curDir, _ := os.Getwd()
			shortPath := strings.TrimPrefix(file, curDir+"/")

			return shortPath + ":" + strconv.Itoa(line)
		}

		Logger = zerolog.New(output).
			Level(logLevel).
			With().
			Timestamp().
			IPAddr("host_ip", getLocalIP()).
			Caller().
			Logger()
		Logger.With().Caller()
	})
}

// UpdateContext is a function that updates the internal logger's context.
//
// Parameters:
// update: A function taking a zerolog.Context as input and then returns a zerolog.Context.
//
// eg:
//
//	log.UpdateContext(func(c zerolog.Context) zerolog.Context {
//		return c.Str("some_default_key", "some_default_value")
//	})
func UpdateContext(update func(c zerolog.Context) zerolog.Context) {
	Logger.UpdateContext(update)
}

// Trace starts a new message with trace level.
//
// You must call Msg on the returned event in order to send the event.
func Trace() *zerolog.Event {
	return Logger.Trace()
}

// Debug starts a new message with debug level.
//
// You must call Msg on the returned event in order to send the event.
func Debug() *zerolog.Event {
	return Logger.Debug()
}

// Info starts a new message with info level.
//
// You must call Msg on the returned event in order to send the event.
func Info() *zerolog.Event {
	return Logger.Info()
}

// Warn starts a new message with warn level.
//
// You must call Msg on the returned event in order to send the event.
func Warn() *zerolog.Event {
	return Logger.Warn()
}

// Error starts a new message with error level.
//
// You must call Msg on the returned event in order to send the event.
func Error() *zerolog.Event {
	return Logger.Error()
}

// Fatal starts a new message with fatal level.
//
// You must call Msg on the returned event in order to send the event.
func Fatal() *zerolog.Event {
	return Logger.Fatal()
}

// Panic starts a new message with panic level.
//
// You must call Msg on the returned event in order to send the event.
func Panic() *zerolog.Event {
	return Logger.Panic()
}
