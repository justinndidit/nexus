package broker

type Publisher interface {
	Publish(interface{}) error
}
