package utility

import (
	"context"
	"fmt"
	"io"
	"log/slog"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
)

// PrettyHandler is a custom slog.Handler that formats logs in a human-readable way
// with colors, aligned columns, and local timestamps.
type PrettyHandler struct {
	opts slog.HandlerOptions
	out  io.Writer
}

// NewPrettyHandler creates a new PrettyHandler
func NewPrettyHandler(out io.Writer, opts *slog.HandlerOptions) *PrettyHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &PrettyHandler{
		opts: *opts,
		out:  out,
	}
}

func (h *PrettyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return level >= minLevel
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	// Format time with date (YY-MM-DD HH:MM:SS)
	timeStr := r.Time.Format("06-01-02 15:04:05")

	// Get level string and color
	var levelStr string
	var color string
	switch r.Level {
	case slog.LevelDebug:
		levelStr = "DEBUG"
		color = colorGray
	case slog.LevelInfo:
		levelStr = "INFO "
		color = colorBlue
	case slog.LevelWarn:
		levelStr = "WARN "
		color = colorYellow
	case slog.LevelError:
		levelStr = "ERROR"
		color = colorRed
	default:
		levelStr = r.Level.String()
		color = colorReset
	}

	// Build the base message with aligned columns
	buf := fmt.Sprintf("%s %s%s%s %s",
		timeStr,
		color,
		levelStr,
		colorReset,
		r.Message,
	)

	// Add attributes (context) with colors
	r.Attrs(func(a slog.Attr) bool {
		buf += fmt.Sprintf(" %s%s%s=%v", colorCyan, a.Key, colorReset, a.Value)
		return true
	})

	// Write to output
	buf += "\n"
	_, err := h.out.Write([]byte(buf))
	return err
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// For simplicity, return the same handler
	// A full implementation would need to track these attrs
	return h
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	// For simplicity, return the same handler
	// A full implementation would need to track the group name
	return h
}
