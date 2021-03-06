package rollsum

import (
	"github.com/Redundancy/go-sync/circularbuffer"
)

func NewRollsum32(blocksize uint) *Rollsum32 {
	return &Rollsum32{
		Rollsum32Base: Rollsum32Base{
			blockSize: blocksize,
		},
		buffer: circularbuffer.MakeC2Buffer(int(blocksize)),
	}
}

// Uses 16bit internal values, 4 byte hashes
type Rollsum32 struct {
	Rollsum32Base
	buffer *circularbuffer.C2
}

// cannot be called concurrently
func (r *Rollsum32) Write(p []byte) (n int, err error) {
	ulen_p := uint(len(p))

	if ulen_p >= r.blockSize {
		// if it's really long, we can just ignore a load of it
		remaining := p[ulen_p-r.blockSize:]
		r.buffer.Write(remaining)
		r.Rollsum32Base.SetBlock(remaining)
	} else {
		b_len := r.buffer.Len()
		r.buffer.Write(p)
		evicted := r.buffer.Evicted()
		r.Rollsum32Base.AddAndRemoveBytes(p, evicted, b_len)
	}

	return len(p), nil
}

func (r *Rollsum32) BlockSize() int {
	return int(r.blockSize)
}

func (r *Rollsum32) Size() int {
	return 4
}

func (r *Rollsum32) Reset() {
	r.Rollsum32Base.Reset()
	r.buffer.Reset()
}

// Sum appends the current hash to b and returns the resulting slice.
// It does not change the underlying hash state.
// Note that this is to allow Sum() to reuse a preallocated buffer
func (r *Rollsum32) Sum(b []byte) []byte {
	if b != nil && cap(b)-len(b) >= 4 {
		p := len(b)
		b = b[:len(b)+4]
		r.Rollsum32Base.GetSum(b[p:])
		return b
	} else {
		result := []byte{0, 0, 0, 0}
		r.Rollsum32Base.GetSum(result)
		return append(b, result...)
	}
}

func (r *Rollsum32) GetLastBlock() []byte {
	return r.buffer.GetBlock()
}
