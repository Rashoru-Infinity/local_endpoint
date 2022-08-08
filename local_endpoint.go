package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"time"
)

var content []byte
var conf Config

type Config struct {
	interval time.Duration
	url string
	file string
}

func setConf() Config {
	var conf Config
	var err error
	interval := os.Getenv("INTERVAL")
	conf.interval, err = time.ParseDuration(interval)
	if err != nil {
		panic(err)
	}
	if conf.interval < time.Second * 10 {
		panic(errors.New("Interval must be greater than or equal to 10s"))
	}
	conf.url = os.Getenv("URL")
	conf.file = os.Getenv("FILENAME")
	return conf
}

func epsvc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Disposition", "Attachment; filename=" + conf.file)
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Content-Type", r.Header.Get("application/octet-stream"))
	http.ServeContent(w, r, os.Getenv("FILENAME"), time.Now(), bytes.NewReader(content))
}

func main() {
	conf = setConf()
	go func() {
		for {
			resp, err := http.Get(conf.url)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()
			content, err = io.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}
			time.Sleep(conf.interval)
		}
	}()
	http.HandleFunc("/", epsvc)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		panic(err)
	}
}