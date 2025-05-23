<!-- gomarkdoc:embed:start -->

<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# sblog

```go
import "github.com/barbell-math/smoothbrain-logging"
```

A very simple library that implements a logger with verbosity levels and an optional rotating file log writer.

<details><summary>Example (Log Debug)</summary>
<p>



```go
s, _ := New(Opts{})
s.Debug("This is a debug message")
s.Debug(
	"This is a debug message",
	"Keys with nil values will be treated as messages", nil,
	"and will be printed on a separate line", nil,
)
s.Debug("This is a debug message logging a value", "value", 10)
s.Debug(
	"This is a debug message logging a multiple values",
	"value1", 10,
	"value2", 11,
)
s.Debug("This is a debug message", "Separate line", nil, "value", 10)
```

</p>
</details>

<details><summary>Example (Log Error)</summary>
<p>



```go
s, _ := New(Opts{})
s.Error("This is an error message")
s.Error(
	"This is an error message",
	"Keys with nil values will be treated as messages", nil,
	"and will be printed on a separate line", nil,
)
s.Error("This is an error message logging a value", "value", 10)
s.Error(
	"This is an error message logging a multiple values",
	"value1", 10,
	"value2", 11,
)
s.Error("This is an error message", "Separate line", nil, "value", 10)
```

</p>
</details>

<details><summary>Example (Log Info)</summary>
<p>



```go
s, _ := New(Opts{})
s.Info("This is an info message")
s.Info(
	"This is an info message",
	"Keys with nil values will be treated as messages", nil,
	"and will be printed on a separate line", nil,
)
s.Info("This is an info message logging a value", "value", 10)
s.Info(
	"This is an info message logging a multiple values",
	"value1", 10,
	"value2", 11,
)
s.Info("This is an info message", "Separate line", nil, "value", 10)
```

</p>
</details>

<details><summary>Example (Log Multi Error)</summary>
<p>



```go
err := sberr.AppendError(
	sberr.Wrap(
		ExpectedDirErr, "This is the first error",
	),
	sberr.Wrap(
		ExpectedDirErr, "This is the second error",
	),
	sberr.WrapValueList(
		ExpectedDirErr,
		"This is the third error",
		sberr.WrapListVal{
			ItemName: "Item 1",
			Item:     1,
		},
		sberr.WrapListVal{
			ItemName: "Item 2",
			Item:     1.0,
		},
		sberr.WrapListVal{
			ItemName: "Item 3",
			Item:     "asdf",
		},
	),
)

s, _ := New(Opts{})
s.Error(err.Error())
```

</p>
</details>

<details><summary>Example (Log Verbose)</summary>
<p>



```go
s, _ := New(Opts{CurVerbosityLevel: 2})
s.Log(
	context.TODO(),
	VLevel(0),
	"This is a level 0 verbose message, a.k.a. a debug message",
)
s.Log(context.TODO(), VLevel(1), "This is a level 1 verbose message")
s.Log(context.TODO(), VLevel(2), "This is a level 2 verbose message")
// This next line should not print
s.Log(context.TODO(), VLevel(3), "This is a level 3 verbose message")

s, _ = New(Opts{CurVerbosityLevel: 0})
s.Log(
	context.TODO(),
	VLevel(0),
	"This is a level 0 verbose message, a.k.a. a debug message",
)
// These next three lines should not print
s.Log(context.TODO(), VLevel(1), "This is a level 1 verbose message")
s.Log(context.TODO(), VLevel(2), "This is a level 2 verbose message")
s.Log(context.TODO(), VLevel(3), "This is a level 3 verbose message")
```

</p>
</details>

<details><summary>Example (Log Warn)</summary>
<p>



```go
s, _ := New(Opts{})
s.Warn("This is a warn message")
s.Warn(
	"This is a warn message",
	"Keys with nil values will be treated as messages", nil,
	"and will be printed on a separate line", nil,
)
s.Warn("This is a warn message logging a value", "value", 10)
s.Warn(
	"This is a warn message logging a multiple values",
	"value1", 10,
	"value2", 11,
)
s.Warn("This is a warn message", "Separate line", nil, "value", 10)
```

</p>
</details>

<details><summary>Example (Persistant Logs)</summary>
<p>



```go
s, _ := New(Opts{
	CurVerbosityLevel: 1,
	RotateWriterOpts: RotateWriterOpts{
		LogDir:  "./testData/",
		LogName: "persistantExample",
	},
})
s.Info("This should be saved to ./testData/persistantExample.0.log")
s.Warn("This should be saved to ./testData/persistantExample.0.log")
s.Error("This should be saved to ./testData/persistantExample.0.log")
s.Debug("This should be saved to ./testData/persistantExample.0.log")
s.Log(
	context.TODO(), VLevel(1),
	"This should be saved to ./testData/persistantExample.0.log",
)
```

</p>
</details>

## Index

- [Variables](<#variables>)
- [func New\(opts Opts\) \(\*slog.Logger, error\)](<#New>)
- [func VLevel\(val uint\) slog.Level](<#VLevel>)
- [type Opts](<#Opts>)
- [type RotateWriter](<#RotateWriter>)
  - [func NewRotateWriter\(opts RotateWriterOpts\) \(\*RotateWriter, error\)](<#NewRotateWriter>)
  - [func \(w \*RotateWriter\) Close\(\) error](<#RotateWriter.Close>)
  - [func \(w \*RotateWriter\) Rotate\(\) \(err error\)](<#RotateWriter.Rotate>)
  - [func \(w \*RotateWriter\) Write\(output \[\]byte\) \(int, error\)](<#RotateWriter.Write>)
- [type RotateWriterOpts](<#RotateWriterOpts>)


## Variables

<a name="ExpectedDirErr"></a>

```go
var (
    ExpectedDirErr = errors.New("A dir was expected")
)
```

<a name="New"></a>
## func [New](<https://github.com/barbell-math/smoothbrain-logging/blob/main/log.go#L49>)

```go
func New(opts Opts) (*slog.Logger, error)
```

Creates a new logger. The logs will always be printed to stdout/stderr. If \`opts.RotateWriterOpts.LogDir\` is an empty string no logs will be written to disk. If \`opts.RotateWriterOpts.LogDir\` is a valid dir then all logs from stdout and stderr will be mirrored to files in the supplied dir using a [RotateWriter](<#RotateWriter>).

<a name="VLevel"></a>
## func [VLevel](<https://github.com/barbell-math/smoothbrain-logging/blob/main/log.go#L40>)

```go
func VLevel(val uint) slog.Level
```

Translates the requested positive verbosity level from a range of \[0,inf\) to \(\-inf, \-4\]. This must be done for the level to be understood by the \[slog\] package.

<a name="Opts"></a>
## type [Opts](<https://github.com/barbell-math/smoothbrain-logging/blob/main/log.go#L26-L34>)



```go
type Opts struct {
    // Sets the allowed maximum verbosity level. There is no bound to the
    // requested verbosity, but any verbosity levels greater than the level
    // set here will not be printed.
    CurVerbosityLevel uint
    // The time format to use when printing log messages.
    TimeFmt string
    RotateWriterOpts
}
```

<a name="RotateWriter"></a>
## type [RotateWriter](<https://github.com/barbell-math/smoothbrain-logging/blob/main/rotatingWriter.go#L27-L34>)

Implements a rotating file writer. The logger will create log files in the following format:

```
<logDir>/<logName>.0.log
<logDir>/<logName>.1.log
...
<logDir>/<logName>.<maxNumLogs>.log
```

Where logDir, logName, and maxNumLogs are parameters that are specified through the [RotateWriterOpts](<#RotateWriterOpts>) struct when [NewRotateWriter](<#NewRotateWriter>) is called. Once the maximum allowed number of logs have been reached the writer will start to cycle back through old logs, clearing each one as needed before writing new log data.

```go
type RotateWriter struct {
    // contains filtered or unexported fields
}
```

<a name="NewRotateWriter"></a>
### func [NewRotateWriter](<https://github.com/barbell-math/smoothbrain-logging/blob/main/rotatingWriter.go#L68>)

```go
func NewRotateWriter(opts RotateWriterOpts) (*RotateWriter, error)
```

Creates a new rotating log writer. \`opts.LogDir\` must exist. Several defaults will be set from opts:

- If opts.LogName=="" then it will be set to a default of "sblog"
- If opts.MaxNumLogs==0 then it will be set to a default of 1
- If opts.MaxLogSizeBytes==0 then it will be set to a default of 1 MB

<a name="RotateWriter.Close"></a>
### func \(\*RotateWriter\) [Close](<https://github.com/barbell-math/smoothbrain-logging/blob/main/rotatingWriter.go#L133>)

```go
func (w *RotateWriter) Close() error
```

Closes the currently open log in the rotation. All other files in the log rotation are already closed.

<a name="RotateWriter.Rotate"></a>
### func \(\*RotateWriter\) [Rotate](<https://github.com/barbell-math/smoothbrain-logging/blob/main/rotatingWriter.go#L110>)

```go
func (w *RotateWriter) Rotate() (err error)
```

Rotates the log, advancing to the next log in the rotation. If the end of the log rotation is reached the rotation will start over at the beginning, clearing out logs as necessary to continue writing data.

<a name="RotateWriter.Write"></a>
### func \(\*RotateWriter\) [Write](<https://github.com/barbell-math/smoothbrain-logging/blob/main/rotatingWriter.go#L90>)

```go
func (w *RotateWriter) Write(output []byte) (int, error)
```

Writes data to the log rotation, advancing to the next log in the rotation as necessary.

<a name="RotateWriterOpts"></a>
## type [RotateWriterOpts](<https://github.com/barbell-math/smoothbrain-logging/blob/main/rotatingWriter.go#L36-L56>)



```go
type RotateWriterOpts struct {
    // The dir that the log files should be placed in. The dir must exist.
    // It will not be created if it does not exist.
    LogDir string
    // The name to give all log files. All log files will then match the
    // following format:
    //
    //   <logName>.<N>.log
    //
    // Where N is the number of the log in the log rotation.
    LogName string
    // The maximum number of log files to allow. Once this limit has been
    // reached the log writer will start to cycle back through old logs,
    // clearing each one as needed.
    MaxNumLogs uint
    // The number of bytes that can be reached before a moving onto the next
    // log file in the log rotation. Log messages will not be split between
    // files, meaning the number of bytes in the file may be slightly more
    // than `MaxLogSizeBytes`.
    MaxLogSizeBytes uint64
}
```

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)


<!-- gomarkdoc:embed:end -->

## Helpful Developer Cmds

To build the build system:

```
go build -o ./bs/bs ./bs
```

The build system can then be used as usual:

```
./bs/bs --help
./bs/bs buildbs # Builds the build system!
```
