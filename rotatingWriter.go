package sblog

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sync"

	sberr "github.com/barbell-math/smoothbrain-errs"
)

type (
	// Implements a rotating file writer. The logger will create log files in
	// the following format:
	//
	//   <logDir>/<logName>.0.log
	//   <logDir>/<logName>.1.log
	//   ...
	//   <logDir>/<logName>.<maxNumLogs>.log
	//
	// Where logDir, logName, and maxNumLogs are parameters that are specified
	// through the [RotateWriterOpts] struct when [NewRotateWriter] is called.
	// Once the maximum allowed number of logs have been reached the writer will
	// start to cycle back through old logs, clearing each one as needed before
	// writing new log data.
	RotateWriter struct {
		lock     sync.Mutex
		fp       *os.File
		curFile  uint
		numBytes uint64

		opts RotateWriterOpts
	}

	RotateWriterOpts struct {
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
)

var (
	ExpectedDirErr = errors.New("A dir was expected")
)

// Creates a new rotating log writer. `opts.LogDir` must exist. Several defaults
// will be set from opts:
//   - If opts.LogName=="" then it will be set to a default of "sblog"
//   - If opts.MaxNumLogs==0 then it will be set to a default of 1
//   - If opts.MaxLogSizeBytes==0 then it will be set to a default of 1 MB
func NewRotateWriter(opts RotateWriterOpts) (*RotateWriter, error) {
	fs, err := os.Stat(opts.LogDir)
	if err != nil {
		return nil, err
	} else if !fs.IsDir() {
		return nil, sberr.Wrap(ExpectedDirErr, opts.LogDir)
	}
	if opts.LogName == "" {
		opts.LogName = "sblog"
	}
	if opts.MaxNumLogs == 0 {
		opts.MaxNumLogs = 1
	}
	if opts.MaxLogSizeBytes == 0 {
		opts.MaxLogSizeBytes = 1000000 // 1 MB
	}
	w := &RotateWriter{opts: opts}
	return w, w.Rotate()
}

// Writes data to the log rotation, advancing to the next log in the rotation as
// necessary.
func (w *RotateWriter) Write(output []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	n, err := w.fp.Write(output)
	if err != nil {
		return n, err
	}

	w.numBytes += uint64(n)
	if w.numBytes > w.opts.MaxLogSizeBytes {
		err = w.Rotate()
		w.numBytes = 0
	}

	return n, err
}

// Rotates the log, advancing to the next log in the rotation. If the end of the
// log rotation is reached the rotation will start over at the beginning,
// clearing out logs as necessary to continue writing data.
func (w *RotateWriter) Rotate() (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.fp != nil {
		err = w.fp.Close()
		w.fp = nil
		if err != nil {
			return
		}
	}

	w.curFile = (w.curFile + 1) % w.opts.MaxNumLogs
	fName := path.Join(
		w.opts.LogDir, fmt.Sprintf("%s.%d.log", w.opts.LogName, w.curFile),
	)

	w.fp, err = os.Create(fName)
	return
}

// Closes the currently open log in the rotation. All other files in the log
// rotation are already closed.
func (w *RotateWriter) Close() error {
	err := w.fp.Close()
	w.fp = nil
	return err
}
