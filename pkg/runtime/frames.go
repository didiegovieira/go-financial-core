package runtime

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// FuncName returns the name of the function provided.
//
// fn must reference a function, any other value will panic.
func FuncName(fn interface{}) string {
	v := reflect.ValueOf(fn)
	if v.Kind() != reflect.Func {
		panic(fmt.Errorf("can't obtain function name of a non function object: %T", fn))
	}
	return normalizeFunctionName(runtime.FuncForPC(v.Pointer()).Name())
}

// CallerName returns the name of the function that called this one
// after skipping some frames given by skip.
func CallerName(skip uint) string {
	return FrameFunctionName(Frame(skip + 1))
}

// FrameFunctionName returns the name of the function contained in the given
// stack frame
func FrameFunctionName(frame runtime.Frame) string {
	return normalizeFunctionName(frame.Function)
}

func normalizeFunctionName(name string) string {
	i := strings.LastIndex(name, "/")
	if i == -1 {
		return name
	}
	return name[i+1:]
}

// Frame returns the stack frame after skipping some frames given by the skip
// argument.
//
// A stack frame is an entry in the call stack, this is useful to see what
// functions, files and lines were used in the chain of calls.
func Frame(skip uint) runtime.Frame {
	var frame runtime.Frame
	for frame = range frames(2, skip+1) {
		// do nothing, just get the last frame
	}
	return frame
}

// Next returns the next frame in the call stack
func Next(frame *runtime.Frame) *runtime.Frame {
	var (
		next runtime.Frame
		save bool = false
	)

	// the stack has a maximum size in bytes, not in entries, any number used
	// here could be smaller or bigger than the real stack size, so a sane value
	// of 32 is used here to not kill runtime performance nor to miss the next
	// frame
	for f := range frames(2, 32) {
		if save {
			next = f
			save = false
		}
		if f.PC == frame.PC {
			save = true
		}
		// do not return after finding next, just continue consuming frames
		// until there is no more frames to consume
	}

	return &next
}

func frames(skip, stackSize uint) chan runtime.Frame {
	ch := make(chan runtime.Frame)

	pc := make([]uintptr, stackSize+1)
	n := runtime.Callers(int(skip), pc)
	if n == 0 {
		close(ch)
		return ch
	}

	frames := runtime.CallersFrames(pc[:n])

	go func() {
		for {
			frame, more := frames.Next()
			ch <- frame
			if !more {
				close(ch)
				return
			}
		}
	}()

	return ch
}
