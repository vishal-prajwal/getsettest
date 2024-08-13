package jungleegames

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
)

func ConvertFromUnixMilliEpochToTime(s string) (time.Time, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}, errors.Wrapf(err, "Error parsing unix epoch milli time %s", s)
	}

	if i != 0 {
		return time.UnixMilli(i), nil
	}

	return time.Time{}, nil
}
