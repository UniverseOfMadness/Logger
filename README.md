[![Build Status](https://travis-ci.com/UniverseOfMadness/logger.svg?branch=master)](https://travis-ci.com/UniverseOfMadness/logger)

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

## Handlers
List of handlers provided with package:
 * [StringWriterHandler](https://github.com/UniverseOfMadness/logger/blob/master/string_writer_handler.go) - takes any struct that implements `io.StringWriter` interface
 and writes all incoming logs to it.
 * [LevelGroupedHandler](https://github.com/UniverseOfMadness/logger/blob/master/level_grouped_handler.go) - groups handlers by level so each handler is 
 able to handle logs for specific level. There is also parameter that accept fallback handler for non-defined levels.
 * [InMemoryHandler](https://github.com/UniverseOfMadness/logger/blob/master/in_memory_handler.go) - stores all logs in-memory (as slice). Each log can be popped from slice individually.
 Handler can also be cleared. Constructor for handler takes `bufferOverflow` as parameter which is max number of logs stored in the handler. Any log added above limit will cause an error.

### Custom handlers
Package includes `Handler` interface that can be used to create custom handlers for
logger. `StringWriterHandler` can be used as example for implementation.

## Formatters
List of formatters provided with package:
 * [BasicFormatter](https://github.com/UniverseOfMadness/logger/blob/master/basic_formatter.go) - standard log formatter which produce easy to read message
 (example: `SimpleWebServer | 2020-08-25T19:06:36+02:00 | INFO | server is listening on 17333 | port:17333`). Allows setting application name and format for log date time.

### Custom formatters
Package includes `Formatter` interface that can be used to create custom formatters for
logger. `BasicFormatter` can be used as example for implementation.

## Clock
Default clock used in `Logger` is only a wrapper for built-in Golang `time.*`.
If application that implements this package requires a special time adjustment then
interface `Clock` can be used to create custom implementation for the clock.

## Error Wrapper
By default, Logger requires checking if error actually occurred before sending it to logs.
`ErrorWrappedLogger` can handle errors directly with `nil` check. For example:
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
    err := DoSomethingThatWillFail()

    // handling error log without wrapper
    if err != nil {
        l.Error(err.Error())
    }
    
    wr := logger.NewErrorWrappedLogger(l)
    // handling error log with wrapper
    wr.OnError(err)
    
    // additionally there is a possibility to wrap error with message
    wr.OnErrorWrapped(err, "something went wrong in %s func: %w", "main")
}
```

### Functions
 * **OnError** - passes Go error message to Logger if error is not `nil`. 
 * **OnErrorWrapped** - passes Go error message wrapped using `fmt.Errorf` to Logger if error is not `nil`.
 * **OnCritical** - works the same way as `OnError` but passes message to `Critical` instead of `Error`.
 * **OnCriticalWrapped** - works the same way as `OnErrorWrapped` but passes message to `Critical` instead of `Error`.

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
