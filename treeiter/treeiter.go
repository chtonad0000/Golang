package treeiter

type Tree[T any] interface {
	Left() *T
	Right() *T
}

type ValuesNode[T any] struct {
	value       T
	left, right *ValuesNode[T]
}

func (t *ValuesNode[T]) Left() *ValuesNode[T] {
	return t.left
}

func (t *ValuesNode[T]) Right() *ValuesNode[T] {
	return t.right
}

func DoInOrder[T any](tree *T, action func(node *T)) {
	if tree == nil {
		return
	}
	if t, ok := interface{}(tree).(Tree[T]); ok {
		DoInOrder(t.Left(), action)
		action(tree)
		DoInOrder(t.Right(), action)
	}
}
