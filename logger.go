package clog

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
	"sync"
	"strings"
)

var (
	LOGFILE_MAXSIZE_DEFAULT int64 = 50 << 20
	_logger_map             map[string]*log.Logger
	_logfile_maxsize        int64 = LOGFILE_MAXSIZE_DEFAULT
	_logfile_base           string
)

var logger_map_max sync.RWMutex

func init() {
	_logger_map = make(map[string]*log.Logger)
}

func GetLogger(typ string) *log.Logger {
	logger_map_max.RLock()
	if logger, ok := _logger_map[typ]; ok {
		logger_map_max.RUnlock()
		return logger
	}
	logger_map_max.RUnlock()

	var dir string

	if typ == "error" {
		tmpDir, filename := path.Split(_logfile_base)
		filenameArr := strings.Split(filename, ".")
		if len(filenameArr) >= 1 {
			dir = tmpDir + filenameArr[0] + "." + typ + ".log"
		} else {
			dir = tmpDir + typ + ".log"
		}
	} else {
		dir = _logfile_base + "." + typ
	}

	file, err := OpenLogFile(dir, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("error opening file %v\n", err)
		return nil
	}

	logger := log.New(file, "", log.LstdFlags)
	logger_map_max.Lock()
	_logger_map[typ] = logger
	logger_map_max.Unlock()
	return logger
}

func InitLogLevel(lv int)  {
	_log_level=lv
}

func InitLogger(logfile string, maxSize int64) {
	var err error
	var fullpath string

	// start with slash, just open
	if filepath.IsAbs(logfile) {
		fullpath = logfile
	} else {
		fullpath = path.Join(_base_path, "", logfile)
	}
	dir, filename := path.Split(logfile)
	if filename == "" {
		fullpath = path.Join(fullpath, "log")
	}
	_logfile_base = fullpath

	dir = path.Join(_base_path, "", dir)

	err = os.MkdirAll(dir, 0777)
	if err != nil {
		LogFatalf("mkdirAll err:%v,%v", err, dir)
		return
	}

	_logfile_maxsize = int64(maxSize) << 10 //k

	if _logfile_maxsize < (1 << 16) { //log file size min 64k
		_logfile_maxsize = LOGFILE_MAXSIZE_DEFAULT
	}

	startLogger(fullpath)
}
func startLogger(logfile string) {
	f, err := OpenLogFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("cannot open logfile %v\n", err)
		os.Exit(-1)
	}
	log.SetOutput(f)
}

func tmpLog(p *[]byte, format string, v ...interface{}) {
	*p = append([]byte(fmt.Sprintf(format, v...)), (*p)...)
}

type logFile struct {
	*os.File
}

func OpenLogFile(name string, flag int, perm os.FileMode) (file *logFile, err error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	lf := logFile{}
	lf.File = f
	return &lf, nil
}

func (f *logFile) Write(p []byte) (int, error) {
	fi, err := f.Stat()
	if err != nil {
		tmpLog(&p, "file.Stat err:%v.", err)
	}

	if fi.Size() >= _logfile_maxsize {
		now := time.Now().Format("2006_01_02-15_04_05")
		curFileName := f.Name()
		newFileName := fmt.Sprintf("%s.%s", f.Name(), now)

		err = os.Rename(curFileName, newFileName)
		if err != nil {
			tmpLog(&p, "[RAW] rename [%s] to [%s] err:%v\n",
				curFileName, newFileName, err)
		}

		newFile, err := os.OpenFile(curFileName,
			os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			tmpLog(&p, "[RAW] open file %s err:%v", curFileName, err)
		} else {
			f.File.Close()
			f.File = newFile
		}
	}

	return f.File.Write(p)
}
