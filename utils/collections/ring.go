package collections

func NewRing[T any](capacity int) *Ring[T] {
	return &Ring[T]{
		elements: make([]T, capacity),
		capacity: capacity,
		front:    0,
		rear:     0,
	}
}

type Ring[T any] struct {
	elements []T // 数据
	capacity int
	front    int // 前指针,负责弹出数据
	rear     int // 尾指针,负责添加数据
	empty    T   // 空值
}

// Push
//
//	@Description: 入队操作
//	@receiver ring
//	@param data
//	@return bool
func (ring *Ring[T]) Push(data T) bool {
	if (ring.rear+1)%ring.capacity == ring.front {
		// 队列已满
		return false
	}
	ring.elements[ring.rear] = data // 放入队列尾部
	ring.rear = (ring.rear + 1) % ring.capacity
	return true
}

// Pop
//
//	@Description: 出队操作
//	@receiver ring
//	@return value
func (ring *Ring[T]) Pop() (value T) {
	if ring.rear == ring.front {
		return value
	}
	elem := ring.elements[ring.front]
	ring.elements[ring.front] = ring.empty
	ring.front = (ring.front + 1) % ring.capacity
	return elem
}

// Length
//
//	@Description: 队列长度
//	@receiver ring
//	@return int
func (ring *Ring[T]) Length() int {
	return (ring.rear - ring.front + ring.capacity) % ring.capacity
}
