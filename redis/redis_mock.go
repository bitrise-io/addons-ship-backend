package redis

// Mock ...
type Mock struct {
	GetStringFn func(string) (string, error)
	GetInt64Fn  func(string) (int64, error)
	SetFn       func(string, interface{}, int) error
}

// GetString ...
func (m *Mock) GetString(key string) (string, error) {
	if m.GetStringFn == nil {
		panic("You have to override GetString function in tests")
	}
	return m.GetStringFn(key)
}

// GetInt64 ...
func (m *Mock) GetInt64(key string) (int64, error) {
	if m.GetInt64Fn == nil {
		panic("You have to override GetInt64 function in tests")
	}
	return m.GetInt64Fn(key)
}

// Set ...
func (m *Mock) Set(key string, value interface{}, ttl int) error {
	if m.SetFn == nil {
		panic("You have to override Set function in tests")
	}
	return m.SetFn(key, value, ttl)
}
