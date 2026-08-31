package handler

import (
	"context"

	"miaomiaowux/internal/license"
	"miaomiaowux/internal/storage"
)

// licenseUserQuotaExceeded is kept as a compatibility helper for existing
// handler call sites. MEO is self-hosted and does not enforce account quotas.
func licenseUserQuotaExceeded(context.Context, *storage.TrafficRepository, *license.Manager) (string, bool) {
	return "", false
}
