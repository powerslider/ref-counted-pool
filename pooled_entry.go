package refcount

import (
"github.com/pkg/errors"
)

// Entry pool
var EntryPool = NewReferenceCountedPool(
	func(counter ReferenceCounter) ReferenceCountable {
		br := new(PooledEntry)
		br.ReferenceCounter = counter
		return br
	}, ResetEntry)

type PooledEntry struct {
	ReferenceCounter
	Field string
	Value interface{}
}

// Method to get new Entry
func AcquireEntry() *PooledEntry {
	return EntryPool.Get().(*PooledEntry)
}

// Method to reset Entry
// Used by reference countable pool
func ResetEntry(i interface{}) error {
	obj, ok := i.(*PooledEntry)
	if !ok {
		return errors.Errorf("illegal object sent to ResetEntry: %s", i)
	}
	obj.Reset()
	return nil
}

// Method to create new Entry
func NewEntry(field string, value interface{}) *PooledEntry {
	e := AcquireEntry()
	e.Field = field
	e.Value = value
	return e
}

func (e *PooledEntry) Reset() {
	e.Field = ""
	e.Value = nil
}

