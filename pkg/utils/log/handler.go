package log

import "github.com/charmbracelet/log"

type Handler struct {
	*log.Logger
	level Level
}

func NewHandler() *Handler {
	return &Handler{
		Logger: log.Default(),
	}
}

func (h *Handler) Debug(msg string, args ...any) {
	h.Logger.Debugf(msg, args...)
}

func (h *Handler) Info(msg string, args ...any) {
	h.Logger.Infof(msg, args...)
}

func (h *Handler) Warn(msg string, args ...any) {
	h.Logger.Warnf(msg, args...)
}

func (h *Handler) Error(msg string, args ...any) {
	h.Logger.Errorf(msg, args...)
}

func (h *Handler) SetLevel(level Level) {
	h.level = level

	h.Logger.SetLevel(log.Level(level))
}

func (h *Handler) Level() Level {
	return h.level
}
