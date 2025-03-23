package collection

type Stack[T any] struct {
	elements []T
}

func New[T any](capacity int) *Stack[T] {
	return &Stack[T]{
		elements: make([]T, 0, capacity),
	}
}

func (s *Stack[T]) Push(el T) {
	s.elements = append(s.elements, el)
}

func (s *Stack[T]) Len() int {
	return len(s.elements)
}

func (s *Stack[T]) Pop() (el T, ok bool) {
	if len(s.elements) == 0 {
		return el, false
	}

	end := len(s.elements) - 1
	el = s.elements[end]
	s.elements = s.elements[:end]

	return el, true
}
