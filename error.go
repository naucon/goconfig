package goconfig

import (
	"fmt"
)

const (
	errConfigEmptyEnv = "configuration: environments can't be empty"
	errConfigInvalid  = "configuration: config file can't be decoded, wrong format"
	errConfigMissing  = "configuration: missing config file"
)

type ConfigError struct {
	msg string
	Err error
}

func (e ConfigError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf(e.msg+": %v", e.Err)
	}
	return e.msg
}

func NewConfigError(msg string, err error) ConfigError {
	return ConfigError{
		msg: msg,
		Err: err,
	}
}
