package interfaces

import (
	"context"

	moduleEntity "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/module/entity"
	versionEntity "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/version/entity"
)

type ModuleClient interface {
	GetLastVersion(ctx context.Context, module moduleEntity.Module) (versionEntity.Version, error)
}
