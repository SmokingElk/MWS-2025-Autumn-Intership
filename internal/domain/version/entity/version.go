package entity

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/version/errors"
)

type Version struct {
	major int
	minor int
	patch int
}

func NewVersion(versionStr string) (Version, error) {
	cleanVersion := strings.TrimPrefix(strings.Split(versionStr, "-")[0], "v")

	parts := strings.Split(cleanVersion, ".")

	if len(parts) != 3 {
		return Version{}, errors.ErrBadVersionStr
	}

	partsInt := make([]int, 0, len(parts))

	for _, part := range parts {
		partInt, err := strconv.Atoi(part)

		if err != nil {
			return Version{}, errors.ErrBadVersionStr
		}

		partsInt = append(partsInt, partInt)
	}

	return Version{
		major: partsInt[0],
		minor: partsInt[1],
		patch: partsInt[2],
	}, nil
}

func (v Version) Less(other Version) bool {
	if v.major < other.major {
		return true
	}
	if v.major > other.major {
		return false
	}

	if v.minor < other.minor {
		return true
	}
	if v.minor > other.minor {
		return false
	}

	return v.patch < other.patch
}

func (v Version) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.major, v.minor, v.patch)
}
