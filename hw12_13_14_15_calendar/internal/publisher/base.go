package publisher

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
	ErrQueue               = errors.New("queue error")
	ErrChannelClosed       = errors.New("channel was closed")
)
