package rabbit

type Consumer interface {
	Consume(chan error)
	Close()
}
