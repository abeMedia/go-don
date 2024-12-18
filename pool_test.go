package don_test

import (
	"reflect"
	"testing"

	"github.com/abemedia/go-don"
)

func TestRequestPool(t *testing.T) {
	type item struct {
		String  string
		Pointer *string
	}

	t.Run("Nil", func(t *testing.T) {
		var zero any
		pool := don.NewRequestPool(zero)

		pool.Put(pool.Get())

		if !reflect.DeepEqual(&zero, pool.Get()) {
			t.Fatal("should be zero value")
		}
	})

	t.Run("Struct", func(t *testing.T) {
		zero := item{}
		pool := don.NewRequestPool(zero)

		for i := 0; i < 100; i++ {
			v := pool.Get()
			v.String = "test"
			v.Pointer = &v.String
			pool.Put(v)
		}

		for i := 0; i < 100; i++ {
			if !reflect.DeepEqual(&zero, pool.Get()) {
				t.Fatal("should be zero value")
			}
		}
	})

	t.Run("Pointer", func(t *testing.T) {
		zero := &item{}
		pool := don.NewRequestPool(zero)

		for i := 0; i < 100; i++ {
			p := pool.Get()
			v := *p
			v.String = "test"
			v.Pointer = &v.String
			pool.Put(p)
		}

		for i := 0; i < 100; i++ {
			if !reflect.DeepEqual(&zero, pool.Get()) {
				t.Fatal("should be zero value")
			}
		}
	})
}

func BenchmarkPool(b *testing.B) {
	type item struct {
		String  string
		Pointer *string
	}
	b.Run("Struct", func(b *testing.B) {
		pool := don.NewRequestPool(item{})
		for i := 0; i < b.N; i++ {
			pool.Put(pool.Get())
		}
	})
	b.Run("Pointer", func(b *testing.B) {
		pool := don.NewRequestPool(new(item))
		for i := 0; i < b.N; i++ {
			pool.Put(pool.Get())
		}
	})
}

func BenchmarkPool_New(b *testing.B) {
	type item struct {
		String  string
		Pointer *string
	}
	b.Run("Struct", func(b *testing.B) {
		pool := don.NewRequestPool(item{})
		for i := 0; i < b.N; i++ {
			don.PoolNew(pool)
		}
	})
	b.Run("Pointer", func(b *testing.B) {
		pool := don.NewRequestPool(new(item))
		for i := 0; i < b.N; i++ {
			don.PoolNew(pool)
		}
	})
}

func BenchmarkPool_Reset(b *testing.B) {
	type item struct {
		String  string
		Pointer *string
	}
	b.Run("Struct", func(b *testing.B) {
		pool := don.NewRequestPool(item{})
		x := item{}
		for i := 0; i < b.N; i++ {
			don.PoolReset(pool, &x)
		}
	})
	b.Run("Pointer", func(b *testing.B) {
		pool := don.NewRequestPool(new(item))
		x := new(item)
		for i := 0; i < b.N; i++ {
			don.PoolReset(pool, &x)
		}
	})
}
