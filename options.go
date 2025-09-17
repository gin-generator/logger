package logger

type Option interface {
	apply(*Logger)
}

type optionFunc func(*Logger)

func (f optionFunc) apply(l *Logger) {
	f(l)
}

func NewLogger(opts ...Option) *Logger {
	logger := newDefaultLogger()

	for _, opt := range opts {
		opt.apply(logger)
	}

	logger.custom()
	return logger
}

func WithFileName(fileName string) Option {
	return optionFunc(func(l *Logger) {
		l.fileName = fileName
	})
}

func WithMaxSize(maxSize int) Option {
	return optionFunc(func(l *Logger) {
		l.maxSize = maxSize
	})
}

func WithMaxBackup(maxBackup int) Option {
	return optionFunc(func(l *Logger) {
		l.maxBackup = maxBackup
	})
}

func WithMaxAge(maxAge int) Option {
	return optionFunc(func(l *Logger) {
		l.maxAge = maxAge
	})
}

func WithCompress(compress bool) Option {
	return optionFunc(func(l *Logger) {
		l.compress = compress
	})
}

func WithLevel(level string) Option {
	return optionFunc(func(l *Logger) {
		l.level = level
	})
}
