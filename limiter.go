package limiter

import "errors"

type Limiter interface {
	Take(int64) error
	fill() error
}

var (
 	ErrNoTEnoughToken = errors.New("Not enough tokens")
)
