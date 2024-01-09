package stack

import "fmt"

type Stack struct {
	items []string
}

func (s *Stack) Push(data string, items ...string) {
	s.items = append(s.items, data)

	for _, item := range items {
		s.items = append(s.items, item)
	}
}

func (s *Stack) Pop() string {
	if s.IsEmpty() {
		return ""
	}
	res := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return res
}

func (s *Stack) Top() (string, error) {
	if s.IsEmpty() {
		return "", fmt.Errorf("stack is empty")
	}
	return s.items[len(s.items)-1], nil
}

func (s *Stack) IsEmpty() bool {
	if len(s.items) == 0 {
		return true
	}
	return false
}
