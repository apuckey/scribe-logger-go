package logger

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

const (
	DEBUG    = 1
	INFO     = 2
	WARNING  = 4
	WARN     = 4
	ERROR    = 8
	NOTICE   = 16 //notice is like info but for really important stuff ;)
	CRITICAL = 32
	QUIET    = ERROR | NOTICE | CRITICAL               //setting for errors only
	NORMAL   = INFO | WARN | ERROR | NOTICE | CRITICAL // default setting - all besides debug
	ALL      = 255
	NOTHING  = 0
)

var levelsAscending = []int{DEBUG, INFO, WARNING, ERROR, NOTICE, CRITICAL}

var LevlelsByName = map[string]int{
	"DEBUG":    DEBUG,
	"INFO":     INFO,
	"WARNING":  WARN,
	"WARN":     WARN,
	"ERROR":    ERROR,
	"NOTICE":   NOTICE,
	"CRITICAL": CRITICAL,
	"QUIET":    QUIET,
	"NORMAL":   NORMAL,
	"ALL":      ALL,
	"NOTHING":  NOTHING,
}

//default logging level is ALL
var level = ALL

// Set the logging level.
//
// Contrary to Python that specifies a minimal level, this logger is set with a bit mask
// of active levels.
//
// e.g. for INFO and ERROR use:
// 		SetLevel(logging.INFO | logging.ERROR)
//
// For everything but debug and info use:
// 		SetLevel(logging.ALL &^ (logging.INFO | logging.DEBUG))
//
func SetLevel(l int) {
	level = l
}

// Set a minimal level for loggin, setting all levels higher than this level as well.
//
// the severity order is DEBUG, INFO, WARNING, ERROR, CRITICAL
func SetMinimalLevel(l int) {

	newLevel := 0
	for _, level := range levelsAscending {
		if level >= l {
			newLevel |= level
		}
	}
	SetLevel(newLevel)

}

// Set minimal level by string, useful for config files and command line arguments. Case insensitive.
//
// Possible level names are DEBUG, INFO, WARNING, ERROR, NOTICE, CRITICAL
func SetMinimalLevelByName(l string) error {
	l = strings.ToUpper(strings.Trim(l, " "))
	level, found := LevlelsByName[l]
	if !found {
		Error("Could not set level - not found level %s")
		return fmt.Errorf("Invalid level %s", l)
	}

	SetMinimalLevel(level)
	return nil
}

// Set the output writer. for now it just wraps log.SetOutput()
func SetOutput(w io.Writer) {
	log.SetOutput(w)
}

//a pluggable logger interface
type LoggingHandler interface {
	SetFormatter(Formatter)
	Emit(ctx *MessageContext, message string) error
}

type standardHandler struct {
	formatter Formatter
}

func (l *standardHandler) SetFormatter(f Formatter) {
	l.formatter = f
}

// default handling interface - just
func (l *standardHandler) Emit(ctx *MessageContext, message string) error {
	fmt.Fprintln(os.Stderr, l.formatter.Format(ctx, message))
	return nil
}

var currentHandler LoggingHandler = &standardHandler{
	DefaultFormatter,
}

// Set the current handler of the library. We currently support one handler, but it might be nice to have more
func SetHandler(h LoggingHandler) {
	currentHandler = h
}

func CurrentHandler() LoggingHandler {
	return currentHandler
}

type MessageContext struct {
	Level     string
	File      string
	Line      int
	TimeStamp time.Time
}

//get the stack (line + file) context to return the caller to the log
func getContext(level string, skipDepth int) *MessageContext {

	_, file, line, _ := runtime.Caller(skipDepth)
	file = path.Base(file)

	return &MessageContext{
		Level:     level,
		File:      file,
		TimeStamp: time.Now(),
		Line:      line,
	}
}

//Output debug logging messages
func Debug(msg string) {
	if level&DEBUG != 0 {
		writeMessage("DEBUG", msg)
	}
}

//format the message
func writeMessage(level string, msg string) {
	writeMessageDepth(4, level, msg)
}

func writeMessageDepth(depth int, level string, msg string) {
	ctx := getContext(level, depth)

	err := currentHandler.Emit(ctx, msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing log message: %s\n", err)
		fmt.Fprintln(os.Stderr, DefaultFormatter.Format(ctx, msg))
	}

}

//output INFO level messages
func Info(msg string) {
	if level&INFO != 0 {
		writeMessage("INFO", msg)
	}
}

// Output WARNING level messages
func Warning(msg string) {
	if level&WARN != 0 {
		writeMessage("WARNING", msg)
	}
}

// Same as Warning() but return a formatted error object, regardless of logging level
func Warningf(msg string) error {
	err := fmt.Errorf(msg)
	if level&WARN != 0 {
		writeMessage("WARNING", err.Error())
	}
	return err
}

// Output ERROR level messages
func Error(msg string) {
	if level&ERROR != 0 {
		writeMessage("ERROR", msg)
	}
}

// Same as Error() but also returns a new formatted error object with the message regardless of logging level
func Errorf(msg string) error {
	err := fmt.Errorf(msg)
	if level&ERROR != 0 {
		writeMessage("ERROR", err.Error())
	}
	return err
}

// Output NOTICE level messages
func Notice(msg string) {
	if level&NOTICE != 0 {
		writeMessage("NOTICE", msg)
	}
}

// Output a CRITICAL level message while showing a stack trace
func Critical(msg string) {
	if level&CRITICAL != 0 {
		writeMessage("CRITICAL", msg)
		log.Println(string(debug.Stack()))
	}
}

// Same as critical but also returns an error object with the message regardless of logging level
func Criticalf(msg string) error {

	err := fmt.Errorf(msg)
	if level&CRITICAL != 0 {
		writeMessage("CRITICAL", err.Error())
		log.Println(string(debug.Stack()))
	}
	return err
}

// Raise a PANIC while writing the stack trace to the log
func Panic(msg string) {
	log.Println(string(debug.Stack()))
	log.Panicf(msg)
}

func init() {
	log.SetFlags(0)
}

// bridge bridges the logger and the default go log, with a given level
type bridge struct {
	level     int
	levelName string
}

func (lb bridge) Write(p []byte) (n int, err error) {
	if level&lb.level != 0 {
		writeMessageDepth(6, lb.levelName, string(bytes.TrimRight(p, "\r\n")))
	}
	return len(p), nil
}

// BridgeStdLog bridges all messages written using the standard library's log.Print* and makes them output
// through this logger, at a given level.
func BridgeStdLog(level int) {

	for k, l := range LevlelsByName {
		if l == level {
			b := bridge{
				level:     l,
				levelName: k,
			}

			log.SetOutput(b)
		}
	}

}
