package don

var NewNilCheck = newNilCheck

func NewRequestPool[T any](v T) pool[T] {
	return newRequestPool(v)
}

func NewRequestDecoder[T any](v T) requestDecoder[T] { //nolint:revive
	return newRequestDecoder(v)
}

func PoolNew[T any](p pool[T]) any {
	return p.(*requestPool[T]).pool.New()
}

func PoolReset[T any](p pool[T], x *T) {
	p.(*requestPool[T]).reset(x)
}
