package rabbit

import (
	"errors"
	"time"
)

const (
	MaxElapsedTime  = time.Minute
	InitialInterval = time.Second
	MaxInterval     = 20 * time.Second
	Multiplier      = 2
)

var (
	ErrReconnect           = errors.New("reconnect error")
	ErrNotValidAddressData = errors.New("not valid host ot port")
	ErrQueue               = errors.New("queue error")
	ErrChannelClosed       = errors.New("channel was closed")
)
