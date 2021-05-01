package logger

import (
	"fmt"
	"sort"
)

type LevelGroup struct {
	Levels  []Level
	Handler Handler
}

type LevelGroupedHandler struct {
	handlers        map[Level][]Handler
	fallbackHandler Handler
}

func NewLevelGroupedHandler(fallbackHandler Handler, groups ...LevelGroup) *LevelGroupedHandler {
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

func (h *LevelGroupedHandler) HandleBatch(logs []Log) error {
	var levels Levels
	levelLogs := make(map[Level][]Log)

	for _, l := range logs {
		levelLogs[l.Level] = append(levelLogs[l.Level], l)
		levels = append(levels, l.Level)
	}

	sort.Sort(levels)

	for _, lvl := range levels {
		if hd, hOk := h.handlers[lvl]; hOk {
			for _, hand := range hd {
				hErr := hand.HandleBatch(levelLogs[lvl])

				if hErr != nil {
					return fmt.Errorf("LevelGroupedHandler - one of handlers returned an error: %w", hErr)
				}
			}
		} else {
			hErr := h.fallbackHandler.HandleBatch(levelLogs[lvl])

			if hErr != nil {
				return fmt.Errorf("LevelGroupedHandler - fallback handler returned an error: %w", hErr)
			}
		}
	}

	return nil
}

func (h *LevelGroupedHandler) createHandlersList(groups []LevelGroup) {
	for _, group := range groups {
		for _, level := range group.Levels {
			h.handlers[level] = append(h.handlers[level], group.Handler)
		}
	}
}
