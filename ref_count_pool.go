package refcount

import (
	"sync"
	"sync/atomic"

)

type referenceCountWrapper func(referenceCounter ReferenceCounter) ReferenceCountable

// Struct representing the entryPool
type referenceCountedPool struct {
	entryPool       *sync.Pool
	refCountWrapper referenceCountWrapper
	returned        uint32
	allocated       uint32
	referenced      uint32
}

// Method to create a new entryPool
func NewReferenceCountedPool(refCountWrapper referenceCountWrapper, reset resetObjFunc) *referenceCountedPool {
	p := new(referenceCountedPool)
	p.entryPool = new(sync.Pool)
	p.entryPool.New = func() interface{} {
		// Incrementing allocated count
		atomic.AddUint32(&p.allocated, 1)
		c := refCountWrapper(ReferenceCounter{
			count:       new(uint32),
			destination: p.entryPool,
			released:    &p.returned,
			reset:       reset,
			id:          p.allocated,
		})
		return c
	}
	return p
}

// Method to get new object.
func (p *referenceCountedPool) Get() ReferenceCountable {
	c := p.entryPool.Get().(ReferenceCountable)
	c.SetInstance(c)
	atomic.AddUint32(&p.referenced, 1)
	c.IncrementReferenceCount()
	return c
}

// Method to return reference counted entryPool stats.
func (p *referenceCountedPool) Stats() map[string]interface{} {
	return map[string]interface{}{"allocated": p.allocated, "referenced": p.referenced, "returned": p.returned}
}
