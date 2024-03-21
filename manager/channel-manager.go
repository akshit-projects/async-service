package manager

type Channel[T any] interface {
	PushToChan(val T) bool
	IsChannelClosed() bool
	GetChannel() chan T
	CloseChannel()
}
type channelManager[T any] struct {
	ch       chan T
	isClosed bool
}

// GetChannel implements Channel.
func (c *channelManager[T]) GetChannel() chan T {
	return c.ch
}

// IsChannelClosed implements Channel.
func (c *channelManager[T]) IsChannelClosed() bool {
	return c.isClosed
}

// PushToChan implements Channel.
func (c *channelManager[T]) PushToChan(val T) bool {
	if c.isClosed {
		return false
	}
	c.ch <- val
	return true
}

// CloseChannel implements Channel.
func (c *channelManager[T]) CloseChannel() {
	if c.isClosed {
		return
	}
	c.isClosed = true
	close(c.ch)
}

func NewChannel[T any](size int16) Channel[T] {
	ch := make(chan T, size)
	return &channelManager[T]{
		ch:       ch,
		isClosed: false,
	}
}

func NewChannelWithChan[T any](ch chan T) Channel[T] {
	return &channelManager[T]{
		ch:       ch,
		isClosed: false,
	}
}
