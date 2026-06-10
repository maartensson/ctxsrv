package jlog

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/coreos/go-systemd/v22/journal"
)

type handler struct {
	attrs   map[string]string
	prefix  string
	options Options
}

type Options struct {
	Level      slog.Level
	WithSource bool
}

func New(options *Options) slog.Handler {
	if options == nil {
		options = &Options{
			Level:      slog.LevelInfo,
			WithSource: false,
		}
	}
	return &handler{
		attrs:  make(map[string]string),
		prefix: "",
		options: Options{
			Level:      options.Level,
			WithSource: options.WithSource,
		},
	}
}

func (h *handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.options.Level
}

func mapLevel(level slog.Level) journal.Priority {
	switch level {
	case slog.LevelDebug:
		return journal.PriDebug
	case slog.LevelInfo:
		return journal.PriInfo
	case slog.LevelWarn:
		return journal.PriWarning
	case slog.LevelError:
		return journal.PriErr
	default:
		return journal.PriInfo
	}
}

func (h *handler) Handle(ctx context.Context, r slog.Record) error {
	priority := mapLevel(r.Level)

	fields := maps.Clone(h.attrs)
	r.Attrs(func(a slog.Attr) bool {
		flattenAttrs(h.prefix, fields, a)
		return true
	})

	if h.options.WithSource && r.PC != 0 {
		fn := runtime.FuncForPC(r.PC)
		if fn != nil {
			file, line := fn.FileLine(r.PC)
			fields["SOURCE"] = fmt.Sprintf("%s:%d", filepath.Base(file), line)
		}
	}

	return journal.Send(r.Message, priority, fields)
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	out := maps.Clone(h.attrs)
	flattenAttrs(h.prefix, out, attrs...)
	return &handler{
		attrs:   out,
		prefix:  h.prefix,
		options: h.options,
	}
}

func (h *handler) WithGroup(name string) slog.Handler {
	return &handler{
		attrs:   maps.Clone(h.attrs),
		prefix:  addPrefix(h.prefix, name),
		options: h.options,
	}
}

func flattenAttrs(prefix string, out map[string]string, attrs ...slog.Attr) {
	for _, a := range attrs {
		key := strings.ToUpper(addPrefix(prefix, a.Key))
		switch a.Value.Kind() {
		case slog.KindLogValuer:
			flattenAttrs(prefix, out, slog.Attr{
				Key:   a.Key,
				Value: a.Value.LogValuer().LogValue(),
			})
		case slog.KindGroup:
			flattenAttrs(key, out, a.Value.Group()...)
		case slog.KindBool:
			out[key] = strconv.FormatBool(a.Value.Bool())
		case slog.KindDuration:
			out[key] = a.Value.Duration().String()
		case slog.KindFloat64:
			out[key] = strconv.FormatFloat(a.Value.Float64(), 'g', -1, 64)
		case slog.KindInt64:
			out[key] = strconv.FormatInt(a.Value.Int64(), 10)
		case slog.KindString:
			out[key] = a.Value.String()
		case slog.KindTime:
			out[key] = a.Value.Time().Format(time.RFC3339Nano)
		case slog.KindUint64:
			out[key] = strconv.FormatUint(a.Value.Uint64(), 10)
		case slog.KindAny:
			out[key] = fmt.Sprint(a.Value.Any())
		default:
			out[key] = fmt.Sprint(a.Value)
		}
	}
}

func addPrefix(prefix, key string) string {
	if prefix != "" {
		return prefix + "_" + key
	}
	return key
}
