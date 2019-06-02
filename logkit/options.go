package logkit

const (
	LOG_LEVEL          = "INFO"
	LOG_PATH           = "log"
	LOG_MAX_SIZE       = 100
	LOG_MAX_BACKUPS    = 10
	LOG_MAX_AGE        = 3
	LOG_ENABLE_CONSOLE = false
	LOG_ENABLE_CALLER  = false
)

func newOptions() *Options {
	return &Options{
		level:         LOG_LEVEL,
		path:          LOG_PATH,
		maxSize:       LOG_MAX_SIZE,
		maxBackups:    LOG_MAX_BACKUPS,
		maxAge:        LOG_MAX_AGE,
		enableConsole: LOG_ENABLE_CONSOLE,
		enableCaller:  LOG_ENABLE_CALLER,
	}
}

type Options struct {
	path          string
	maxSize       int
	maxBackups    int
	maxAge        int
	level         string
	enableConsole bool
	enableCaller  bool
}

type Option func(*Options)

func Path(p string) Option {
	return func(o *Options) {
		o.path = p
	}
}

func MaxSize(m int) Option {
	return func(o *Options) {
		o.maxSize = m
	}
}

func MaxBackups(m int) Option {
	return func(o *Options) {
		o.maxBackups = m
	}
}

func MaxAge(m int) Option {
	return func(o *Options) {
		o.maxAge = m
	}
}

func Level(l string) Option {
	return func(o *Options) {
		o.level = l
	}
}

func EnableConsole(e bool) Option {
	return func(o *Options) {
		o.enableConsole = e
	}
}

func EnableCaller(e bool) Option {
	return func(o *Options) {
		o.enableCaller = e
	}
}
