package adapters

import "context"

// ServerHTTP Interface for service to use
type ServerHTTP interface {
	ListenAndServe(port string) error
	Shutdown(ctx context.Context) error
}
