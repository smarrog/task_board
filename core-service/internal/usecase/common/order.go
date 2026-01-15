package common

type IntLike interface {
	~int
}

type Shift[T IntLike] struct {
	FromPosition T
	ToPosition   T
	Delta        int
}

func CalcShift[T IntLike](from, to T) (Shift[T], bool) {
	if to == from {
		return Shift[T]{}, false
	}

	if to < from {
		return Shift[T]{
			FromPosition: to,
			ToPosition:   from - 1,
			Delta:        1,
		}, true
	}

	return Shift[T]{
		FromPosition: from + 1,
		ToPosition:   to,
		Delta:        -1,
	}, true
}

func Clamp[T IntLike](pos, min, max T) T {
	if pos < min {
		return min
	}
	if pos > max {
		return max
	}
	return pos
}
