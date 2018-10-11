package posix_mq

/*
#include <fcntl.h>
#include <mqueue.h>

// Expose non-variadic function requires 4 arguments.
mqd_t mq_open4(const char *name, int oflag, int mode, struct mq_attr *attr) {
	return mq_open(name, oflag, mode, attr);
}
*/
import "C"
import "unsafe"

const (
	O_RDONLY = C.O_RDONLY
	O_WRONLY = C.O_WRONLY
	O_RDWR   = C.O_RDWR

	O_CLOEXEC  = C.O_CLOEXEC
	O_CREAT    = C.O_CREAT
	O_EXCL     = C.O_EXCL
	O_NONBLOCK = C.O_NONBLOCK
)

func mq_open(name string, oflag int, mode int) (int, error) {
	// TODO Support mq_attr
	h, err := C.mq_open4(C.CString(name), C.int(oflag), C.int(mode), nil)
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
