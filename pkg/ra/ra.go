package ra

import (
	"slices"
	"sync"
)

type AdvAction int

const (
	Add AdvAction = iota
	Remove
)

type RouteAdv[R any] struct {
	Route  R
	Action AdvAction
	Reason string
}

// RaQueue manages concurrent incoming route advertisements from peers
type RaQueue[R any] struct {
	queues map[string][]RouteAdv[R] // peer -> queue
	mutex  sync.Mutex
}

func NewRaQueue[R any]() *RaQueue[R] {
	return &RaQueue[R]{
		queues: make(map[string][]RouteAdv[R]),
	}
}

func (rq *RaQueue[R]) PopAll(neighbor string) []RouteAdv[R] {
	rq.mutex.Lock()
	defer rq.mutex.Unlock()

	original := rq.queues[neighbor]
	if original == nil {
		return nil
	}
	cp := slices.Clone(original)
	delete(rq.queues, neighbor)

	return cp
}

func (rq *RaQueue[R]) BeginTx() *Tx[R] {
	return &Tx[R]{
		enqueued: make(map[string][]RouteAdv[R]),
		rq:       rq,
	}
}

type Tx[R any] struct {
	rq       *RaQueue[R]
	enqueued map[string][]RouteAdv[R] // Cached enqueues
}

func (tx *Tx[R]) Push(peer string, ra RouteAdv[R]) {
	tx.enqueued[peer] = append(tx.enqueued[peer], ra)
}

func (tx *Tx[R]) Commit() {
	// Aquire mutex
	tx.rq.mutex.Lock()
	defer tx.rq.mutex.Unlock()

	for k, v := range tx.enqueued {
		tx.rq.queues[k] = append(tx.rq.queues[k], v...)
	}
}
