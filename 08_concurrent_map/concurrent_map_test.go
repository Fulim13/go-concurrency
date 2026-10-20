package concurrentmap

import (
	"testing"
)

func TestConcurrentMap(t *testing.T) {
	cm1 := NewConcurrentMap[string, int](50)
	cm1.Store("FuLim", 24)
	if v, exists := cm1.Load("FuLim"); !exists {
		t.Fail()
	} else {
		if v != 24 {
			t.Fail()
		}
	}

	if _, exists := cm1.Load("HappyLim"); exists {
		t.Fail()
	}

	cm2 := NewConcurrentMap[int, bool](50)
	cm2.Store(24, true)
	if v, exists := cm2.Load(24); !exists {
		t.Fail()
	} else {
		if v != true {
			t.Fail()
		}
	}
	if _, exists := cm2.Load(99); exists {
		t.Fail()
	}
}
// cd 08_concurrent_map
// go test -v . -run '^TestConcurrentMap$'
