package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/aquilahkj/maport/internal/log"
	"github.com/aquilahkj/maport/internal/version"
	"github.com/aquilahkj/maport/maportd"
	"github.com/judwhite/go-svc/svc"
)

type service struct {
	once sync.Once
	mpd  *maportd.MPD
	log.Logger
}

func main() {
	options := maportd.NewOptions()

	flagSet := getFlagSet()
	flagSet.Parse(os.Args[1:])

	if flagSet.Lookup("version").Value.(flag.Getter).Get().(bool) {
		fmt.Println(version.FullVersion("maportd"))
		os.Exit(0)
	}
	resolve(options, flagSet)

	logConfig := &log.Config{
		Target:       options.Log,
		Level:        options.LogLevel,
		PrintCaller:  options.LogCaller,
		Formatter:    options.LogFormat,
		DisableColor: false,
	}
	log.SetConfig(logConfig)
	logger := log.NewLogger("service")

	mpd, err := maportd.New(options)
	if err != nil {
		logger.Fatal("New maportd error, %s", err)
	}
	ser := &service{
		mpd:    mpd,
		Logger: logger,
	}
	if err := svc.Run(ser, syscall.SIGINT, syscall.SIGTERM); err != nil {
		ser.Fatal("Run service error, %s", err)
	}
}

func (ser *service) Init(env svc.Environment) error {
	if env.IsWindowsService() {
		dir := filepath.Dir(os.Args[0])
		return os.Chdir(dir)
	}
	return nil
}

func (ser *service) Start() error {
	go func() {
		err := ser.mpd.Run()
		if err != nil {
			ser.Error("Start service error, %s", err)
			ser.Stop()
			os.Exit(1)
		}
	}()
	return nil
}

func (ser *service) Stop() error {
	ser.once.Do(func() {
		ser.mpd.Exit()
	})
	return nil
}

func getFlagSet() *flag.FlagSet {
	flagSet := flag.NewFlagSet("maportd", flag.ExitOnError)
	flagSet.Bool("version", false, "Show the version.")
	flagSet.Int("port", 0, "The source address.")
	flagSet.String("dest", "", "The destination address.")
	flagSet.String("log", "", "Write log messages to this file. the default 'stdout'")
	flagSet.String("log-level", "", "The level of messages to log.")
	flagSet.Bool("log-caller", false, "Whether to log the caller.")
	flagSet.String("log-format", "", "The format of messages to log.")
	return flagSet
}

func resolve(options *maportd.Options, flagSet *flag.FlagSet) error {
	port := flagSet.Lookup("port").Value.(flag.Getter).Get().(int)
	dest := flagSet.Lookup("dest").Value.String()
	if port > 0 && dest != "" {
		options.MapInfos = []*maportd.MapInfo{&maportd.MapInfo{Port: port, DestAddr: dest}}
	}
	options.Log = flagSet.Lookup("log").Value.String()
	options.LogLevel = flagSet.Lookup("log-level").Value.String()
	options.LogCaller = flagSet.Lookup("log-caller").Value.(flag.Getter).Get().(bool)
	options.LogFormat = flagSet.Lookup("log-format").Value.String()
	return nil
}
