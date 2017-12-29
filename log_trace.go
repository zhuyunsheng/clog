package clog

import (
	"fmt"
	"log"
	"runtime"
)

const (
	LOG_DEBUG = 0
	LOG_INFO  = 1
	LOG_WARN  = 2
	LOG_ERROR = 3
	LOG_FATAL = 4
)

const (
	debug_str = "\033[0;32m[DEBUG]\033[0;0m"
	info_str  = "\033[0m[INFO]\033[0;0m"
	warn_str  = "\033[1;33m[WARN]\033[0;0m"
	error_str = "\033[0;31m[ERROR]\033[0;0m"
	fatal_str = "\033[0;31m[FATAL]\033[0;0m"
)

const (
	Llongfile  = 1 << iota
	Lshortfile
	Lfuncname
	LerrorExit
	Lfilemask  = Llongfile | Lshortfile | Lfuncname

	default_calldepth = 3
)

// debug
const (
	LDM_NONE  = 0
	LDM_ALL   = 1
	LDM_QUEST = 2
	LDM_LOGIN = 3
)

var (
	_log_level          = LOG_INFO
	_log_flag           = Lshortfile | Lfuncname // | LerrorExit
	_log_debug_modules  map[int]bool
	_log_debug_all      bool = true
	_log_current_module int  = LDM_ALL
)

func init() {
	_log_debug_modules = make(map[int]bool)
	_log_debug_modules[1] = true
}

func logPrefix(calldepth int) (ret string) {
	if (_log_flag & Lfilemask) != 0 {
		var ok bool
		var funcName string
		pc, file, line, ok := runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		if (_log_flag & Llongfile) != 0 {
			ret += file + fmt.Sprintf(":%d ", line)
		} else if (_log_flag & Lshortfile) != 0 {
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					file = file[i+1:]
					break
				}
			}
			ret += file + fmt.Sprintf(":%d ", line)
		}
		if (_log_flag & Lfuncname) != 0 {
			fc := runtime.FuncForPC(pc)
			if fc != nil {
				funcName = fc.Name() + "()"
			} else {
				funcName = "?()"
			}
			ret += funcName
		}
	}
	return
}

func SetFlag(flag int) {
	_log_flag = flag
}

func AddDebugModule(module int) {
	_log_debug_modules[module] = true
	if module == LDM_ALL {
		_log_debug_all = true
	}
}

func ClearDebugModules() {
	_log_debug_modules = make(map[int]bool)
	_log_debug_all = false
}

func SetCurrentDebugModule(m int) {
	_log_current_module = m
}

func ClearCurrentDebugModule() {
	_log_current_module = LDM_NONE
}

func NaLog(v ...interface{}) {
	if _log_level > LOG_INFO {
		return
	}
	str := formatLog(LOG_INFO, default_calldepth, v...)
	log.Println(str)
}

func NaLogf(format string, v ...interface{}) {
	if _log_level > LOG_INFO {
		return
	}
	str := formatLogf(LOG_INFO, default_calldepth, format, v...)
	log.Println(str)
}

func LogDebug(v ...interface{}) {
	if _log_level > LOG_DEBUG {
		return
	}
	if _log_debug_all == false {
		if _, exist := _log_debug_modules[_log_current_module]; !exist {
			return
		}
	}
	str := formatLog(LOG_DEBUG, default_calldepth, v...)
	log.Println(str)
}
func LogDebugc(calldepth int, v ...interface{}) {
	if _log_level > LOG_DEBUG {
		return
	}
	if _log_debug_all == false {
		if _, exist := _log_debug_modules[_log_current_module]; !exist {
			return
		}
	}
	str := formatLog(LOG_DEBUG, calldepth, v...)
	log.Println(str)
}
func LogDebugf(format string, v ...interface{}) {
	if _log_level > LOG_DEBUG {
		return
	}
	if _log_debug_all == false {
		if _, exist := _log_debug_modules[_log_current_module]; !exist {
			return
		}
	}
	str := formatLogf(LOG_DEBUG, default_calldepth, format, v...)
	log.Println(str)
}

func LogInfo(v ...interface{}) {
	if _log_level > LOG_INFO {
		return
	}
	str := formatLog(LOG_INFO, default_calldepth, v...)
	log.Println(str)
}
func LogInfoc(calldepth int, v ...interface{}) {
	if _log_level > LOG_INFO {
		return
	}
	str := formatLog(LOG_INFO, calldepth, v...)
	log.Println(str)
}
func LogInfof(format string, v ...interface{}) {
	if _log_level > LOG_INFO {
		return
	}
	str := formatLogf(LOG_INFO, default_calldepth, format, v...)
	log.Println(str)
}

func LogWarn(v ...interface{}) {
	if _log_level > LOG_WARN {
		return
	}
	str := formatLog(LOG_WARN, default_calldepth, v...)
	log.Println(str)
}
func LogWarnc(calldepth int, v ...interface{}) {
	if _log_level > LOG_WARN {
		return
	}
	str := formatLog(LOG_WARN, calldepth, v...)
	log.Println(str)
}
func LogWarnf(format string, v ...interface{}) {
	if _log_level > LOG_WARN {
		return
	}
	str := formatLogf(LOG_WARN, default_calldepth, format, v...)
	log.Println(str)
}

func LogErr(v ...interface{}) {
	if _log_level > LOG_ERROR {
		return
	}
	str := formatLog(LOG_ERROR, default_calldepth, v...)
	if logger := GetLogger("error"); logger != nil {
		logger.Println(str)
	}
	if (_log_flag & LerrorExit) != 0 {
		log.Fatalln(str)
	} else {
		log.Println(str)
	}
}
func LogErrc(calldepth int, v ...interface{}) {
	if _log_level > LOG_ERROR {
		return
	}
	str := formatLog(LOG_ERROR, calldepth, v...)
	if logger := GetLogger("error"); logger != nil {
		logger.Println(str)
	}
	if (_log_flag & LerrorExit) != 0 {
		log.Fatalln(str)
	} else {
		log.Println(str)
	}
}
func LogErrf(format string, v ...interface{}) {
	if _log_level > LOG_ERROR {
		return
	}
	str := formatLogf(LOG_ERROR, default_calldepth, format, v...)
	if logger := GetLogger("error"); logger != nil {
		logger.Println(str)
	}
	if (_log_flag & LerrorExit) != 0 {
		log.Fatalln(str)
	} else {
		log.Println(str)
	}
}

func LogFatal(v ...interface{}) {
	str := formatLog(LOG_FATAL, default_calldepth, v...)
	log.Println(str)
	if logger := GetLogger("error"); logger != nil {
		logger.Println(str)
	}
}
func LogFatalc(calldepth int, v ...interface{}) {
	str := formatLog(LOG_FATAL, calldepth, v...)
	log.Println(str)
	if logger := GetLogger("error"); logger != nil {
		logger.Println(str)
	}
}
func LogFatalf(format string, v ...interface{}) {
	str := formatLogf(LOG_FATAL, default_calldepth, format, v...)
	log.Println(str)
	if logger := GetLogger("error"); logger != nil {
		logger.Println(str)
	}
}

func formatLog(level int, calldepth int, v ...interface{}) string {
	var pre_str string
	var post_str string

	switch level {
	case LOG_DEBUG:
		pre_str = debug_str
	case LOG_INFO:
		pre_str = info_str
	case LOG_WARN:
		pre_str = warn_str
	case LOG_ERROR:
		pre_str = error_str
	case LOG_FATAL:
		pre_str = fatal_str
	}
	post_str = " [" + logPrefix(calldepth) + "]"
	return fmt.Sprint(pre_str, fmt.Sprint(v...), post_str)
}

func formatLogf(level int, calldepth int, format string, v ...interface{}) string {
	var pre_str string
	var post_str string

	switch level {
	case LOG_DEBUG:
		pre_str = debug_str
	case LOG_INFO:
		pre_str = info_str
	case LOG_WARN:
		pre_str = warn_str
	case LOG_ERROR:
		pre_str = error_str
	case LOG_FATAL:
		pre_str = fatal_str
	}
	post_str = " [" + logPrefix(calldepth) + "]"
	return fmt.Sprint(pre_str, fmt.Sprintf(format, v...), post_str)
}