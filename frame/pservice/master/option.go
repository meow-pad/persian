package master

import (
	"errors"
	"github.com/meow-pad/persian/frame/pservice/cache"
	"github.com/meow-pad/persian/utils/timewheel"
)

type Options struct {
	SrvName                    string
	TickIntervalSec            int
	DistributionCacheKey       string
	ServiceId                  string
	DistributionCacheSignature [8]byte
	DistributionCacheExpireSec int64
	TWTimer                    *timewheel.TimeWheel
	Cache                      *cache.Cache
	Handler                    MSHandler
}

func (opts *Options) check() error {
	if opts.SrvName == "" {
		return errors.New("less SrvName")
	}
	if opts.TickIntervalSec <= 0 {
		return errors.New("less TickIntervalSec")
	}
	if opts.DistributionCacheKey == "" {
		return errors.New("less DistributionCacheKey")
	}
	if opts.ServiceId == "" {
		return errors.New("less ServiceId")
	}
	if opts.DistributionCacheSignature == [8]byte{} {
		return errors.New("less DistributionCacheSignature")
	}
	if opts.DistributionCacheExpireSec <= 0 {
		return errors.New("less DistributionCacheExpireSec")
	}
	if opts.TWTimer == nil {
		return errors.New("less TWTimer")
	}
	if opts.Cache == nil {
		return errors.New("less Cache")
	}
	if opts.Handler == nil {
		return errors.New("less Handler")
	}
	return nil
}

type Option func(*Options)

func WithSrvName(srvName string) Option {
	return func(options *Options) {
		options.SrvName = srvName
	}
}

func WithTickIntervalSec(tickIntervalSec int) Option {
	return func(options *Options) {
		options.TickIntervalSec = tickIntervalSec
	}
}

func WithDistributionCacheKey(distributionCacheKey string) Option {
	return func(options *Options) {
		options.DistributionCacheKey = distributionCacheKey
	}
}

func WithServiceId(serviceId string) Option {
	return func(options *Options) {
		options.ServiceId = serviceId
	}
}

func WithDistributionCacheSignature(distributionCacheSignature [8]byte) Option {
	return func(options *Options) {
		options.DistributionCacheSignature = distributionCacheSignature
	}
}

func WithDistributionCacheExpireSec(distributionCacheExpireSec int64) Option {
	return func(options *Options) {
		options.DistributionCacheExpireSec = distributionCacheExpireSec
	}
}

func WithTWTimer(twTimer *timewheel.TimeWheel) Option {
	return func(options *Options) {
		options.TWTimer = twTimer
	}
}

func WithCache(cache *cache.Cache) Option {
	return func(options *Options) {
		options.Cache = cache
	}
}

func WithHandler(handler MSHandler) Option {
	return func(options *Options) {
		options.Handler = handler
	}
}
