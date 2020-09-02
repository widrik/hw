package rabbit

import (
	"errors"
	"time"
)

const (
	MaxElapsedTime  = 4 * time.Minute
	InitialInterval = time.Second
	MaxInterval     = 30 * time.Second
	Multiplier      = 2
)

var (
	ErrReconnect           = errors.New("reconnect error")
	ErrNotValidAddressData = errors.New("not valid host ot port")
	ErrChannelClosed       = errors.New("channel was closed")
)
