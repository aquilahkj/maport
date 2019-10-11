package maportd

// Options server options
type Options struct {
	MapInfos  []*MapInfo
	Log       string
	LogLevel  string
	LogCaller bool
	LogFormat string
}

// MapInfo map ports info
type MapInfo struct {
	Port     int
	DestAddr string
}

// NewOptions Create a options
func NewOptions() *Options {
	return &Options{}
}
