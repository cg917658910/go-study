package main

import "fmt"

type MinStack struct {
	data    []int
	minData []int
}

func Constructor() MinStack {
	return MinStack{}
}

func (s *MinStack) Push(val int) {
	s.data = append(s.data, val)
	if len(s.minData) == 0 || val < s.minData[len(s.minData)-1] {
		s.minData = append(s.minData, val)
	} else {
		s.minData = append(s.minData, s.minData[len(s.minData)-1])
	}
}

func (s *MinStack) Pop() {
	s.data = s.data[:len(s.data)-1]
	s.minData = s.minData[:len(s.minData)-1]
}

func (s *MinStack) Top() int {
	return s.data[len(s.data)-1]
}
func (s *MinStack) printData() {
	fmt.Println("data: ", s.data)
	fmt.Println("minData: ", s.minData)
}
func (s *MinStack) GetMin() int {
	return s.minData[len(s.minData)-1]
}
func main() {
	stack := Constructor()
	stack.Push(3)
	stack.printData()
	stack.Push(1)
	stack.printData()

	stack.Push(2)
	stack.printData()

	println(stack.GetMin()) // 输出 2
	stack.Pop()
	println(stack.GetMin()) // 输出 2
	stack.Pop()
	println(stack.GetMin()) // 输出 5
}
