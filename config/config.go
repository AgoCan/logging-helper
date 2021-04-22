package config

// 配置文件导入yaml文件是configstruct.go
//
// 配置文件可以使用 -c 的参数
// https://github.com/go-yaml/yaml

import (
	"flag"
	"log"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/spf13/viper"
)

// 设置配置文件的 环境变量
var (
	// LogDirector 日志目录
	LogDirector      string
	LogInfoFilename  string
	LogMaxSize       int
	LogMaxBackups    int
	LogMaxAge        int
	LogLevel         string
	LogInfoFilePath  string
	LogErrorFilePath string
	EsHost           string
	Sniff            bool
	TaskName         string
	EsIndex          string
)

// 获取文件绝对路径
func getCurrPath() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(1)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

// InitConfig 初始化配置项
func init() {
	var configFile = flag.String("c", "../config/config.yaml", "config fime path.")
	testing.Init()
	flag.Parse()

	if _, err := os.Stat(*configFile); os.IsNotExist(err) {
		newFile := path.Join(path.Dir(getCurrPath()), "config/config.yaml")
		configFile = &newFile
	}

	config := viper.New()
	config.AutomaticEnv()
	config.SetConfigFile(*configFile)
	err := config.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	LogDirector = config.GetString("LOG_DIRECTOR")
	if LogDirector == "" {
		LogDirector = path.Join(path.Dir(getCurrPath()), "log")
	}
	LogInfoFilename = config.GetString("LOG_INFO_FILENAME")
	LogMaxSize = config.GetInt("LOG_MAX_SIZE")
	LogMaxBackups = config.GetInt("LOG_MAX_BACKUPS")
	LogMaxAge = config.GetInt("LOG_MAX_AGE")
	LogLevel = config.GetString("LOG_LEVEL")
	LogInfoFilePath = path.Join(LogDirector, LogInfoFilename)
	LogErrorFilePath = path.Join(LogDirector, LogInfoFilename)
	EsHost = config.GetString("ES_HOST")
	Sniff = config.GetBool("SNIFF")
	TaskName = config.GetString("TASK_NAME")
	EsIndex = config.GetString("ES_INDEX")
}
