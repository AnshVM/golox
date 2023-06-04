package Utils

type Stack[T any] struct {
	values []T
}

func NewStack[T any]() Stack[T] {
	return Stack[T]{values: []T{}}
}

func (s *Stack[T]) Push(value T) {
	s.values = append(s.values, value)
}

func (s *Stack[T]) Pop(value T) {
	if s.IsEmpty() {
		return
	}
	s.values = s.values[0 : len(s.values)-1]
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.values) == 0
}

func (s *Stack[T]) Get(i int) (T, error) {
	if len(s.values)-1 <= i {
		return s.values[i], nil
	}
	return *new(T), nil
}

func (s *Stack[T]) Peek() (T, error) {
	return s.Get(len(s.values) - 1)
}
