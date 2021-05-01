package logger

import (
	"errors"
)

type (
	Level     uint32
	LevelName string
	Levels    []Level
)

const (
	LevelDebug    = Level(0)
	LevelInfo     = Level(1000)
	LevelWarning  = Level(2000)
	LevelError    = Level(3000)
	LevelCritical = Level(9001)

	LevelNameDebug    = LevelName("debug")
	LevelNameInfo     = LevelName("info")
	LevelNameWarning  = LevelName("warning")
	LevelNameError    = LevelName("error")
	LevelNameCritical = LevelName("critical")
)

var (
	ErrLevelMappingNotFound     = errors.New("level cannot be mapped to level name")
	ErrLevelNameMappingNotFound = errors.New("level name cannot be mapped to level")

	levelsMapping = map[LevelName]Level{
		LevelNameDebug:    LevelDebug,
		LevelNameInfo:     LevelInfo,
		LevelNameWarning:  LevelWarning,
		LevelNameError:    LevelError,
		LevelNameCritical: LevelCritical,
	}
	levelsReverseMapping = map[Level]LevelName{
		LevelDebug:    LevelNameDebug,
		LevelInfo:     LevelNameInfo,
		LevelWarning:  LevelNameWarning,
		LevelError:    LevelNameError,
		LevelCritical: LevelNameCritical,
	}
)

func (l Level) EqualOrGreaterThan(level Level) bool {
	return l >= level
}

func (l Level) Name() (LevelName, error) {
	val, ok := levelsReverseMapping[l]

	if !ok {
		return "", ErrLevelNameMappingNotFound
	}

	return val, nil
}

func (ln LevelName) Level() (Level, error) {
	val, ok := levelsMapping[ln]

	if !ok {
		return 0, ErrLevelMappingNotFound
	}

	return val, nil
}

func (l Level) Int() int {
	return int(l)
}

func (ln LevelName) String() string {
	return string(ln)
}

func (l Levels) Len() int {
	return len(l)
}

func (l Levels) Less(i, j int) bool {
	return l[i] < l[j]
}

func (l Levels) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]

}
