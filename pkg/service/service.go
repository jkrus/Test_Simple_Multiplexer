package service

type (
	// Service describe common service.
	Service interface {
		// Start tries start service.
		Start() error

		// Stop tries stop service.
		Stop() error
	}
)
