package logging

import (
	"context"
	"testing"
)

func setup(t *testing.T) Logger {
	t.Helper()
	logger, err := New(ConsoleHandler, Options{
		Development:  true,
		Debug:        true,
		EnableCaller: true,
		SkipCaller:   3,
	})
	if err != nil {
		t.Fatal(err)
	}

	return logger
}

func TestLog_Error(t *testing.T) {
	logger := setup(t)
	logger.Error(false, "error example", "test", 2, "test2", 2.5)
}

func TestLog_Debug(t *testing.T) {
	logger := setup(t)
	logger.Debug(false, "test", "msg", "hello")
}

func TestLog_Info(t *testing.T) {
	logger := setup(t)
	logger.Info(false, "error example", "test", 2, "test2", 2.5)
}

func TestLog_Warn(t *testing.T) {
	logger := setup(t)
	logger.Warn(false, "error example", "test", 2, "test2", 2.5)
}

func TestLog_ErrorContext(t *testing.T) {
	logger := setup(t)
	logger.ErrorContext(context.TODO(), false, "error example", "test", 2, "test2", 2.5)
}

func TestLog_DebugContext(t *testing.T) {
	logger := setup(t)
	logger.DebugContext(context.TODO(), false, "error example", "test", 2, "test2", 2.5)
}

func TestLog_InfoContext(t *testing.T) {
	logger := setup(t)
	logger.InfoContext(context.TODO(), false, "error example", "test", 2, "test2", 2.5)
}

func TestLog_WarnContext(t *testing.T) {
	logger := setup(t)
	logger.WarnContext(context.TODO(), false, "error example", "test", 2, "test2", 2.5)
}
