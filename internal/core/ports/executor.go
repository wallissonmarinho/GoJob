package ports

// Executor is the main port that defines how to execute a command
type Executor interface {
	Execute() error
}

// ServiceCommand is a port that defines the interface for a service command
type ServiceCommand interface {
	Execute() error
}
