package logx

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLogxStdout(t *testing.T) {
	// 默认 Dir 为空，应输出到 stdout，不会报错
	l := NewLogger(&Option{
		Dir:   "",
		Level: Debug,
	})

	l.Debug("test debug to stdout")
	l.Infof("test info to stdout: %d", 123)
	l.Warn("test warn to stdout")
	l.Errorf("test error to stdout")
}

func TestLogxFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "logx_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	l := NewLogger(&Option{
		Dir:        tmpDir,
		Level:      Info,
		Compress:   false,
		ReserveDay: 1,
	})

	l.Debug("this should not be printed")
	l.Info("this is info log")
	l.Errorf("this is error log: %s", "err")

	files, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("failed to read temp dir: %v", err)
	}
	if len(files) == 0 {
		t.Fatalf("expected log files created in %s, got none", tmpDir)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, files[0].Name()))
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}
	t.Logf("file log output:\n%s", string(content))
}
