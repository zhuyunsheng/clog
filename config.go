package clog

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"sync"
	"github.com/clog/goini"
)

var (
	_map          map[string]string
	_lock         sync.RWMutex
	_file         goini.File
	_base_path    = path.Join(path.Dir(os.Args[0]), "./conf")
	_default_file = path.Join(_base_path, "app.conf")
	_config_file  = flag.String("config", _default_file, "config filename")
	_force        = flag.Bool("f", false, "force")
)

func init() {
	flag.Parse()
	Reload()
}

func Get(name string) map[string]string {
	_lock.RLock()
	defer _lock.RUnlock()
	return _file.Section(name)
}

func FileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func Reload() {
	if FileExist(_default_file) == false {
		_base_path = path.Join(path.Dir(os.Args[0]), "../")
		if *_config_file == _default_file {
			_default_file = path.Join(_base_path, "app.conf")
			_config_file = &_default_file
		} else {
			_default_file = path.Join(_base_path, "app.conf")
		}
	}

	tpath := *_config_file
	log.Printf("Loading Config: %v.", tpath)
	defer log.Println("Config Load Completed.")

	_lock.Lock()
	_file = _load_config(tpath)
	_lock.Unlock()

	if v, ok := _file.Get("LOG", "log_level"); ok {
		if lvl, err := strconv.Atoi(v); err == nil {
			_log_level = lvl
		}
	}
	if v, ok := _file.Get("LOG", "log_flag"); ok {
		if lvl, err := strconv.Atoi(v); err == nil {
			_log_flag = lvl
		}
	}
}

func _load_config(tpath string) goini.File {
	var err error
	log.Println("Loading Config.")
	file, err := goini.LoadFile(tpath)
	if err != nil {
		fmt.Printf("load [%s] error: %s", tpath, err)
		os.Exit(-1)
	}
	return file
}

func GetLogLevel() int {
	return _log_level
}
