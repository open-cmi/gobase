package stack

import "container/list"

// Stack 栈结构
type Stack struct {
	list *list.List
}

// NewStack 创建新的栈
func NewStack() *Stack {
	list := list.New()
	return &Stack{list}
}

// Push 压入栈顶
func (stack *Stack) Push(value interface{}) {
	stack.list.PushBack(value)
}

// Pop 获取栈顶，弹出
func (stack *Stack) Pop() interface{} {
	e := stack.list.Back()
	if e != nil {
		stack.list.Remove(e)
		return e.Value
	}
	return nil
}

// Peak 获取栈顶，但是不弹出
func (stack *Stack) Peak() interface{} {
	e := stack.list.Back()
	if e != nil {
		return e.Value
	}

	return nil
}

// Len 获取栈的深度
func (stack *Stack) Len() int {
	return stack.list.Len()
}

// Empty 判断是否为空
func (stack *Stack) Empty() bool {
	return stack.list.Len() == 0
}
