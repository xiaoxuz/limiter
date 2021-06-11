package limiter

import "errors"

type Limiter interface {
	Take() error
	Cnt() int64
}

var (
 	ErrNoTEnoughToken = errors.New("Not enough tokens")
)
