package posix_mq_test

import (
	"syscall"
	"testing"

	"github.com/skill215/posix_mq"
	"github.com/stretchr/testify/assert"
)

func TestOpenMQWithOutCreatePermission(t *testing.T) {
	oflag := posix_mq.O_WRONLY
	mqt, err := posix_mq.NewMessageQueue("/testName", oflag, 666, nil)
	assert.Nil(t, mqt)
	mqErr, ok := err.(syscall.Errno)
	assert.True(t, ok)
	assert.Equal(t, syscall.ENOENT, mqErr)
}

func TestOpenMQWithWrongName(t *testing.T) {
	oflag := posix_mq.O_WRONLY | posix_mq.O_CREAT
	mqt, err := posix_mq.NewMessageQueue("wrongName", oflag, 666, nil)
	assert.Nil(t, mqt)
	assert.NotNil(t, err)
	mqErr, ok := err.(syscall.Errno)
	assert.True(t, ok)
	assert.Equal(t, syscall.EINVAL, mqErr)
}
