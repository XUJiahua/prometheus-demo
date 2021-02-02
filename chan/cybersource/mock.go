package cybersource

import (
	"errors"
	"math/rand"
	"time"
)

func NewMock() Service {
	return &mockService{}
}

type mockService struct {
}

// simulate:
// 1. response time
// 2. success rate
func (m mockService) Call() error {
	randomSleep()
	return randomError()
}

// sleep time range: 100ms ~ 300ms
func randomSleep() {
	ms := rand.Int()%200 + 100
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

var errorList = []error{
	errors.New("merchant not found"),
	errors.New("credit card has risk"),
}

// error rate: 20%
func randomError() error {
	index := rand.Int() % (len(errorList) * 5)
	if index < len(errorList) {
		return errorList[index]
	}
	return nil
}
