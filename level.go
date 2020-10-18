package logger

import "fmt"

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

func (l Level) EqualOrGreaterThan(level Level) bool {
	return l >= level
}

func (l Level) Name() LevelName {
	switch l {
	case LevelDebug:
		return LevelNameDebug
	case LevelInfo:
		return LevelNameInfo
	case LevelWarning:
		return LevelNameWarning
	case LevelError:
		return LevelNameError
	case LevelCritical:
		return LevelNameCritical
	}

	panic(fmt.Errorf("unknown level %d", l))
}

func (ln LevelName) Level() Level {
	switch ln {
	case LevelNameDebug:
		return LevelDebug
	case LevelNameInfo:
		return LevelInfo
	case LevelNameWarning:
		return LevelWarning
	case LevelNameError:
		return LevelError
	case LevelNameCritical:
		return LevelCritical
	}

	panic(fmt.Errorf("unknown level name %s", ln))
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
