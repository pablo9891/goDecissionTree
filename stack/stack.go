package stack

type Stack struct {
	stack []interface{}
	idx   int
}

func Create() *Stack {
	return &Stack{make([]interface{}, 0), -1}
}

func Push(s *Stack, item interface{}) {
	s.stack = append(s.stack, item)
	s.idx++
}

func Pop(s *Stack) interface{} {
	curr := s.stack[s.idx]
	s.stack = removeIndex(s.stack, s.idx)
	s.idx--
	return curr
}

func Size(s *Stack) int {
	return len(s.stack)
}

func IsEmpty(s *Stack) bool {
	return Size(s) == 0
}

func removeIndex(s []interface{}, index int) []interface{} {
	return append(s[:index], s[index+1:]...)
}
