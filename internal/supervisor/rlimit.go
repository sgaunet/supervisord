// +build !windows,!freebsd

package supervisor

import (
	"syscall"

	apperrors "github.com/sgaunet/supervisord/internal/errors"
)

func (s *Supervisor) checkRequiredResources() error {
	if minfds, vErr := s.getMinRequiredRes("minfds"); vErr == nil {
		return s.checkMinLimit(syscall.RLIMIT_NOFILE, "NOFILE", minfds)
	}
	if minprocs, vErr := s.getMinRequiredRes("minprocs"); vErr == nil {
		// RPROC = 6
		return s.checkMinLimit(6, "NPROC", minprocs)
	}
	return nil
}

func (s *Supervisor) getMinRequiredRes(resourceName string) (uint64, error) {
	if entry, ok := s.config.GetSupervisord(); ok {
		intVal := entry.GetInt(resourceName, 0)
		if intVal < 0 {
			return 0, apperrors.NewNegativeValueError(resourceName) //nolint:wrapcheck // Internal error type with context
		}
		value := uint64(intVal)
		if value > 0 {
			return value, nil
		}
		return 0, apperrors.NewNoSuchKeyError(resourceName) //nolint:wrapcheck // Internal error type with context
	}
	return 0, apperrors.ErrNoSupervisordSection
}

func (s *Supervisor) checkMinLimit(resource int, resourceName string, minRequiredSource uint64) error {
	var limit syscall.Rlimit

	if syscall.Getrlimit(resource, &limit) != nil {
		return apperrors.NewFailedToGetLimitError(resourceName) //nolint:wrapcheck // Internal error type with context
	}

	if minRequiredSource > limit.Max {
		//nolint:gosec // G115: Conversion validated by system limits
		return apperrors.NewLimitExceedsHardError(resourceName, int64(minRequiredSource), int64(limit.Max)) //nolint:wrapcheck // Internal error type with context
	}

	if limit.Cur >= minRequiredSource {
		return nil
	}

	limit.Cur = limit.Max
	if syscall.Setrlimit(syscall.RLIMIT_NOFILE, &limit) != nil {
		//nolint:gosec // G115: Conversion validated by system limits
		return apperrors.NewFailedToSetLimitError(resourceName, int64(limit.Cur)) //nolint:wrapcheck // Internal error type with context
	}
	return nil
}
