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

var pinMap = make(map[string]Sample)

type SampleCallback func(sample Sample, number int) error
type sampleCallback func(number C.int) C.int

type sample struct {
	sample   unsafe.Pointer
	callback sampleCallback
}

type Sample interface {
	Destroy() error
	InvokeCallback() error
	RegisterCallback(callback SampleCallback) error
}

func NewSample(number int) Sample {
	result := sample{
		sample: unsafe.Pointer(C.create_sample(C.int(number))),
	}
	return &result
}

func (s *sample) Destroy() error {
	return NewErrno(int(C.destroy_sample(s.sample)))
}

//export e_callback
func e_callback(sample unsafe.Pointer, number C.int, data unsafe.Pointer) C.int {
	f := *(*sampleCallback)(data)
	return f(number)
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
	return NewErrno(int(C.register_callback(s.sample, (C.SampleCallback)(C.e_callback), unsafe.Pointer(&s.callback))))
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
