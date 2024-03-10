package interfaces

type DatabaseConfig interface {
	InitDb() (interface{}, error)
	Close() error
	PingDb() error
}
