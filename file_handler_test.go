package logger

import (
	"github.com/stretchr/testify/assert"
	"os"
	"sync"
	"testing"
	"time"
)

func TestFileHandler_Handle(t *testing.T) {
	t.Parallel()

	t.Run("without formatter", func(t *testing.T) {
		f := &mockFile{}
		f.On("WriteString", "test message\n").Return(5, nil)

		fs := &mockFilesystem{}
		fs.On("OpenFile", "/path/to/log.file", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.FileMode(0644)).Return(f, nil)

		fh := NewFileHandler("/path/to/log.file")
		fh.WithFilesystem(fs)

		hErr := fh.Handle(Log{
			Level:     LevelInfo,
			Message:   "test message",
			Data:      make(Data),
			CreatedAt: time.Now(),
		})

		if assert.NoError(t, hErr) {
			f.AssertExpectations(t)
			fs.AssertExpectations(t)
		}
	})

	t.Run("with formatter", func(t *testing.T) {
		lg := Log{
			Level:     LevelInfo,
			Message:   "test message",
			Data:      make(Data),
			CreatedAt: time.Now(),
		}

		fr := &mockFormatter{}
		fr.On("Format", lg).Return(FormattedLog{
			Log:              lg,
			FormattedMessage: "formatted message",
		})

		f := &mockFile{}
		f.On("WriteString", "formatted message\n").Return(5, nil)

		fs := &mockFilesystem{}
		fs.On("OpenFile", "/path/to/log.file", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.FileMode(0644)).Return(f, nil)

		fh := NewFileHandler("/path/to/log.file")
		fh.UseFormatter(fr)
		fh.WithFilesystem(fs)

		hErr := fh.Handle(lg)

		if assert.NoError(t, hErr) {
			f.AssertExpectations(t)
			fs.AssertExpectations(t)
		}
	})

	t.Run("async", func(t *testing.T) {
		f := &mockFile{}
		f.On("WriteString", "test message\n").Return(5, nil)

		fs := &mockFilesystem{}
		fs.On("OpenFile", "/path/to/log.file", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.FileMode(0644)).Return(f, nil)

		fh := NewFileHandler("/path/to/log.file")
		fh.WithFilesystem(fs)

		wg := &sync.WaitGroup{}
		wg.Add(100)

		for i := 1; i <= 100; i++ {
			go func() {
				hErr := fh.Handle(Log{
					Level:     LevelInfo,
					Message:   "test message",
					Data:      make(Data),
					CreatedAt: time.Now(),
				})

				assert.NoError(t, hErr)
				wg.Done()
			}()
		}

		wg.Wait()

		fs.AssertExpectations(t)
		fs.AssertNumberOfCalls(t, "OpenFile", 1)

		f.AssertExpectations(t)
		f.AssertNumberOfCalls(t, "WriteString", 100)
	})
}

func TestFileHandler_HandleBatch(t *testing.T) {
	t.Parallel()

	t.Run("without formatter", func(t *testing.T) {
		f := &mockFile{}
		f.On("Write", []byte("test message\ntest message 2\n")).Return(30, nil)

		fs := &mockFilesystem{}
		fs.On("OpenFile", "/path/to/log.file", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.FileMode(0644)).Return(f, nil)

		fh := NewFileHandler("/path/to/log.file")
		fh.WithFilesystem(fs)

		hErr := fh.HandleBatch([]Log{
			{
				Level:     LevelInfo,
				Message:   "test message",
				Data:      make(Data),
				CreatedAt: time.Now(),
			},
			{
				Level:     LevelDebug,
				Message:   "test message 2",
				Data:      make(Data),
				CreatedAt: time.Now(),
			},
		})

		if assert.NoError(t, hErr) {
			f.AssertExpectations(t)
			f.AssertNumberOfCalls(t, "Write", 1)

			fs.AssertExpectations(t)
			fs.AssertNumberOfCalls(t, "OpenFile", 1)
		}
	})

	t.Run("with formatter", func(t *testing.T) {
		lg := []Log{
			{
				Level:     LevelInfo,
				Message:   "test message",
				Data:      make(Data),
				CreatedAt: time.Now(),
			},
			{
				Level:     LevelDebug,
				Message:   "test message 2",
				Data:      make(Data),
				CreatedAt: time.Now(),
			},
		}

		fr := &mockFormatter{}
		fr.On("Format", lg[0]).Return(FormattedLog{
			Log:              lg[0],
			FormattedMessage: "formatted message",
		}).Once()
		fr.On("Format", lg[1]).Return(FormattedLog{
			Log:              lg[1],
			FormattedMessage: "formatted message 2",
		}).Once()

		f := &mockFile{}
		f.On("Write", []byte("formatted message\nformatted message 2\n")).Return(42, nil)

		fs := &mockFilesystem{}
		fs.On("OpenFile", "/path/to/log.file", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.FileMode(0644)).Return(f, nil)

		fh := NewFileHandler("/path/to/log.file")
		fh.UseFormatter(fr)
		fh.WithFilesystem(fs)

		hErr := fh.HandleBatch(lg)

		if assert.NoError(t, hErr) {
			f.AssertExpectations(t)
			f.AssertNumberOfCalls(t, "Write", 1)

			fs.AssertExpectations(t)
			fs.AssertNumberOfCalls(t, "OpenFile", 1)

			fr.AssertExpectations(t)
			fr.AssertNumberOfCalls(t, "Format", 2)
		}
	})

	t.Run("async", func(t *testing.T) {
		f := &mockFile{}
		f.On("Write", []byte("test message\ntest message 2\n")).Return(30, nil)

		fs := &mockFilesystem{}
		fs.On("OpenFile", "/path/to/log.file", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.FileMode(0644)).Return(f, nil)

		fh := NewFileHandler("/path/to/log.file")
		fh.WithFilesystem(fs)

		wg := &sync.WaitGroup{}
		wg.Add(100)

		for i := 1; i <= 100; i++ {
			go func() {
				hErr := fh.HandleBatch([]Log{
					{
						Level:     LevelInfo,
						Message:   "test message",
						Data:      make(Data),
						CreatedAt: time.Now(),
					},
					{
						Level:     LevelDebug,
						Message:   "test message 2",
						Data:      make(Data),
						CreatedAt: time.Now(),
					},
				})

				assert.NoError(t, hErr)
				wg.Done()
			}()
		}

		wg.Wait()

		fs.AssertExpectations(t)
		fs.AssertNumberOfCalls(t, "OpenFile", 1)

		f.AssertExpectations(t)
		f.AssertNumberOfCalls(t, "Write", 100)
	})
}
