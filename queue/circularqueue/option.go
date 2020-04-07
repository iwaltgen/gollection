package circularqueue

// Option is create queue option
type Option func(queue *Queue)

// GrowthFactor set growth, shrink factor
// growth when elements full then growth by `factor` percent (0.0 means never grow)
// shrink when size is `factor` percent of capacity (0.0 means never shrink)
func GrowthFactor(growth, shrink float32) Option {
	return func(queue *Queue) {
		queue.growthFactor = growth
		queue.shrinkFactor = shrink
	}
}

// FixedSize is not grow queue size
// overwrite in order oldest element when element full
func FixedSize(size int) Option {
	return func(queue *Queue) {
		queue.growthFactor = 0.0
		queue.shrinkFactor = 0.0
		queue.guaranteedSize = size
		queue.elements = make([]interface{}, size)
	}
}

// GuaranteedSize is guaranteed queue minium size
func GuaranteedSize(size int) Option {
	return func(queue *Queue) {
		queue.guaranteedSize = size
		queue.elements = make([]interface{}, size)
	}
}
