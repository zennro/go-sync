package blocksources

import (
	"fmt"
	"github.com/Redundancy/go-sync/patcher"
)

// errorWatcher is a small helper object
// sendIfSet will only return a channel if there is an error set
// so w.sendIfSet() <- w.Err() is always safe in a select statement
// even if there is no error set
type errorWatcher struct {
	errorChannel chan error
	lastError    error
}

func (w *errorWatcher) setError(e error) {
	if w.lastError != nil {
		panic("cannot set a new error when one is already set!")
	}
	w.lastError = e
}

func (w *errorWatcher) clear() {
	w.lastError = nil
}

func (w *errorWatcher) Err() error {
	return w.lastError
}

func (w *errorWatcher) sendIfSet() chan<- error {
	if w.lastError != nil {
		return w.errorChannel
	} else {
		return nil
	}
}

type pendingResponseHelper struct {
	responseChannel chan patcher.BlockReponse
	pendingResponse *patcher.BlockReponse
}

func (w *pendingResponseHelper) setResponse(r *patcher.BlockReponse) {
	if w.pendingResponse != nil {
		p := fmt.Sprintf("Setting a response when one is already set! Had startblock %v, got %v", r.StartBlock, w.pendingResponse.StartBlock)
		panic(p)
	}
	w.pendingResponse = r
}

func (w *pendingResponseHelper) clear() {
	w.pendingResponse = nil
}

func (w *pendingResponseHelper) Response() patcher.BlockReponse {
	if w.pendingResponse == nil {
		return patcher.BlockReponse{}
	}
	return *w.pendingResponse
}

func (w *pendingResponseHelper) sendIfPending() chan<- patcher.BlockReponse {
	if w.pendingResponse != nil {
		return w.responseChannel
	} else {
		return nil
	}

}

type UintSlice []uint

func (r UintSlice) Len() int {
	return len(r)
}

func (r UintSlice) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r UintSlice) Less(i, j int) bool {
	return r[i] < r[j]
}

type asyncResult struct {
	blockID uint
	data    []byte
	err     error
}

type queuedRequest struct {
	startBlockID uint
	startOffset  int64
	endOffset    int64
}

type queuedRequestList []queuedRequest

func (r queuedRequestList) Len() int {
	return len(r)
}

func (r queuedRequestList) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r queuedRequestList) Less(i, j int) bool {
	return r[i].startBlockID < r[j].startBlockID
}
