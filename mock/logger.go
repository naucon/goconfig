package mock

import "github.com/stretchr/testify/mock"

type LoggerMock struct {
	mock.Mock
}

func NewLoggerMock() *LoggerMock {
	return &LoggerMock{}
}

func (m *LoggerMock) Fatal(v ...interface{}) {
	m.Called(v)
}

func (m *LoggerMock) Fatalf(format string, v ...interface{}) {
	m.Called(format, v)
}

func (m *LoggerMock) Fatalln(v ...interface{}) {
	m.Called(v)
}

func (m *LoggerMock) Panic(v ...interface{}) {
	m.Called(v)
}

func (m *LoggerMock) Panicf(format string, v ...interface{}) {
	m.Called(format, v)
}

func (m *LoggerMock) Panicln(v ...interface{}) {
	m.Called(v)
}

func (m *LoggerMock) Print(v ...interface{}) {
	m.Called(v)
}

func (m *LoggerMock) Printf(format string, v ...interface{}) {
	m.Called(format, v)
}

func (m *LoggerMock) Println(v ...interface{}) {
	m.Called(v)
}
