package server

import (
	"sync/atomic"
)

type Counter struct {
	OldClientQPS    int64
	OldErrLogTotal  int64
	OldSlowLogTotal int64

	ClientConns  int64
	ClientQPS    int64
	ErrLogTotal  int64
	SlowLogTotal int64
}

func (counter *Counter) IncrClientConns() {
	atomic.AddInt64(&counter.ClientConns, 1)
}

func (counter *Counter) DecrClientConns() {
	atomic.AddInt64(&counter.ClientConns, -1)
}

func (counter *Counter) IncrClientQPS() {
	atomic.AddInt64(&counter.ClientQPS, 1)
}

func (counter *Counter) IncrErrLogTotal() {
	atomic.AddInt64(&counter.ErrLogTotal, 1)
}

func (counter *Counter) IncrSlowLogTotal() {
	atomic.AddInt64(&counter.SlowLogTotal, 1)
}

//flush the count per second
func (counter *Counter) FlushCounter() {
	atomic.StoreInt64(&counter.OldClientQPS, counter.ClientQPS)
	atomic.StoreInt64(&counter.OldErrLogTotal, counter.ErrLogTotal)
	atomic.StoreInt64(&counter.OldSlowLogTotal, counter.SlowLogTotal)

	atomic.StoreInt64(&counter.ClientQPS, 0)
}
