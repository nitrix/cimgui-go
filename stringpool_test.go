//go:build cgo

package imgui

import "testing"

func TestStringPoolStoreCStringArray(t *testing.T) {
	var pool StringPool
	items := pool.StoreCStringArray([]string{"one", "two"})
	defer pool.Reset()

	if items == nil {
		t.Fatal("expected C string array")
	}
}

func TestStringPoolStoreCStringArrayEmpty(t *testing.T) {
	var pool StringPool
	if got := pool.StoreCStringArray(nil); got != nil {
		t.Fatalf("empty array = %v, want nil", got)
	}
}

func BenchmarkStringPoolStoreCString(b *testing.B) {
	var pool StringPool
	for i := 0; i < b.N; i++ {
		_ = pool.StoreCString("button label")
		if i%1024 == 0 {
			pool.Reset()
		}
	}
	pool.Reset()
}
