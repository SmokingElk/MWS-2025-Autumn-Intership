package proxygolang

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	moduleEntity "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/module/entity"
	"github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/module/errors"
	"github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/module/interfaces"
	versionEntity "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/version/entity"
)

const urlPrefix = "http://proxy.golang.org"

type ModuleClientProxyGolang struct {
}

func NewModuleProxyGolang() interfaces.ModuleClient {
	return &ModuleClientProxyGolang{}
}

func (c ModuleClientProxyGolang) GetLastVersion(
	ctx context.Context,
	module moduleEntity.Module,
) (versionEntity.Version, error) {
	url := fmt.Sprintf("%s/%s/@latest", urlPrefix, module.Name)

	resp, err := http.Get(url)

	if err != nil {
		return versionEntity.Version{}, fmt.Errorf("failed to get module: %w", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Fatal("failed to close response body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return versionEntity.Version{}, errors.ErrModuleNotFound
		default:
			return versionEntity.Version{},
				fmt.Errorf("failed to get module: server responded with status %d", resp.StatusCode)
		}
	}

	var response moduleResponseDTO

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return versionEntity.Version{}, fmt.Errorf("failed to parse response: %w", err)
	}

	version, err := versionEntity.NewVersion(response.Version)

	if err != nil {
		return versionEntity.Version{}, err
	}

	return version, nil
}
