package stack

import (
	"errors"
)

type Stack struct {
	items []string
}

func New() *Stack {
	return &Stack{items: []string{}}
}

func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

func (s *Stack) Push(value string) {
	s.items = append(s.items, value)
}

func (s *Stack) Pop() (string, error) {
	if s.IsEmpty() {
		return "", errors.New("stack is empty")
	}
	top := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return top, nil
}

func (s *Stack) Top() (string, error) {
	if s.IsEmpty() {
		return "", errors.New("stack is empty")
	}
	return s.items[len(s.items)-1], nil
}
