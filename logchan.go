package logchan

import (
	"log"
	"fmt"
	"strings"
)

type Level uint64

type Channel struct {
	Level Level
	Key byte
	Desc string
}

type Channels []Channel

const (
	LOG_NONE  Level = 0x0
	LOG_DEBUG Level =  0x800000000000000
	LOG_INFO  Level = 0x1000000000000000
	LOG_WARN  Level = 0x2000000000000000
	LOG_ERROR Level = 0x4000000000000000
	LOG_FATAL Level = 0x8000000000000000
	LOG_ALL   Level = 0xFFFFFFFFFFFFFFFF
)

const (
	CHANNEL_NONE byte  = '0'
	CHANNEL_DEBUG byte = 'D'
	CHANNEL_INFO byte  = 'I'
	CHANNEL_WARN byte  = 'W'
	CHANNEL_ERROR byte = 'E'
	CHANNEL_FATAL byte = 'F'
	CHANNEL_ALL byte   = 'A'
)

var defaultChannels = []Channel {
	Channel{LOG_NONE, CHANNEL_NONE, "none"},
	Channel{LOG_DEBUG, CHANNEL_DEBUG, "debug"},
	Channel{LOG_INFO, CHANNEL_INFO, "info"},
	Channel{LOG_WARN, CHANNEL_WARN, "warn"},
	Channel{LOG_ERROR, CHANNEL_ERROR, "error"},
	Channel{LOG_FATAL, CHANNEL_FATAL, "fatal"},
	Channel{LOG_ALL, CHANNEL_ALL, "all"}}

type Logger struct {
	level Level
	channels Channels
	chanmap map[byte]Channel
	bitmap map[Level]Channel
}

func NewLogger(ch Channels, def Level) *Logger {
	ret := new(Logger)
	ret.level = def

	for _, c := range defaultChannels {
		ret.chanmap[c.Key] = c
		ret.bitmap[c.Level] = c
	}

	for _, c := range ch {
		ret.chanmap[c.Key] = c
		ret.bitmap[c.Level] = c
	}

	ret.channels = ch
	return ret
}

func (logger *Logger) AtLevel (l Level) bool {
	return (l & logger.level) != 0
}

func (logger *Logger) LevelToPrefix (l Level) string {
	s := logger.LevelToString(l)
	if len(s) > 0 {
		s = fmt.Sprintf("[%s] ", s)
	}
	return s
}

func (logger *Logger) LevelToString(l Level) string {
	descs := make([]string,0)
	
	for i := 0; l != 0 && i < len(logger.channels); i++ {
		c := logger.channels[i]
		if (c.Level & l) != 0 {
			descs = append (descs, c.Desc)
			l = l&(^c.Level)
		}
	}
	return strings.Join(descs, ",")
}

func (logger *Logger) SetChannelsEz (which string, s string, setIfEmpty bool) bool {
	ret := true
	if len(s) > 0 || setIfEmpty {
		if s, e := logger.SetChannels(s); e == nil {
			log.Printf("Setting %s logging to '%s'\n", which, s);
		} else {
			log.Printf("Failed to set %s logging: %s\n", which, e);
			ret = false
		}
	}
	return ret
}

func (logger *Logger) SetChannels (s string) (newdesc string, e error) {
	var newlev Level = 0
	for _, c := range []byte(s) {
		if ch, found := logger.chanmap[c]; found {
			newlev |= ch.Level
		} else {
			e = fmt.Errorf("bad logger channel found: '%c'", c)
			break
		}
	}
	if e == nil {
		newdesc = logger.LevelToString(newlev)
		logger.level = newlev
	}
	return
}


func (logger *Logger) Printf(l Level, fmt string, v ...interface{}) {
	if logger.AtLevel (l) {
		s := logger.LevelToPrefix (l)
		if len(s) > 0 {
			fmt = s + fmt
		}
		log.Printf (fmt, v...)
	}
}

func (logger *Logger) Print(l Level, v ...interface{}) {
	if logger.AtLevel (l) {
		s := logger.LevelToPrefix (l)
		if len(s) > 0 {
			log.Print (s)
		}
		log.Print(v)
	}
}

func (logger *Logger) Println(l Level, v ...interface{}) {
	if logger.AtLevel (l) {
		s := logger.LevelToPrefix (l)
		if len(s) > 0 {
			log.Print (s)
		}
		log.Println (v...)
	}
}

var std = NewLogger(defaultChannels, LOG_ALL)

func AtLevel (l Level) bool {
	return (l & std.level) != 0
}

func LevelToPrefix (l Level) string {
	return std.LevelToPrefix(l)
}

func LevelToString(l Level) string {
	return std.LevelToString(l)
}

func SetChannelsEz (which string, s string, setIfEmpty bool) bool {
	return std.SetChannelsEz(which, s, setIfEmpty)
}

func SetChannels (s string) (newdesc string, e error) {
	return std.SetChannels(s)
}

func Printf(l Level, v ...interface{}) {
	std.Print(l, v...)
}

func Print(l Level, v ...interface{}) {
	std.Print(l, v...)
}

func Println(l Level, v ...interface{}) {
	std.Println(l, v...)
}

