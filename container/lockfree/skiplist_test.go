package lockfree

import (
	"slices"
	"sync/atomic"
	"testing"
	"time"

	"github.com/nikmy/algo/testx/ptest"
	"github.com/stretchr/testify/assert"
)

func generateInts(start, end, step int) []int64 {
	ints := make([]int64, 0, (end-start)/step+1)
	for i := start; i < end; i += step {
		ints = append(ints, int64(i))
	}
	return ints
}

func TestSkipList_Find(t *testing.T) {
	t.Parallel()

	tests := [...]struct {
		name string
		list *SkipList[int64]
		elem int64
		want bool
	}{
		{
			name: "found max",
			list: MakeSkipList(generateInts(1, 9, 1)...),
			elem: 8,
			want: true,
		},
		{
			name: "found min",
			list: MakeSkipList(generateInts(1, 9, 1)...),
			elem: 1,
			want: true,
		},
		{
			name: "found middle",
			list: MakeSkipList(generateInts(1, 9, 1)...),
			elem: 4,
			want: true,
		},
		{
			name: "not found",
			list: MakeSkipList(generateInts(1, 9, 1)...),
			elem: 9,
			want: false,
		},
		{
			name: "found negative",
			list: MakeSkipList[int64](1, -2, -5, 7, -8),
			elem: -5,
			want: true,
		},
		{
			name: "not found negative",
			list: MakeSkipList[int64](1, -2, -5, 7, -8),
			elem: -4,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("dump:\n%s\n", tt.list.SDump())
			if got := tt.list.Lookup(tt.elem); got != tt.want {
				t.Errorf("List.Lookup(%d) = %v, want %v", tt.elem, got, tt.want)
			}
		})
	}
}

func TestSkipList_Insert(t *testing.T) {
	t.Parallel()

	l := MakeSkipList(1, 2, 3)
	if !l.Insert(4) {
		t.Error("cannot insert new element: expected true, got false")
	}
	if l.Insert(3) {
		t.Error("double insertion: expected false, got true")
	}
}

func TestSkipList_Delete(t *testing.T) {
	t.Parallel()

	l := MakeSkipList(1, 2, 3)
	if l.Delete(4) {
		t.Errorf("deleted element that does not exist")
	}
	if !l.Delete(1) {
		t.Errorf("cannot delete element that exists")
	}
	if l.Delete(1) {
		t.Errorf("double deleted element")
	}
}

func TestSkipList_Elements(t *testing.T) {
	t.Parallel()

	t.Run("empty list", func(t *testing.T) {
		list := NewSkipList[int64]()
		for elem := range list.Elements {
			t.Errorf("Unexpected element %d", elem)
		}
	})

	t.Run("traversal", func(t *testing.T) {
		l := MakeSkipList(generateInts(1, 11, 1)...)
		want := int64(1)
		for elem := range l.Elements {
			if elem != want {
				t.Errorf("List.Elements() = %v, want %v", elem, want)
			}
			want++
		}
	})
}

func TestSkipList_Visual(t *testing.T) {
	t.Parallel()

	tests := [...]struct {
		name string
		list *SkipList[int64]
	}{
		{
			name: "empty list",
			list: NewSkipList[int64](),
		},
		{
			name: "one element",
			list: MakeSkipList[int64](42),
		},
		{
			name: "small",
			list: MakeSkipList(generateInts(6, 0, -1)...),
		},
		{
			name: "medium",
			list: MakeSkipList(generateInts(-1, 11, 1)...),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("fmt: %s\n", tt.list)
			t.Logf("dump:\n%s\n", tt.list.SDump())
		})
	}
}

func TestSkipList_IsEmpty(t *testing.T) {
	t.Parallel()

	list := NewSkipList[int64]()
	assert.True(t, list.IsEmpty())

	list.Insert(1)
	assert.False(t, list.IsEmpty())
}

func TestSkipList_RaceFree(t *testing.T) {
	t.Run("find-find", func(t *testing.T) {
		ctrl := ptest.NewController(t)
		list := MakeSkipList(generateInts(2, 11, 2)...)
		for i := int64(2); i <= 10; i += 2 {
			ctrl.Spawn(200, func() {
				assert.Truef(t, list.Lookup(i), "missing element %d", i)
			})
		}
		for i := int64(1); i <= 9; i += 2 {
			ctrl.Spawn(200, func() { assert.Falsef(t, list.Lookup(i), "unexpected element %d", i) })
		}

		ctrl.Run(5 * time.Second)
	})

	t.Run("insert-insert", func(t *testing.T) {
		ctrl := ptest.NewController(t)
		list := NewSkipList[int64]()
		inserted := [50]atomic.Int64{}
		for i := int64(0); i < 50; i++ {
			ctrl.Spawn(400, func() {
				ok := list.Insert(i)
				if ok {
					n := int(inserted[i].Add(1))
					assert.Lessf(t, n, 2, "multiple insertion: %d", n)
				}
			})
		}

		ctrl.Run(5 * time.Second)

		assert.Equal(t, slices.Collect(list.Elements), generateInts(0, 50, 1))

		t.Logf("list after all operations:\n%s\n", list.SDump())
	})

	t.Run("find-insert", func(t *testing.T) {
		ctrl := ptest.NewController(t)
		list := NewSkipList[int64]()
		inserted := [50]atomic.Int64{}
		for i := int64(0); i < 50; i++ {
			ctrl.Spawn(400, func() {
				_ = list.Lookup(i)

				ok := list.Insert(i)
				if ok {
					n := int(inserted[i].Add(1))
					assert.Lessf(t, n, 2, "multiple insertion: %d", n)
				}

				ok = list.Lookup(i)
				assert.Truef(t, ok, "can't found element %d after insertion", i)
			})
		}

		ctrl.Run(5 * time.Second)

		t.Logf("list after all operations:\n%s\n", list.SDump())
	})

	t.Run("delete-delete", func(t *testing.T) {
		ctrl := ptest.NewController(t)
		list := MakeSkipList(generateInts(0, 60, 1)...)

		t.Logf("list before all operations:\n%s\n", list.SDump())

		for i := int64(0); i < 60; i++ {
			ctrl.Spawn(400, func() {
				_ = list.Delete(i * 2)
				ok := list.Delete(i * 2)
				if ok {
					assert.Failf(t, "double delete", "element %d deleted after delete", i*2)
				}

				_ = list.Delete(i*2 + 1)
				ok = list.Delete(i*2 + 1)
				if ok {
					assert.Failf(t, "double delete", "element %d deleted after delete", i*2+1)
				}
			})
		}

		ctrl.Run(5 * time.Second)

		t.Logf("list after all operations:\n%s\n", list.SDump())
	})

	t.Run("find-delete", func(t *testing.T) {
		ctrl := ptest.NewController(t)
		list := MakeSkipList(generateInts(0, 40, 2)...)

		t.Logf("list before all operations:\n%s\n", list.SDump())

		deleted := [50]atomic.Int64{}
		for i := int64(0); i < 50; i++ {
			ctrl.Spawn(400, func() {
				_ = list.Lookup(i * 2)

				ok := list.Delete(i*2 + 1)
				assert.Falsef(t, ok, "deleted element %d that does not exist", i*2+1)

				ok = list.Delete(i * 2)
				if ok {
					n := int(deleted[i].Add(1))
					assert.Lessf(t, n, 2, "multiple delete of %d: %d", i*2, n)
				}

				assert.Falsef(t, list.Lookup(i*2), "element %d found after deletion", i*2)
			})
		}

		ctrl.Run(5 * time.Second)

		t.Logf("list after all operations:\n%s\n", list.SDump())
	})

	t.Run("insert-delete", func(t *testing.T) {
		ctrl := ptest.NewController(t)
		list := NewSkipList[int64]()

		t.Logf("list before all operations:\n%s\n", list.SDump())

		counts := [100][2]atomic.Int64{}
		inc := func(i int64, b bool) {
			if b {
				counts[i][0].Add(1)
			}
		}
		dec := func(i int64, b bool) {
			if b {
				counts[i][1].Add(1)
			}
		}

		for i := int64(0); i < 50; i++ {
			ctrl.Spawn(10, func() {
				inc(i*2+1, list.Insert(i*2+1))
				inc(i*2, list.Insert(i*2))
				dec(i*2, list.Delete(i*2))
				dec(i*2+1, list.Delete(i*2+1))
			})
		}

		ctrl.Run(5 * time.Second)

		for i := range counts {
			inserts := counts[i][0].Load()
			deletes := counts[i][1].Load()
			assert.Equal(t, inserts, deletes, "inserts (expected) and deletes (actual) disbalanced")
			assert.Greaterf(t, inserts, int64(0), "no successful inserts / deletes")
		}

		assert.Len(t, slices.Collect(list.Elements), 0, "list must be empty")

		t.Logf("list after all operations:\n%s\n", list.SDump())
	})
}

func TestSkipList_Serializability(t *testing.T) {
	//  Serializability guarantees:
	//   ___________________________________
	//  |    Operations    |  Second result |
	//  |-----------------------------------|
	//  | Insert -> Lookup |     true       |
	//  | Insert -> Delete |     true       |
	//  | Insert -> Insert |     false      |
	//  | Delete -> Lookup |     false      |
	//  | Delete -> Delete |     false      |
	//  | Delete -> Insert |     false      |
	//  |__________________|________________|

	t.Parallel()

	t.Run("insert then lookup", func(t *testing.T) {
		ctrl := ptest.NewController(t)

		list := NewSkipList[int64]()
		for i := int64(0); i < 50; i++ {
			ctrl.Spawn(100, func() {
				_ = list.Insert(i)
				assert.True(t, list.Lookup(i))
			})
		}

		ctrl.Run(5 * time.Second)
	})

	t.Run("insert then insert", func(t *testing.T) {
		ctrl := ptest.NewController(t)

		list := NewSkipList[int64]()
		inserted := [51]atomic.Int64{}
		for i := int64(0); i < 50; i++ {
			ctrl.Spawn(100, func() {
				if list.Insert(i) {
					inserted[i].Add(1)
				}
				if list.Insert(i + 1) {
					inserted[i+1].Add(1)
				}
			})
		}

		ctrl.Run(5 * time.Second)

		for i := range inserted {
			assert.EqualValues(t, 1, inserted[i].Load(), "element %d is not inserted once", i)
		}
	})

	t.Run("delete then delete", func(t *testing.T) {
		ctrl := ptest.NewController(t)

		list := MakeSkipList(generateInts(0, 51, 1)...)
		deleted := [51]atomic.Int64{}

		for i := int64(0); i < 50; i++ {
			ctrl.Spawn(100, func() {
				if list.Delete(i) {
					deleted[i].Add(1)
				}
				if list.Delete(i + 1) {
					deleted[i+1].Add(1)
				}
			})
		}

		ctrl.Run(5 * time.Second)

		for i := range deleted {
			assert.EqualValues(t, 1, deleted[i].Load(), "element %d is not deleted once", i)
		}
	})

	t.Run("insert then delete", func(t *testing.T) {
		ctrl := ptest.NewController(t)

		list := NewSkipList[int64]()
		inserted := [51]atomic.Int64{}
		deleted := [51]atomic.Int64{}
		for i := int64(0); i < 50; i++ {
			ctrl.Spawn(100, func() {
				if list.Insert(i) {
					inserted[i].Add(1)
				}
				if list.Delete(i) {
					deleted[i].Add(1)
				}
			})
		}

		ctrl.Run(5 * time.Second)

		for i := range len(inserted) {
			assert.EqualValues(t, inserted[i].Load(), deleted[i].Load(), "inserted (expected) != deleted (actual)", i)
		}
	})

	t.Run("delete then insert", func(t *testing.T) {
		ctrl := ptest.NewController(t)

		list := MakeSkipList(generateInts(0, 51, 1)...)
		inserted := [51]atomic.Int64{}
		deleted := [51]atomic.Int64{}
		for i := int64(0); i < 50; i++ {
			ctrl.Spawn(100, func() {
				if list.Delete(i) {
					deleted[i].Add(1)
				}
				if list.Insert(i) {
					inserted[i].Add(1)
				}
			})
		}

		ctrl.Run(5 * time.Second)

		for i := range len(inserted) {
			assert.EqualValues(t, deleted[i].Load(), inserted[i].Load(), "inserted (expected) != deleted (actual)", i)
		}
	})

	t.Run("delete then lookup", func(t *testing.T) {
		ctrl := ptest.NewController(t)

		list := MakeSkipList(generateInts(0, 51, 1)...)
		for i := int64(0); i < 50; i++ {
			ctrl.Spawn(100, func() {
				_ = list.Delete(i)
				assert.False(t, list.Lookup(i))
			})
		}

		ctrl.Run(5 * time.Second)
	})
}
