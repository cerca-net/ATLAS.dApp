package vm

// Storage provides persistent key-value storage for contracts.
type Storage struct {
	data map[string]interface{}
}

// NewStorage creates a new contract storage instance.
func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]interface{}),
	}
}

// Get retrieves a value from storage.
func (s *Storage) Get(key string) (interface{}, bool) {
	val, ok := s.data[key]
	return val, ok
}

// Set stores a value in storage.
func (s *Storage) Set(key string, value interface{}) {
	s.data[key] = value
} 