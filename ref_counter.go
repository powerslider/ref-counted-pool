package refcount

import (
	"reflect"
	"sync"
	"sync/atomic"
)

// Interface following reference countable interface.
// We have provided inbuilt embeddable implementation of the reference countable entryPool.
// This interface just provides the extensibility for the implementation.
type ReferenceCountable interface {
	// Method to set the current instance
	SetInstance(i interface{})
	// Method to increment the reference count
	IncrementReferenceCount()
	// Method to decrement reference count
	DecrementReferenceCount()
}

type resetObjFunc func(interface{}) error

// Struct representing a reference.
// This struct is supposed to be embedded inside the object to be pooled.
// Along with that incrementing and decrementing the references is highly important specifically around goroutines.
type ReferenceCounter struct {
	count       *uint32
	destination *sync.Pool
	released    *uint32
	Instance    interface{}
	reset       resetObjFunc
	id          uint32
}

// Method to increment a reference.
func (r ReferenceCounter) IncrementReferenceCount() {
	atomic.AddUint32(r.count, 1)
}

// Method to decrement a reference.
// If the reference count goes to zero, the object is put back inside the entryPool.
func (r ReferenceCounter) DecrementReferenceCount() {
	if atomic.LoadUint32(r.count) == 0 {
		panic("this should not happen =>" + reflect.TypeOf(r.Instance).String())
	}
	decrementedCount := atomic.AddUint32(r.count, ^uint32(0))
	if decrementedCount == 0 {
		// Mark that object is released by incrementing the released count.
		atomic.AddUint32(r.released, 1)
		// Reset object to its zero values.
		if err := r.reset(r.Instance); err != nil {
			panic("error while resetting an instance => " + err.Error())
		}
		// Put object in the entry entryPool.
		r.destination.Put(r.Instance)
		// Stop tracking this current instance.
		r.Instance = nil
	}
}

// Method to set the current instance
func (r *ReferenceCounter) SetInstance(i interface{}) {
	r.Instance = i
}

