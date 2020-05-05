package utils

import (
	"fmt"
	"testing"
	"time"
)

// go test -c -v utils/map-ttl_test.go utils/map-ttl.go
func TestMapT(t *testing.T) {
	fmt.Println(`begin`)
	mapT := NewMapT()
	mapT.Store(`test`, `test`, 1*time.Second)
	value, ok := mapT.Load(`test`)
	if !ok {
		t.Error(`我存了，但是拿取不到`)
		// return
	}
	fmt.Println(value)
	time.Sleep(3 * time.Second)
	value, ok = mapT.Load(`test`)
	if ok {
		fmt.Println(value)
		t.Error(`应该拿取不到才对`)
	}
	fmt.Println(value)
}
