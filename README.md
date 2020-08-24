# Golang Logger
Simple logger for Golang application that allows creating custom handling
methods as well as final format of the message.

## Requirements
 * Golang version >= 1.14

## How to use?
Example of basic usage.
```go
package main

import (
    "github.com/UniverseOfMadness/logger"
    "os"
    "time"
)

func main() {
    f := logger.NewBasicFormatter("MyApp", time.RFC3339)

    h := logger.NewStringWriterHandler(os.Stdout)
    h.UseFormatter(f)

    l := logger.New(h)
    l.SetLevel(logger.LevelInfo)

    l.Info("app started at {now}", "now", time.Now().Format(time.RFC3339))
    l.Infof("hello %s", "world")
}
```

### Logging functions.
There are two types of logging functions:
 * Standard type which allow setting data parameters in Log.
 * Formatted type which uses "fmt.Sprintf" to create final message in Log from provided message and values.

List of standard logging functions:
 * Logger.Debug
 * Logger.Info
 * Logger.Warning
 * Logger.Error
 * Logger.Critical

List of formatted logging functions:
 * Logger.Debugf
 * Logger.Infof
 * Logger.Warningf
 * Logger.Errorf
 * Logger.Criticalf

## Custom handlers
Package includes `Handler` interface that can be used to create custom handlers for
logger. `StringWriterHandler` can be used as example for implementation.

## Custom formatters
Package includes `Formatter` interface that can be used to create custom formatters for
logger. `BasicFormatter` can be used as example for implementation.

## Custom clock
Default clock used in `Logger` is only a wrapper for built-in Golang `time.*`.
If application that implements this package requires a special time adjustment then
interface `Clock` can be used to create custom implementation for the clock.

## Log Levels
 * **Debug** [0] - detailed information, mostly for development or debugging.
 * **Info** [1000] - basic info message for normal application flow (new account, finished process etc.).
 * **Warning** [2000] - warning message for "expected errors" (use of deprecated functions, running outdated version of application etc.).
 * **Error** [3000] - unexpected errors that should not occur but does not break flow of application.
 * **Critical** [9001] - unexpected errors that break flow of application.

### Critical Handler
Logs with critical level can trigger some additional events in application with `CriticalHandleFunc`
set in logger using `WithCriticalHandler` function. Logger does nothing in case of critical errors by default.

### Failure Handler
Logger will ignore all handlers errors by default. There is `FailureHandleFunc` that can be implemented and added to logger
using `WithFailureHandler` function. Whenever the handler returns error, this function will be triggered with `Log` and 
error as parameters.
