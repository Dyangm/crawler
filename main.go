package main

import (
	"github.com/Dyangm/crawler/command"
	"github.com/Dyangm/crawler/config"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

func main() {
	handler := command.NewHandler()
	handler.CommandHandler()
}

func init() {
	formatter := &log.JSONFormatter{}
	formatter.TimestampFormat = "2006-01-02 15:04:05.000"
	log.SetFormatter(formatter)
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
	c, err := config.GetConfig()
	if err != nil {
		log.WithField("error", err).Panic("get config failed")
		panic("get config failed")
	}

	level, err := log.ParseLevel(c.Log.Level)
	if err != nil {
		log.Error(err)
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(level)
	}

	if len(c.Log.Path) > 0 {
		f, err := os.OpenFile(c.Log.Path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			log.WithField("error", err).Panic("open log file failed")
			panic("open log file failed")
		}
		log.SetOutput(io.MultiWriter(os.Stdout, f))
	} else {
		log.SetOutput(os.Stdout)
	}

	log.WithField("config", c).Info("config info")

}
