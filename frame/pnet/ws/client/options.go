package client

func newOptions(opts ...Option) *Options {
	options := &Options{
		WriteQueueCap: 100,
	}
	for _, opt := range opts {
		opt(options)
	}
	return options
}

type Options struct {
	Name string
	// 同时提交写队列上限，队列满时异步写直接返回错误
	WriteQueueCap int
	// 最大消息长度
	MaxMessageLength int64
}

type Option func(*Options)

func WithName(value string) Option {
	return func(options *Options) {
		options.Name = value
	}
}

func WithWriteQueueCap(cap int) Option {
	return func(options *Options) {
		options.WriteQueueCap = cap
	}
}

func WithMaxMessageLength(value int64) Option {
	return func(options *Options) {
		options.MaxMessageLength = value
	}
}
