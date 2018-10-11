package posix_mq

/*
#include <mqueue.h>

// Expose non-variadic function requires 2 arguments.
mqd_t mq_open2(const char *name, int oflag) {
	return mq_open(name, oflag);
}
*/
import "C"
import "unsafe"

func mq_open(name string, oflag int) (int, error) {
	h, err := C.mq_open2(C.CString(name), C.int(oflag))
	if err != nil {
		return 0, err
	}

	return int(h), nil
}

func mq_send(h int, data []byte, priority uint) (int, error) {
	byteStr := *(*string)(unsafe.Pointer(&data))
	rv, err := C.mq_send(C.int(h), C.CString(byteStr), C.size_t(len(data)), C.uint(priority))
	return int(rv), err
}

func mq_unlink(name string) (int, error) {
	rv, err := C.mq_unlink(C.CString(name))
	return int(rv), err
}
