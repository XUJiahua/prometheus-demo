package cybersource

import (
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"math/rand"
	"time"
)

type MockService struct {
	LatenciesHistogram *prometheus.HistogramVec
}

// simulate:
// 1. response time
// 2. success rate
func (m MockService) Call(request *Request) error {
	var err error

	// use histogram record latency
	start := time.Now()
	defer func() {
		elapsed := float64(time.Since(start) / time.Millisecond)
		if err != nil {
			m.LatenciesHistogram.WithLabelValues(err.Error(), request.Op).Observe(elapsed)
		} else {
			m.LatenciesHistogram.WithLabelValues("success", request.Op).Observe(elapsed)
		}
	}()

	err = randomError()
	randomSleep(err)
	return err
}

var errorList = []error{
	errors.New("merchant_not_found"),
	errors.New("card_has_risk"),
	errors.New("system_error"),
}

func punishFactor(err error) int {
	for i, er := range errorList {
		if er == err {
			return i + 2
		}
	}
	// no punish
	return 1
}

// normal sleep time range: [20ms, 200ms)
// each err has punish factor
func randomSleep(err error) {
	ms := (rand.Int()%180 + 20) * punishFactor(err)
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

// error rate: 20%
func randomError() error {
	index := rand.Int() % (len(errorList) * 5)
	if index < len(errorList) {
		return errorList[index]
	}
	return nil
}
