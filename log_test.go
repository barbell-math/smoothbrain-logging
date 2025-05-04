package sblog

import (
	"context"
	"testing"
)

func TestExampleOutput(t *testing.T) {
	if !testing.Verbose() {
		return
	}
	Example_logDebug()
	Example_logInfo()
	Example_logWarn()
	Example_logError()
	Example_logVerbose()
	Example_persistantLogs()
}

func Example_logDebug() {
	s, _ := New(Opts{})
	s.Debug("This is a debug message")
	s.Debug(
		"This is a debug message",
		"Keys with nil values will be treated as messages", nil,
		"and will be printed on a separate line", nil,
	)
	s.Debug("This is a debug message logging a value", "value", 10)
	s.Debug(
		"This is a debug message logging a multiple values",
		"value1", 10,
		"value2", 11,
	)
	s.Debug("This is a debug message", "Separate line", nil, "value", 10)
}

func Example_logInfo() {
	s, _ := New(Opts{})
	s.Info("This is an info message")
	s.Info(
		"This is an info message",
		"Keys with nil values will be treated as messages", nil,
		"and will be printed on a separate line", nil,
	)
	s.Info("This is an info message logging a value", "value", 10)
	s.Info(
		"This is an info message logging a multiple values",
		"value1", 10,
		"value2", 11,
	)
	s.Info("This is an info message", "Separate line", nil, "value", 10)
}

func Example_logWarn() {
	s, _ := New(Opts{})
	s.Warn("This is a warn message")
	s.Warn(
		"This is a warn message",
		"Keys with nil values will be treated as messages", nil,
		"and will be printed on a separate line", nil,
	)
	s.Warn("This is a warn message logging a value", "value", 10)
	s.Warn(
		"This is a warn message logging a multiple values",
		"value1", 10,
		"value2", 11,
	)
	s.Warn("This is a warn message", "Separate line", nil, "value", 10)
}

func Example_logError() {
	s, _ := New(Opts{})
	s.Error("This is an error message")
	s.Error(
		"This is an error message",
		"Keys with nil values will be treated as messages", nil,
		"and will be printed on a separate line", nil,
	)
	s.Error("This is an error message logging a value", "value", 10)
	s.Error(
		"This is an error message logging a multiple values",
		"value1", 10,
		"value2", 11,
	)
	s.Error("This is an error message", "Separate line", nil, "value", 10)
}

func Example_logVerbose() {
	s, _ := New(Opts{CurVerbosityLevel: 2})
	s.Log(
		context.TODO(),
		VLevel(0),
		"This is a level 0 verbose message, a.k.a. a debug message",
	)
	s.Log(context.TODO(), VLevel(1), "This is a level 1 verbose message")
	s.Log(context.TODO(), VLevel(2), "This is a level 2 verbose message")
	// This next line should not print
	s.Log(context.TODO(), VLevel(3), "This is a level 3 verbose message")

	s, _ = New(Opts{CurVerbosityLevel: 0})
	s.Log(
		context.TODO(),
		VLevel(0),
		"This is a level 0 verbose message, a.k.a. a debug message",
	)
	// These next three lines should not print
	s.Log(context.TODO(), VLevel(1), "This is a level 1 verbose message")
	s.Log(context.TODO(), VLevel(2), "This is a level 2 verbose message")
	s.Log(context.TODO(), VLevel(3), "This is a level 3 verbose message")
}

func Example_persistantLogs() {
	s, _ := New(Opts{
		CurVerbosityLevel: 1,
		RotateWriterOpts: RotateWriterOpts{
			LogDir:  "./testData/",
			LogName: "persistantExample",
		},
	})
	s.Info("This should be saved to ./testData/persistantExample.0.log")
	s.Warn("This should be saved to ./testData/persistantExample.0.log")
	s.Error("This should be saved to ./testData/persistantExample.0.log")
	s.Debug("This should be saved to ./testData/persistantExample.0.log")
	s.Log(
		context.TODO(), VLevel(1),
		"This should be saved to ./testData/persistantExample.0.log",
	)
}
