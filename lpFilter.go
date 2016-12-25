package main

// Filter provide a low-pass filter
type Filter struct {
	bufsize   int
	buf       []int
	sum       int
	nextindex int
}

// NewFilter return a new Filter given a UDPWriter and name
func NewFilter(size int) *Filter {
	return &Filter{
		bufsize:   size,
		buf:       make([]int, size),
		sum:       0,
		nextindex: 0,
	}
}

// Add a value into the filter
func (f *Filter) Add(vin int) int {
	f.sum -= f.buf[f.nextindex]                 // deduct the oldest value from the sum
	f.sum += vin                                // and add the newest
	f.buf[f.nextindex] = vin                    // save the newest in the buffer
	f.nextindex = (f.nextindex + 1) % f.bufsize // and move the index
	return f.sum / f.bufsize                    // return the average of the values in the buffer
}

// Reset the filter
func (f *Filter) Reset(v int) {
	f.sum = v * bufsize
	for i := 0; i < f.bufsize; i++ {
		f.buf[i] = v
	}
}
