package models_test

type testRedis struct {
	getFn func(string) (string, error)
	setFn func(string, string, int) error
}

func (r *testRedis) Get(key string) (string, error) {
	if r.getFn == nil {
		panic("You have to override Get function in tests")
	}
	return r.getFn(key)
}

func (r *testRedis) Set(key, value string, ttl int) error {
	if r.setFn == nil {
		panic("You have to override Set function in tests")
	}
	return r.setFn(key, value, ttl)
}
