package fdw

// https://www.postgresql.org/docs/current/static/error-message-reporting.html
//
// Errors

/*
#include "go_fdw.h"

static inline ErrorData *
goCreateError(
	_GoString_ filename, int line,
	_GoString_ function,
	_GoString_ message)
{
	char *buffer = (char *) palloc0(_GoStringLen(filename)+_GoStringLen(function)+_GoStringLen(message)+3);
	ErrorData *edata;

	// initialize an error on the stack, copy it off, and reset the stack
	errstart(ERROR, NULL, line, NULL, NULL);
	edata = CopyErrorData();
	FlushErrorState();

	edata->filename = memcpy(buffer, _GoStringPtr(filename), _GoStringLen(filename));
	edata->funcname = memcpy(buffer+_GoStringLen(filename)+1, _GoStringPtr(function), _GoStringLen(function));
	edata->message = memcpy(buffer+_GoStringLen(filename)+_GoStringLen(function)+2, _GoStringPtr(message), _GoStringLen(message));

	return edata;
}
*/
import "C"
import (
	"path"
	"runtime"
)

type errorData struct{ edata *C.ErrorData }

func (e errorData) Error() string { return C.GoString(e.edata.message) }

func goErrorData(edata *C.ErrorData) error {
	if edata == nil {
		return nil
	}
	return errorData{edata}
}

func pgErrorData(err error) *C.ErrorData {
	if ed, ok := err.(errorData); ok {
		return ed.edata
	}

	message := err.Error()
	function := "unknown"

	pc, filename, line, _ := runtime.Caller(1)
	if fn := runtime.FuncForPC(pc); fn != nil {
		filename = path.Base(filename)
		function = fn.Name()
	}

	return C.goCreateError(filename, C.int(line), function, message)
}
