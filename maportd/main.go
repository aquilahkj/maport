package maportd

import (
	"errors"
	"fmt"

	"github.com/aquilahkj/maport/internal/log"
)

// MPD the main program model
type MPD struct {
	portMaps map[string]*Mapper
	log.Logger
}

// New create MappingPort instance
func New(options *Options) (*MPD, error) {
	if len(options.MapInfos) == 0 {
		return nil, errors.New("there are not any maport info")
	}
	mpd := make(map[string]*Mapper)
	logger := log.NewLogger("main")
	logger.Info("maportd create")
	for _, info := range options.MapInfos {
		mapper, err := NewMapper(info.Port, info.DestAddr)
		if err != nil {
			return nil, err
		}
		if _, ok := mpd[mapper.name]; ok {
			return nil, fmt.Errorf("the mapper %s is duplicate", mapper.name)
		}
		mpd[mapper.name] = mapper
	}
	return &MPD{mpd, logger}, nil
}

// Exit exit the MappingPort program
func (mpd *MPD) Exit() {
	mpd.Info("maportd exit")
	for key, mapper := range mpd.portMaps {
		mpd.Info("shuting down %s", key)
		mapper.Shutdown()
	}
}

// Run run the MappingPort program
func (mpd *MPD) Run() error {
	mpd.Info("maportd run")
	for key, mapper := range mpd.portMaps {
		mpd.Info("starting %s", key)
		if err := mapper.Start(); err != nil {
			return err
		}
	}
	return nil
}
