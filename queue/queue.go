package queue

type Queue struct {
	queue []interface{}
	idx   int
}

func Create() *Queue {
	return &Queue{make([]interface{}, 0), -1}
}

func Add(q *Queue, item interface{}) {
	q.queue = append(q.queue, item)
	q.idx++
}

func Get(q *Queue) interface{} {
	curr := q.queue[0]
	q.queue = removeIndex(q.queue, 0)
	q.idx--
	return curr
}

func Size(q *Queue) int {
	return len(q.queue)
}

func IsEmpty(q *Queue) bool {
	return Size(q) == 0
}

func removeIndex(s []interface{}, index int) []interface{} {
	return append(s[:index], s[index+1:]...)
}
