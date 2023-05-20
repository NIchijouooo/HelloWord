package setting

import (
	"encoding/json"
	"fmt"
	"gateway/device/eventBus"
	"io"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type WSWriteTemplate struct {
}

var ZAPS *zap.SugaredLogger

var ZapEventBus = eventBus.NewBus()

var ZapWrite WSWriteTemplate

const (
	// DefaultLevel the default log level
	DefaultLevel = zapcore.InfoLevel

	// DefaultTimeLayout the default time layout;
	DefaultTimeLayout = time.RFC3339
)

func InitLogger() {
	var zapLogger *zap.Logger
	var err error

	if LogToFile == true {
		zapLogger, err = NewJSONLogger(WithTimeLayout("2006-01-02 15:04:05.000"), WithFileRotationP(LogFile))
	} else {
		zapLogger, err = NewJSONLogger(WithTimeLayout("2006-01-02 15:04:05.000"))
	}
	if err != nil {
		fmt.Printf("zap日志模块初始化失败 %v", err)
		panic(err)
	}

	ZAPS = zapLogger.Sugar()
	ZAPS.Infof("zap日志模块初始化成功!日志等级:%s", LogLevel)
}

func Level(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return DefaultLevel
	}
}

// Option custom setup config
type Option func(*option)

type option struct {
	level          zapcore.Level
	fields         map[string]string
	file           io.Writer
	timeLayout     string
	disableConsole bool
}

// WithDebugLevel only greater than 'level' will output
func WithDebugLevel() Option {
	return func(opt *option) {
		opt.level = zapcore.DebugLevel
	}
}

// WithInfoLevel only greater than 'level' will output
func WithInfoLevel() Option {
	return func(opt *option) {
		opt.level = zapcore.InfoLevel
	}
}

// WithWarnLevel only greater than 'level' will output
func WithWarnLevel() Option {
	return func(opt *option) {
		opt.level = zapcore.WarnLevel
	}
}

// WithErrorLevel only greater than 'level' will output
func WithErrorLevel() Option {
	return func(opt *option) {
		opt.level = zapcore.ErrorLevel
	}
}

// WithField add some field(s) to log
func WithField(key, value string) Option {
	return func(opt *option) {
		opt.fields[key] = value
	}
}

// WithFileP write log to some file
func WithFileP(file string) Option {
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, 0766); err != nil {
		panic(err)
	}

	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0766)
	if err != nil {
		panic(err)
	}

	return func(opt *option) {
		opt.file = zapcore.Lock(f)
	}
}

// WithFileRotationP write log to some file with rotation
func WithFileRotationP(file string) Option {
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, 0766); err != nil {
		panic(err)
	}

	return func(opt *option) {
		opt.file = &lumberjack.Logger{ // concurrent-safed
			Filename:   file,           // 文件路径
			MaxSize:    LogFileMaxSize, // 单个文件最大尺寸，默认单位 M
			MaxBackups: LogFileBackup,  // 最多保留3个备份
			MaxAge:     7,              // 日志文件最多保存天数
			LocalTime:  true,           // 使用本地时间
			Compress:   true,           // 是否压缩
		}
	}
}

// WithTimeLayout custom time format
func WithTimeLayout(timeLayout string) Option {
	return func(opt *option) {
		opt.timeLayout = timeLayout
	}
}

// WithEnableConsole write log to os.Stdout or os.Stderr
func WithDisableConsole() Option {
	return func(opt *option) {
		opt.disableConsole = true
	}
}

// NewJSONLogger return a json-encoder zap logger,
func NewJSONLogger(opts ...Option) (*zap.Logger, error) {

	opt := &option{
		level:          Level(LogLevel),
		fields:         make(map[string]string),
		disableConsole: false,
	}
	for _, f := range opts {
		f(opt)
	}

	timeLayout := DefaultTimeLayout
	if opt.timeLayout != "" {
		timeLayout = opt.timeLayout
	}

	// similar to zap.NewProductionEncoderConfig()
	fileConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger", // used by logger.Named(key); optional; useless
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace", // use by zap.AddStacktrace; optional; useless
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(timeLayout))
		},
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 全路径编码器
	}

	webConfig := zapcore.EncoderConfig{
		//TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger", // used by logger.Named(key); optional; useless
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace", // use by zap.AddStacktrace; optional; useless
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder, // 小写编码器
		//EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder, // 全路径编码器
	}

	consoleConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger", // used by logger.Named(key); optional; useless
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace", // use by zap.AddStacktrace; optional; useless
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(timeLayout))
		},
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 全路径编码器
	}

	fileEncoder := zapcore.NewJSONEncoder(fileConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(consoleConfig)
	webEncoder := zapcore.NewJSONEncoder(webConfig)

	/* lowPriority usd by info\debug\warn */
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= opt.level && lvl <= zapcore.ErrorLevel
	})

	// highPriority usd by error\panic\fatal
	//highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
	//	return lvl >= opt.level && lvl >= zapcore.ErrorLevel
	//})

	stdout := zapcore.Lock(os.Stdout) // lock for concurrent safe
	//stderr := zapcore.Lock(os.Stderr) // lock for concurrent safe

	core := zapcore.NewTee()
	if !opt.disableConsole {
		core = zapcore.NewTee(
			zapcore.NewCore(consoleEncoder,
				zapcore.NewMultiWriteSyncer(stdout),
				lowPriority,
			),
			zapcore.NewCore(webEncoder,
				zapcore.AddSync(&ZapWrite),
				lowPriority,
			),
		)
	}
	if opt.file != nil {
		core = zapcore.NewTee(core,
			zapcore.NewCore(fileEncoder,
				zapcore.AddSync(opt.file),
				zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
					return lvl >= opt.level
				}),
			),
		)
	}

	zapLog := zap.New(core,
		zap.AddCaller(),
		//zap.ErrorOutput(stderr),
	)

	for key, value := range opt.fields {
		zapLog = zapLog.WithOptions(zap.Fields(zapcore.Field{
			Key:    key,
			Type:   zapcore.StringType,
			String: value,
		}))
	}
	return zapLog, nil
}

func (w *WSWriteTemplate) Write(p []byte) (n int, err error) {

	logParam := struct {
		Level  string `json:"level"`
		Caller string `json:"caller"`
		Msg    string `json:"msg"`
	}{}

	_ = json.Unmarshal(p, &logParam)

	msg := struct {
		Content   string `json:"content"`
		TimeStamp string `json:"timeStamp"` //时间戳
		Direction int    `json:"direction"` //数据方向
		Label     string `json:"label"`     //数据标识
	}{
		Content:   "[" + logParam.Caller + "]    " + logParam.Msg,
		TimeStamp: time.Now().Format("2006-01-02 15:04:05.000"),
		Direction: 0,
		Label:     logParam.Level,
	}

	_ = ZapEventBus.Publish("zap", msg)

	return len(p), nil
}
