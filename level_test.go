package logger

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestLevel_EqualOrGreaterThan(t *testing.T) {
	t.Parallel()

	assert.True(t, LevelDebug.EqualOrGreaterThan(LevelDebug))
	assert.True(t, LevelError.EqualOrGreaterThan(LevelInfo))
	assert.False(t, LevelInfo.EqualOrGreaterThan(LevelError))
}

func TestLevel_Int(t *testing.T) {
	t.Parallel()

	assert.IsType(t, 0, Level(9999).Int())
}

func TestLevelName_String(t *testing.T) {
	t.Parallel()

	assert.IsType(t, "test", LevelName("testing").String())
}

func TestLevel_Name(t *testing.T) {
	t.Parallel()

	assert.Equal(t, LevelDebug.Name(), LevelNameDebug)
	assert.Equal(t, LevelInfo.Name(), LevelNameInfo)
	assert.Equal(t, LevelWarning.Name(), LevelNameWarning)
	assert.Equal(t, LevelError.Name(), LevelNameError)
	assert.Equal(t, LevelCritical.Name(), LevelNameCritical)

	assert.PanicsWithError(t, "unknown level 9999", func() {
		Level(9999).Name()
	})
}

func TestLevelName_Level(t *testing.T) {
	t.Parallel()

	assert.Equal(t, LevelNameDebug.Level(), LevelDebug)
	assert.Equal(t, LevelNameInfo.Level(), LevelInfo)
	assert.Equal(t, LevelNameWarning.Level(), LevelWarning)
	assert.Equal(t, LevelNameError.Level(), LevelError)
	assert.Equal(t, LevelNameCritical.Level(), LevelCritical)

	assert.PanicsWithError(t, "unknown level name testing", func() {
		LevelName("testing").Level()
	})
}

func TestLevels_Sorting(t *testing.T) {
	t.Parallel()

	lvl := Levels{LevelInfo, LevelCritical, LevelDebug, LevelError, LevelWarning}
	sort.Sort(lvl)

	assert.Equal(t, Levels{LevelDebug, LevelInfo, LevelWarning, LevelError, LevelCritical}, lvl)
}
