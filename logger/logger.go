package logger

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var log *zap.Logger

type LogEntryCtx struct {
	Referer   zap.Field
	RequestId zapcore.Field
	Protocol  zap.Field
	Ip        zap.Field
	Ips       zap.Field
	Host      zap.Field
	Method    zap.Field
	Path      zap.Field
	Url       zap.Field
	Route     zap.Field
	Source    zap.Field
	ExtraInfo zapcore.Field
}

type LogEntryWithoutCtx struct {
	Source    zap.Field
	ExtraInfo zapcore.Field
}

type LogReqResp struct {
	RequestBody  interface{} `json:"requestBody"`
	ResponseData interface{} `json:"responseData"`
}

func getLogEntry(c *fiber.Ctx, extraInfo interface{}) LogEntryCtx {
	logEntry := LogEntryCtx{
		Referer:   zap.String("referer", c.Get(fiber.HeaderReferer)),
		RequestId: zap.Any("requestId", c.Locals("requestId")),
		Protocol:  zap.String("protocol", c.Protocol()),
		Ip:        zap.String("ip", c.IP()),
		Ips:       zap.Strings("ips", c.IPs()),
		Host:      zap.String("host", c.Hostname()),
		Method:    zap.String("method", c.Method()),
		Path:      zap.String("path", c.Path()),
		Url:       zap.String("url", c.OriginalURL()),
		Route:     zap.String("route", c.Route().Path),
		Source:    zap.String("source", "fiber middleware/controllers"),
		ExtraInfo: zap.Any("ExtraInfo", extraInfo),
	}
	return logEntry
}

func getLogEntryWithoutContext(extraInfo interface{}) LogEntryWithoutCtx {
	logEntry := LogEntryWithoutCtx{
		Source:    zap.String("source", "Async Processes like go routines"),
		ExtraInfo: zap.Any("ExtraInfo", extraInfo),
	}
	return logEntry
}

func getLogWriter(hostname string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename: "/var/log/" + hostname + "/multiBot.log",
		Compress: false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func LoggerInit() error {
	cfg := zap.NewProductionEncoderConfig()

	cfg.TimeKey = "time"
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder

	cfg.LevelKey = "level"
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder

	cfg.MessageKey = "message"

	fileEncoder := zapcore.NewJSONEncoder(cfg)
	consoleEncoder := zapcore.NewConsoleEncoder(cfg)

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	writer := getLogWriter(hostname)
	defaultLogLevel := zapcore.DebugLevel

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)

	log = zap.New(core)
	return nil
}

func LogInfo(c *fiber.Ctx, logMsg string, extraInfo interface{}) {
	if c != nil {
		logEntry := getLogEntry(c, extraInfo)
		log.Info(
			logMsg,
			logEntry.Referer,
			logEntry.RequestId,
			logEntry.Protocol,
			logEntry.Ip,
			logEntry.Ips,
			logEntry.Host,
			logEntry.Method,
			logEntry.Path,
			logEntry.Url,
			logEntry.Route,
			logEntry.Source,
			logEntry.ExtraInfo,
		)
	} else {
		logEntryNoCtx := getLogEntryWithoutContext(extraInfo)
		log.Info(
			logMsg,
			logEntryNoCtx.Source,
			logEntryNoCtx.ExtraInfo,
		)
	}
}

func LogError(c *fiber.Ctx, logMsg string, extraInfo interface{}) {
	if c != nil {
		logEntry := getLogEntry(c, extraInfo)
		log.Error(
			logMsg,
			logEntry.Referer,
			logEntry.RequestId,
			logEntry.Protocol,
			logEntry.Ip,
			logEntry.Ips,
			logEntry.Host,
			logEntry.Method,
			logEntry.Path,
			logEntry.Url,
			logEntry.Route,
			logEntry.Source,
			logEntry.ExtraInfo,
		)
	} else {
		logEntryNoCtx := getLogEntryWithoutContext(extraInfo)
		log.Error(
			logMsg,
			logEntryNoCtx.Source,
			logEntryNoCtx.ExtraInfo,
		)
	}
}

func LogDebug(c *fiber.Ctx, logMsg string, extraInfo interface{}) {
	if c != nil {
		logEntry := getLogEntry(c, extraInfo)
		log.Debug(
			logMsg,
			logEntry.Referer,
			logEntry.RequestId,
			logEntry.Protocol,
			logEntry.Ip,
			logEntry.Ips,
			logEntry.Host,
			logEntry.Method,
			logEntry.Path,
			logEntry.Url,
			logEntry.Route,
			logEntry.Source,
			logEntry.ExtraInfo,
		)
	} else {
		logEntryNoCtx := getLogEntryWithoutContext(extraInfo)
		log.Debug(
			logMsg,
			logEntryNoCtx.Source,
			logEntryNoCtx.ExtraInfo,
		)
	}
}

func LogWarn(c *fiber.Ctx, logMsg string, extraInfo interface{}) {
	if c != nil {
		logEntry := getLogEntry(c, extraInfo)
		log.Warn(
			logMsg,
			logEntry.Referer,
			logEntry.RequestId,
			logEntry.Protocol,
			logEntry.Ip,
			logEntry.Ips,
			logEntry.Host,
			logEntry.Method,
			logEntry.Path,
			logEntry.Url,
			logEntry.Route,
			logEntry.Source,
			logEntry.ExtraInfo,
		)
	} else {
		logEntryNoCtx := getLogEntryWithoutContext(extraInfo)
		log.Warn(
			logMsg,
			logEntryNoCtx.Source,
			logEntryNoCtx.ExtraInfo,
		)
	}
}

func LogPanic(c *fiber.Ctx, logMsg string, extraInfo interface{}) {
	if c != nil {
		logEntry := getLogEntry(c, extraInfo)
		log.Panic(
			logMsg,
			logEntry.Referer,
			logEntry.RequestId,
			logEntry.Protocol,
			logEntry.Ip,
			logEntry.Ips,
			logEntry.Host,
			logEntry.Method,
			logEntry.Path,
			logEntry.Url,
			logEntry.Route,
			logEntry.Source,
			logEntry.ExtraInfo,
		)
	} else {
		logEntryNoCtx := getLogEntryWithoutContext(extraInfo)
		log.Panic(
			logMsg,
			logEntryNoCtx.Source,
			logEntryNoCtx.ExtraInfo,
		)
	}
}

func LogFatal(c *fiber.Ctx, logMsg string, extraInfo interface{}) {
	if c != nil {
		logEntry := getLogEntry(c, extraInfo)
		log.Fatal(
			logMsg,
			logEntry.Referer,
			logEntry.RequestId,
			logEntry.Protocol,
			logEntry.Ip,
			logEntry.Ips,
			logEntry.Host,
			logEntry.Method,
			logEntry.Path,
			logEntry.Url,
			logEntry.Route,
			logEntry.Source,
			logEntry.ExtraInfo,
		)
	} else {
		logEntryNoCtx := getLogEntryWithoutContext(extraInfo)
		log.Fatal(
			logMsg,
			logEntryNoCtx.Source,
			logEntryNoCtx.ExtraInfo,
		)
	}
}
