package Log

import (
	"ayachanV2/Config"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	"time"
)

var Log = logrus.New()

func init() {
	var logLevel = logrus.WarnLevel
	if Config.Config.Debug {
		logLevel = logrus.DebugLevel
	}
	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   "logs/warnings.log",
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
		Level:      logrus.WarnLevel,
		Formatter: &logrus.TextFormatter{
			DisableColors:   true,
			TimestampFormat: time.RFC3339,
		},
	})

	rotateFileHook2, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   "logs/console.log",
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
		Level:      logrus.InfoLevel,
		Formatter: &logrus.TextFormatter{
			DisableColors:   true,
			TimestampFormat: time.RFC3339,
		},
	})

	if err != nil {
		Log.Fatalf("Failed to initialize file rotate hook: %v", err)
	}

	Log.SetLevel(logLevel)
	Log.SetOutput(colorable.NewColorableStdout())
	Log.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
	})
	Log.AddHook(rotateFileHook)
	Log.AddHook(rotateFileHook2)
}
