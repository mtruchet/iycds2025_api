package ping

// UseCase interface para ping
type UseCase interface {
	Execute() string
}

// PingImpl implementación del caso de uso ping
type PingImpl struct{}

// Execute ejecuta el caso de uso ping
func (u *PingImpl) Execute() string {
	return "pong"
}
