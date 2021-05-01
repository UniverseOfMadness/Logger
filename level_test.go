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

	ld, ldErr := LevelDebug.Name()
	assert.Equal(t, ld, LevelNameDebug)
	assert.NoError(t, ldErr)

	li, liErr := LevelInfo.Name()
	assert.Equal(t, li, LevelNameInfo)
	assert.NoError(t, liErr)

	lw, lwErr := LevelWarning.Name()
	assert.Equal(t, lw, LevelNameWarning)
	assert.NoError(t, lwErr)

	le, leErr := LevelError.Name()
	assert.Equal(t, le, LevelNameError)
	assert.NoError(t, leErr)

	lc, lcErr := LevelCritical.Name()
	assert.Equal(t, lc, LevelNameCritical)
	assert.NoError(t, lcErr)

	ul, ulErr := Level(9999).Name()
	assert.Empty(t, ul)
	assert.Error(t, ulErr)
	assert.EqualError(t, ulErr, ErrLevelNameMappingNotFound.Error())
}

func TestLevelName_Level(t *testing.T) {
	t.Parallel()

	ld, ldErr := LevelNameDebug.Level()
	assert.Equal(t, ld, LevelDebug)
	assert.NoError(t, ldErr)

	li, liErr := LevelNameInfo.Level()
	assert.Equal(t, li, LevelInfo)
	assert.NoError(t, liErr)

	lw, lwErr := LevelNameWarning.Level()
	assert.Equal(t, lw, LevelWarning)
	assert.NoError(t, lwErr)

	le, leErr := LevelNameError.Level()
	assert.Equal(t, le, LevelError)
	assert.NoError(t, leErr)

	lc, lcErr := LevelNameCritical.Level()
	assert.Equal(t, lc, LevelCritical)
	assert.NoError(t, lcErr)

	ul, ulErr := LevelName("testing").Level()
	assert.Equal(t, ul, Level(0))
	assert.Error(t, ulErr)
	assert.EqualError(t, ulErr, ErrLevelMappingNotFound.Error())
}

func TestLevels_Sorting(t *testing.T) {
	t.Parallel()

	lvl := Levels{LevelInfo, LevelCritical, LevelDebug, LevelError, LevelWarning}
	sort.Sort(lvl)

	assert.Equal(t, Levels{LevelDebug, LevelInfo, LevelWarning, LevelError, LevelCritical}, lvl)
}
