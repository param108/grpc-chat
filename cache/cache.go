package cache

type Cache interface {
	Connect() error
	SetKey(key string, val string) error
	SetKeyWithExpiry(key string, val string, ttlSeconds int) error
	GetKey(key string) (string, error)
}
