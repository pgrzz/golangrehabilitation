package mq

type Broker struct {
	ID   int
	Name string
	// other fields
}

// example function that returns a Broker object
func getBroker() (Broker, error) {
	broker := Broker{ID: 123, Name: "MyBroker"}
	return broker, nil
}
