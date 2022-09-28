package posix_mq

/*
#cgo LDFLAGS: -lrt

#include <stdlib.h>
#include <signal.h>
#include <fcntl.h>
#include <mqueue.h>
#include <errno.h>
// Expose non-variadic function requires 4 arguments.
mqd_t mq_open4(const char *name, int oflag, int mode, struct mq_attr *attr) {
	return mq_open(name, oflag, mode, attr);
}

*/
import "C"
import (
	"fmt"
	"time"
	"unsafe"
)

const (
	O_RDONLY = C.O_RDONLY
	O_WRONLY = C.O_WRONLY
	O_RDWR   = C.O_RDWR

	O_CLOEXEC  = C.O_CLOEXEC
	O_CREAT    = C.O_CREAT
	O_EXCL     = C.O_EXCL
	O_NONBLOCK = C.O_NONBLOCK

	S_IRUSR = C.S_IRUSR
	S_IWUSR = C.S_IWUSR
	S_IRGRP = C.S_IRGRP
	S_IWGRP = C.S_IWGRP

	// Based on Linux 3.5+
	MSGSIZE_MAX     = 16777216
	MSGSIZE_DEFAULT = MSGSIZE_MAX
)

var (
	MemoryAllocationError = fmt.Errorf("Memory Allocation Error")
)

type receiveBuffer struct {
	buf  *C.char
	size C.size_t
}

func newReceiveBuffer(bufSize int) (*receiveBuffer, error) {
	buf := (*C.char)(C.malloc(C.size_t(bufSize)))
	if buf == nil {
		return nil, MemoryAllocationError
	}

	return &receiveBuffer{
		buf:  buf,
		size: C.size_t(bufSize),
	}, nil
}

func (rb *receiveBuffer) free() {
	C.free(unsafe.Pointer(rb.buf))
}

func timeToTimespec(t time.Time) C.struct_timespec {
	return C.struct_timespec{
		tv_sec:  C.long(t.Unix()),
		tv_nsec: C.long(t.Nanosecond() % 1000000000),
	}
}

/*Errors
EACCES
    The queue exists, but the caller does not have permission to open it in the specified mode.
EACCES
    name contained more than one slash.
EEXIST
    Both O_CREAT and O_EXCL were specified in oflag, but a queue with this name already exists.
EINVAL
    O_CREAT was specified in oflag, and attr was not NULL, but attr->mq_maxmsg or attr->mq_msqsize was invalid. Both of these fields must be greater than zero. In a process that is unprivileged (does not have the CAP_SYS_RESOURCE capability), attr->mq_maxmsg must be less than or equal to the msg_max limit, and attr->mq_msgsize must be less than or equal to the msgsize_max limit. In addition, even in a privileged process, attr->mq_maxmsg cannot exceed the HARD_MAX limit. (See mq_overview(7) for details of these limits.)
EMFILE
    The process already has the maximum number of files and message queues open.
ENAMETOOLONG
    name was too long.
ENFILE
    The system limit on the total number of open files and message queues has been reached.
ENOENT
    The O_CREAT flag was not specified in oflag, and no queue with this name exists.
ENOENT
    name was just "/" followed by no other characters.
ENOMEM
    Insufficient memory.
ENOSPC
    Insufficient space for the creation of a new message queue. This probably occurred because the queues_max limit was encountered; see mq_overview(7).
*/
func mq_open(name string, oflag int, mode int, attr *MessageQueueAttribute) (int, error) {
	var cAttr *C.struct_mq_attr
	if attr != nil {
		cAttr = &C.struct_mq_attr{
			mq_flags:   C.long(attr.Flags),
			mq_maxmsg:  C.long(attr.MaxMsg),
			mq_msgsize: C.long(attr.MsgSize),
		}
	}
	// On success, mq_open() returns a message queue descriptor for use by other message queue functions.
	// On error, mq_open() returns (mqd_t) -1, with errno set to indicate the error.
	ret, err := C.mq_open4(C.CString(name), C.int(oflag), C.int(mode), cAttr)
	if ret == -1 {
		// return the syscall.Errno
		return 0, err
	}
	return int(ret), nil
}

/*Errors
EAGAIN
	The queue was full, and the O_NONBLOCK flag was set for the message queue description referred to by mqdes.
EBADF
	The descriptor specified in mqdes was invalid.
EINTR
	The call was interrupted by a signal handler; see signal(7).
EINVAL
	The call would have blocked, and abs_timeout was invalid, either because tv_sec was less than zero, or because tv_nsec was less than zero or greater than 1000 million.
EMSGSIZE
    msg_len was greater than the mq_msgsize attribute of the message queue.
ETIMEDOUT
    The call timed out before a message could be transferred.
*/
func mq_send(h int, data []byte, priority uint) error {
	byteStr := *(*string)(unsafe.Pointer(&data))
	// On success, mq_send() and mq_timedsend() return zero;
	// on error, -1 is returned, with errno set to indicate the error.
	rv, err := C.mq_send(C.int(h), C.CString(byteStr), C.size_t(len(data)), C.uint(priority))
	if rv == -1 {
		return err
	} else {
		return nil
	}
}

func mq_timedsend(h int, data []byte, priority uint, t time.Time) error {
	timeSpec := timeToTimespec(t)
	l := uint32(len(data))
	byteStr := *(*string)(unsafe.Pointer(&data))
	// On success, mq_send() and mq_timedsend() return zero;
	// on error, -1 is returned, with errno set to indicate the error.
	rv, err := C.mq_timedsend(C.int(h), C.CString(byteStr), C.size_t(l), C.uint(priority), &timeSpec)
	if rv == -1 {
		return err
	}
	return nil
}

/*Errors
EAGAIN
	The queue was empty, and the O_NONBLOCK flag was set for the message queue description referred to by mqdes.
EBADF
	The descriptor specified in mqdes was invalid.
EINTR
	The call was interrupted by a signal handler; see signal(7).
EINVAL
	The call would have blocked, and abs_timeout was invalid, either because tv_sec was less than zero, or because tv_nsec was less than zero or greater than 1000 million.
EMSGSIZE
    msg_len was less than the mq_msgsize attribute of the message queue.
ETIMEDOUT
    The call timed out before a message could be transferred.
*/
func mq_receive(h int, recvBuf *receiveBuffer) ([]byte, uint, error) {
	var msgPrio C.uint
	// On success, mq_receive() and mq_timedreceive() return the number of bytes in the received message;
	// On error, -1 is returned, with errno set to indicate the error.
	size, err := C.mq_receive(C.int(h), recvBuf.buf, recvBuf.size, &msgPrio)
	if size == -1 {
		return nil, 0, err
	}
	return C.GoBytes(unsafe.Pointer(recvBuf.buf), C.int(size)), uint(msgPrio), nil
}

func mq_timedreceive(h int, recvBuf *receiveBuffer, t time.Time) ([]byte, uint, error) {
	var (
		msgPrio  C.uint
		timeSpec = timeToTimespec(t)
	)
	// On success, mq_receive() and mq_timedreceive() return the number of bytes in the received message;
	// On error, -1 is returned, with errno set to indicate the error.
	size, err := C.mq_timedreceive(C.int(h), recvBuf.buf, recvBuf.size, &msgPrio, &timeSpec)
	if size == -1 {
		return nil, 0, err
	}

	return C.GoBytes(unsafe.Pointer(recvBuf.buf), C.int(size)), uint(msgPrio), nil
}

func mq_notify(h int, sigNo int) error {
	sigEvent := &C.struct_sigevent{
		sigev_notify: C.SIGEV_SIGNAL, // posix_mq supports only signal.
		sigev_signo:  C.int(sigNo),
	}
	// On success mq_notify() returns 0;
	// On error, -1 is returned, with errno set to indicate the error.
	rv, err := C.mq_notify(C.int(h), sigEvent)
	if rv == -1 {
		return err
	}
	return nil
}

func mq_close(h int) error {
	// On success mq_unlink() returns 0;
	// On error, -1 is returned, with errno set to indicate the error.
	rv, err := C.mq_close(C.int(h))
	if rv == -1 {
		return err
	} else {
		return nil
	}
}

func mq_unlink(name string) error {
	// On success mq_close() returns 0;
	// On error, -1 is returned, with errno set to indicate the error.
	rv, err := C.mq_unlink(C.CString(name))
	if rv == -1 {
		return err
	} else {
		return nil
	}
}
