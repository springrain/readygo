package webext

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

func WebRecovery() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				stack := stack(3)
				httpRequest := c.Request.Body()

				if brokenPipe {
					//logger.Panic(err.(error), logger.String("httpRequest", string(httpRequest)), logger.String("reset", reset))
					hlog.CtxErrorf(ctx, "{err:%v,httpRequest:%v,reset:%v}", err, httpRequest, reset)
					//logger.Printf("%s\n%s%s", err, string(httpRequest), reset)
				} else {
					//logger.Printf("[Recovery] %s panic recovered:\n%s\n%s%s",timeFormat(time.Now()), err, stack, reset)
					//logger.Panic(err.(error), logger.ByteString("stack", stack), logger.String("reset", reset))
					hlog.CtxErrorf(ctx, "{err:%v,stack:%v,reset:%v}", err, stack, reset)
				}

				// If the connection is dead, we can't write a status to it.
				if brokenPipe {
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
				} else {
					c.AbortWithStatus(http.StatusInternalServerError)
				}
			}
		}()
		c.Next(ctx)
	}
}

// stack returns a nicely formatted stack frame, skipping skip frames.
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		//fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		//logger.Panic(errors.New("[GIN]WebRecovery-->stack"), logger.String("file", file), logger.Int("line", line), logger.Uintptr("pc", pc))
		hlog.Errorf("{err:%v,file:%v,line:%v,pc:%v}", errors.New("[GIN]WebRecovery-->stack"), file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		//fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
		//logger.Panic(errors.New("[GIN]WebRecovery-->stack"), logger.ByteString("pc", function(pc)), logger.ByteString("lines", source(lines, line)))
		hlog.Errorf("{err:%v,lines:%v,pc:%v}", errors.New("[GIN]WebRecovery-->stack"), source(lines, line), pc)
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

func timeFormat(t time.Time) string {
	var timeString = t.Format("2006-01-02 15:04:05")
	return timeString
}