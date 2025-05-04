// A very simple library that implements a logger with verbosity levels and an
// optional rotating file log writer.
package sblog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"
)

type (
	// The custom log handler that formats all of the log messages, enforces
	// verbosity levels, and handles writing to both stdout/stderr and the
	// rotating writer.
	handler struct {
		curVerbosityLevel slog.Level
		timeFmt           string
		stdoutMultiWriter io.Writer
		stderrMultiWriter io.Writer
	}

	Opts struct {
		// Sets the allowed maximum verbosity level. There is no bound to the
		// requested verbosity, but any verbosity levels greater than the level
		// set here will not be printed.
		CurVerbosityLevel uint
		// The time format to use when printing log messages.
		TimeFmt string
		RotateWriterOpts
	}
)

// Translates the requested positive verbosity level from a range of [0,inf) to
// (-inf, -4]. This must be done for the level to be understood by the [slog]
// package.
func VLevel(val uint) slog.Level {
	return slog.Level(-int(val) + int(slog.LevelDebug))
}

// Creates a new logger. The logs will always be printed to stdout/stderr. If
// `opts.RotateWriterOpts.LogDir` is an empty string no logs will be written to
// disk. If `opts.RotateWriterOpts.LogDir` is a valid dir then all logs from
// stdout and stderr will be mirrored to files in the supplied dir using a
// [RotateWriter].
func New(opts Opts) (*slog.Logger, error) {
	var err error
	var persistentLogs *RotateWriter

	if opts.LogDir != "" {
		if persistentLogs, err = NewRotateWriter(
			opts.RotateWriterOpts,
		); err != nil {
			persistentLogs.Close()
			return nil, err
		}
	}

	if opts.TimeFmt == "" {
		opts.TimeFmt = time.StampNano
	}
	handler := newHandler(&opts, persistentLogs)
	return slog.New(&handler), nil
}

func newHandler(
	opts *Opts,
	persistentLogs *RotateWriter,
) handler {
	rv := handler{
		curVerbosityLevel: VLevel(opts.CurVerbosityLevel),
		timeFmt:           opts.TimeFmt,
		stdoutMultiWriter: io.MultiWriter(),
	}

	if persistentLogs != nil {
		rv.stdoutMultiWriter = io.MultiWriter(os.Stdout, persistentLogs)
	} else {
		rv.stdoutMultiWriter = io.MultiWriter(os.Stdout)
	}

	if persistentLogs != nil {
		rv.stderrMultiWriter = io.MultiWriter(os.Stderr, persistentLogs)
	} else {
		rv.stderrMultiWriter = io.MultiWriter(os.Stderr)
	}

	return rv
}

// Returns true if the supplied logging level is enabled. Error, warn, info, and
// debug log levels will always return true. Verbosity levels will return true
// if they are within the allowed max verbosity.
func (h *handler) Enabled(ctxt context.Context, level slog.Level) bool {
	switch level {
	case slog.LevelError, slog.LevelWarn, slog.LevelInfo, slog.LevelDebug:
		return true
	default:
		if level < slog.LevelDebug {
			return level >= h.curVerbosityLevel
		}
		return false
	}
}

// Formats and writes the log messages.
func (h *handler) Handle(ctxt context.Context, record slog.Record) error {
	message := record.Message
	timestamp := record.Time.Format(h.timeFmt)

	zeroVal := slog.Value{}
	record.Attrs(func(attr slog.Attr) bool {
		if attr.Value.Equal(zeroVal) {
			message += fmt.Sprintf("\n\t→ %s", attr.Key)
		} else if record.NumAttrs() == 1 {
			message += fmt.Sprintf(
				" \u001b[90m%s=\u001b[0m%v",
				attr.Key, attr.Value.String(),
			)
		} else {
			message += fmt.Sprintf(
				"\n\t→ \u001b[90m%s=\u001b[0m%v",
				attr.Key, attr.Value.String(),
			)
		}
		return true
	})

	switch record.Level {
	case slog.LevelDebug:
		fmt.Fprintf(
			h.stdoutMultiWriter, "\u001b[35m[%v] (%v)\u001b[0m %v\n",
			record.Level, timestamp, message,
		)
	case slog.LevelInfo:
		fmt.Fprintf(
			h.stdoutMultiWriter, "\u001b[36m[%v] (%v)\u001b[0m %v\n",
			record.Level, timestamp, message,
		)
	case slog.LevelWarn:
		fmt.Fprintf(
			h.stdoutMultiWriter, "\u001b[33m[%v] (%v)\u001b[0m %v\n",
			record.Level, timestamp, message,
		)
	case slog.LevelError:
		fmt.Fprintf(
			h.stderrMultiWriter, "\u001b[31m[%v] (%v)\u001b[0m %v\n",
			record.Level, timestamp, message,
		)
	default:
		if record.Level < h.curVerbosityLevel {
			return nil
		}
		fmt.Fprintf(
			h.stdoutMultiWriter, "\u001b[90m[VERBOS][lvl %d/%d] (%v)\u001b[0m %v\n",
			-record.Level+slog.LevelDebug, -h.curVerbosityLevel+slog.LevelDebug,
			timestamp, message,
		)
	}

	return nil
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	panic("unimplemented")
}

func (h *handler) WithGroup(name string) slog.Handler {
	panic("unimplemented")
}
