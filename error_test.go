package goconfig

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigError_NewError(t *testing.T) {
	t.Run("TestConfigError_NewError_ShouldReturnMsg", func(t *testing.T) {
		err := NewConfigError(errConfigMissing, nil)

		assert.Error(t, err)
		assert.Equal(t, errConfigMissing, err.Error())
	})

	t.Run("TestConfigError_NewError_ShouldReturnMsgAndWrappedMsg", func(t *testing.T) {
		innerErrMsg := "inner error"
		innerErr := errors.New(innerErrMsg)
		expectedMsg := errConfigMissing + ": " + innerErrMsg
		err := NewConfigError(errConfigMissing, innerErr)

		assert.Error(t, err)
		assert.Equal(t, expectedMsg, err.Error())
	})
}
