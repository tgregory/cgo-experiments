package main

/*
#cgo LDFLAGS: libsample.a
#include "sample.h"

extern int e_callback(Sample* sample, int number, void *arbitrary_data);
*/
import "C"
import "unsafe"
import "fmt"

type Errno int

func (e *Errno) Error() string {
	return fmt.Sprintf("errno %d", int(*e))
}

func NewErrno(code int) *Errno {
	if 0 == code {
		return nil
	}
	e := Errno(code)
	return &e
}

// not goroutine safe
var pinMap = make(map[int]*sample)
var currentId = 0

type SampleCallback func(sample Sample, number int) error
type sampleCallback func(number C.int) C.int

type sample struct {
	id       int
	sample   unsafe.Pointer
	callback sampleCallback
}

type Sample interface {
	Destroy() error
	InvokeCallback() error
	RegisterCallback(callback SampleCallback) error
}

// this is not safe in case of multiple goroutines trying to do this at the same time
// also possibility of overflow exists
func NewSample(number int) Sample {
	result := sample{
		id:     currentId,
		sample: unsafe.Pointer(C.create_sample(C.int(number))),
	}
	currentId++
	pinMap[result.id] = &result
	return &result
}

func (s *sample) Destroy() error {
	delete(pinMap, s.id)
	return NewErrno(int(C.destroy_sample(s.sample)))
}

//export e_callback
func e_callback(sample unsafe.Pointer, number C.int, data unsafe.Pointer) C.int {
	id := int(uintptr(data))
	s := pinMap[id]
	return s.callback(number)
}

func (s *sample) InvokeCallback() error {
	return NewErrno(int(C.invoke_callback(s.sample)))
}

func (s *sample) RegisterCallback(callback SampleCallback) error {
	s.callback = func(number C.int) C.int {
		if e := callback(s, int(number)); nil == e {
			return 0
		}
		return 1
	}
	return NewErrno(int(C.register_callback(s.sample, (C.SampleCallback)(C.e_callback), unsafe.Pointer(uintptr(s.id)))))
}

func main() {
	samp := NewSample(10)
	defer samp.Destroy()
	samp.RegisterCallback(func(s Sample, num int) error {
		fmt.Printf("Hello from callback magic number received %d \n", num)
		return nil
	})
	samp.InvokeCallback()
	samp.InvokeCallback()
}
