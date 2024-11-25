//go:build !change

package lrucache

type Cache interface {
	// Get returns valueElem associated with the keyElem.
	//
	// The second valueElem is a bool that is true if the keyElem exists in the cache,
	// and false if not.
	Get(key int) (int, bool)
	// Set updates valueElem associated with the keyElem.
	//
	// If there is no keyElem in the cache new (keyElem, valueElem) pair is created.
	Set(key, value int)
	// Range calls function f on all elements of the cache
	// in increasing access time order.
	//
	// Stops earlier if f returns false.
	Range(f func(key, value int) bool)
	// Clear removes all keys and values from the cache.
	Clear()
}
