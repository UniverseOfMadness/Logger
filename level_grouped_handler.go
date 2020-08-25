package logger

import "fmt"

type LevelGroup struct {
	Levels  []Level
	Handler Handler
}

type LevelGroupedHandler struct {
	handlers        map[Level][]Handler
	fallbackHandler Handler
}

func NewLevelGroupedHandler(fallbackHandler Handler, groups ...*LevelGroup) *LevelGroupedHandler {
	h := &LevelGroupedHandler{fallbackHandler: fallbackHandler, handlers: make(map[Level][]Handler)}
	h.createHandlersList(groups)

	return h
}

func (h *LevelGroupedHandler) Handle(log Log) error {
	handlers, ok := h.handlers[log.Level]

	if !ok || len(handlers) == 0 {
		err := h.fallbackHandler.Handle(log)

		if err != nil {
			return fmt.Errorf("LevelGroupedHandler - fallback handler returned an error: %w", err)
		}

		return nil
	}

	for _, handler := range handlers {
		err := handler.Handle(log)

		if err != nil {
			return fmt.Errorf("LevelGroupedHandler - one of handlers returned an error: %w", err)
		}
	}

	return nil
}

func (h *LevelGroupedHandler) createHandlersList(groups []*LevelGroup) {
	for _, group := range groups {
		for _, level := range group.Levels {
			h.handlers[level] = append(h.handlers[level], group.Handler)
		}
	}
}
