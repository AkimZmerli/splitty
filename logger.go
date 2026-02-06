package splitty

import "github.com/charmbracelet/log"

// logger wraps charmbracelet/log for optional debug logging.
// All methods are nil-safe — a nil logger silently discards messages.
type logger struct {
	l *log.Logger
}

func newLogger(l *log.Logger) *logger {
	if l == nil {
		return nil
	}
	return &logger{l: l}
}

func (lg *logger) debug(msg string, keyvals ...interface{}) {
	if lg != nil && lg.l != nil {
		lg.l.Debug(msg, keyvals...)
	}
}

func (lg *logger) info(msg string, keyvals ...interface{}) {
	if lg != nil && lg.l != nil {
		lg.l.Info(msg, keyvals...)
	}
}

func (lg *logger) error(msg string, keyvals ...interface{}) {
	if lg != nil && lg.l != nil {
		lg.l.Error(msg, keyvals...)
	}
}
