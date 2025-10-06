package sl

import "log/slog"

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func Operation(msg string) slog.Attr {
	return slog.Attr{
		Key:   "op",
		Value: slog.StringValue(msg),
	}
}
