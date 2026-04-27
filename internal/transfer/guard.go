package transfer

import (
	"fmt"
	"strings"

	"dbear/internal/config"
)

var allowedDestinationHosts = map[string]struct{}{
	"localhost":            {},
	"host.docker.internal": {},
}

func ValidateDestination(destination config.Connection) error {
	if destination.Type == config.TypeSQLite {
		return nil
	}

	normalizedHost := strings.ToLower(strings.TrimSpace(destination.Host))
	if _, ok := allowedDestinationHosts[normalizedHost]; ok {
		return nil
	}

	return fmt.Errorf(
		"refusing to transfer to destination host %q: only %q and %q are allowed as destinations to prevent accidental writes to remote databases",
		destination.Host,
		"localhost",
		"host.docker.internal",
	)
}
