package gofiglet

import (
	"image/color"
	"strings"
	"testing"
)

func TestNewCmdBanner_DefaultPadding(t *testing.T) {
	b, err := NewCmdBanner([]string{"Test"}, WithColors([]color.Color{ColorCyan}))
	if err != nil {
		t.Fatalf("NewCmdBanner() error = %v", err)
	}
	if b.TopPadding != 1 {
		t.Errorf("default TopPadding = %d, want 1", b.TopPadding)
	}
	if b.BottomPadding != 0 {
		t.Errorf("default BottomPadding = %d, want 0", b.BottomPadding)
	}
}

func TestWithPadding(t *testing.T) {
	b, err := NewCmdBanner([]string{"Test"}, WithPadding(3, 2), WithColors([]color.Color{ColorCyan}))
	if err != nil {
		t.Fatalf("NewCmdBanner() error = %v", err)
	}
	if b.TopPadding != 3 {
		t.Errorf("TopPadding = %d, want 3", b.TopPadding)
	}
	if b.BottomPadding != 2 {
		t.Errorf("BottomPadding = %d, want 2", b.BottomPadding)
	}
}

func TestCmdBanner_TopPadding(t *testing.T) {
	b, err := NewCmdBanner([]string{"Test"}, WithPadding(2, 0), WithColors([]color.Color{ColorCyan}))
	if err != nil {
		t.Fatalf("NewCmdBanner() error = %v", err)
	}
	result, err := CmdBanner(b)
	if err != nil {
		t.Fatalf("CmdBanner() error = %v", err)
	}
	if !strings.HasPrefix(result, "\n\n") {
		t.Errorf("expected 2 leading newlines from TopPadding=2, got: %q", result[:min(20, len(result))])
	}
	if strings.HasPrefix(result, "\n\n\n") {
		t.Errorf("expected exactly 2 leading newlines from TopPadding=2, got 3")
	}
}

func TestCmdBanner_BottomPadding(t *testing.T) {
	b, err := NewCmdBanner([]string{"Test"}, WithPadding(0, 2), WithColors([]color.Color{ColorCyan}))
	if err != nil {
		t.Fatalf("NewCmdBanner() error = %v", err)
	}
	result, err := CmdBanner(b)
	if err != nil {
		t.Fatalf("CmdBanner() error = %v", err)
	}
	if !strings.HasSuffix(result, "\n\n\n") {
		t.Errorf("expected 3 trailing newlines (renderer \\n + BottomPadding=2), got: %q", result[max(0, len(result)-20):])
	}
}

func TestCmdBanner_NoPadding(t *testing.T) {
	b, err := NewCmdBanner([]string{"Test"}, WithPadding(0, 0), WithColors([]color.Color{ColorCyan}))
	if err != nil {
		t.Fatalf("NewCmdBanner() error = %v", err)
	}
	result, err := CmdBanner(b)
	if err != nil {
		t.Fatalf("CmdBanner() error = %v", err)
	}
	if strings.HasPrefix(result, "\n") {
		t.Errorf("expected no leading newline with TopPadding=0, got: %q", result[:min(20, len(result))])
	}
}
