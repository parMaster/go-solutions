package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
)

// Some notes on checking out the slog package
// Sources:
// - https://go.dev/blog/slog

func main() {

	logger := slog.Default()
	logger.Info("hello, world", "user", os.Getenv("USER"))

	// Output:
	// 2023/08/23 02:33:06 INFO hello, world user=macmini

	logger.Warn("longer values", "like this", "are allowed?")
	// Output:
	// 2023/08/23 02:49:26 WARN longer values "like this"="are allowed?"

	logger.Error("one parameter", "one")
	logger.Error("two parameters", "one", "two")
	logger.Error("three parameters", "one", "two", "three")
	logger.Error("four parameters", "one", "two", "three", "four")
	// 2023/08/23 02:55:23 ERROR one parameter !BADKEY=one
	// 2023/08/23 02:55:23 ERROR two parameters one=two
	// 2023/08/23 02:55:23 ERROR three parameters one=two !BADKEY=three
	// 2023/08/23 02:55:23 ERROR four parameters one=two three=four

	// Besides Info, there are functions for three other levels—Debug, Warn,
	// and Error—as well as a more general Log function that takes the level as an argument.

	logger.Info("slog.LevelDebug", "implements stringer", slog.LevelDebug)
	// Output:
	// 2023/08/23 03:03:48 INFO slog.LevelDebug "implements stringer?"=DEBUG

	l2 := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	l2.Info("hello, world", "user", os.Getenv("USER"))
	// {"time":"2023-08-23T03:06:53.709660371+03:00","level":"INFO","msg":"hello, world","user":"macmini"}

	l3 := slog.NewLogLogger(slog.NewJSONHandler(os.Stdout, nil), slog.LevelInfo)
	l3.Printf("hello, world %s, %s", "user", os.Getenv("USER"))
	// {"time":"2023-08-23T03:19:11.54195029+03:00","level":"INFO","msg":"hello, world user, macmini"}

	opts := &slog.HandlerOptions{Level: slog.LevelWarn}
	l4 := slog.New(slog.NewTextHandler(os.Stdout, opts))

	// Levels go: Debug, Info, Warn, Error. Setting the level to Warn means that
	// levels that on the left of Warn will be filtered out.
	l4.Debug("will be filtered out", "user", os.Getenv("USER"))
	l4.Info("will be filtered out", "user", os.Getenv("USER"))
	l4.Warn("will be printed", "user", os.Getenv("USER"))
	l4.Error("will be printed", "user", os.Getenv("USER"))
	// Output:
	// time=2023-08-23T03:59:47.334+03:00 level=WARN msg="will be printed" user=macmini
	// time=2023-08-23T03:59:47.334+03:00 level=ERROR msg="will be printed" user=macmini

	// The slog package provides a Handler interface that can be implemented to
	// write log records to a variety of destinations. The slog package provides
	// several implementations of the Handler interface, including a TextHandler
	// that writes records to an io.Writer in a human-readable format and a
	// JSONHandler that writes records in JSON format.
	l5 := slog.New(NewEmptyHandler(os.Stdout, nil))

	l5.Info("hello, world", "user", os.Getenv("USER"))
	// Output:
	// {2023-08-23 15:31:33.064435942 +0300 EEST m=+0.000217892 hello, world INFO 4914901 [{user {[] 5 0xc000012075}} { {[] 0 <nil>}} { {[] 0 <nil>}} { {[] 0 <nil>}} { {[] 0 <nil>}}] 1 []}

}

type Handler interface {
	Enabled(context.Context, slog.Level) bool
	Handle(context.Context, slog.Record) error
	WithAttrs(attrs []slog.Attr) slog.Handler
	WithGroup(name string) slog.Handler
}

type EmptyHandler struct {
	Handler
}

func NewEmptyHandler(w io.Writer, opts *slog.HandlerOptions) Handler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &EmptyHandler{}
}

func (h *EmptyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (h *EmptyHandler) Handle(ctx context.Context, r slog.Record) error {
	_, err := fmt.Printf("%v\n", r)
	return err
}

func (h *EmptyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *EmptyHandler) WithGroup(name string) slog.Handler {
	return h
}
