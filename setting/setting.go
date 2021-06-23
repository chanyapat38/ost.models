package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

type ApiMsg struct{}

type App struct {
	JwtSecret string
	AppId     string
	AppKey    string
	AppName   string

	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string

	HTTP200 string
	HTTP201 string
	HTTP206 string
	HTTP302 string
	HTTP400 string
	HTTP401 string
	HTTP403 string
	HTTP404 string
	HTTP405 string
	HTTP408 string
	HTTP415 string
	HTTP429 string
	HTTP500 string
	HTTP502 string
	HTTP503 string
	HTTP504 string
}

var AppSetting = &App{}

type Server struct {
	RunMode      string
	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

var cfg *ini.File

// Setup initialize the configuration instance
func (a ApiMsg) Setup() {
	var err error
	cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
	}

	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)

}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}

// func (a ApiMsg) EnvSetup() {
// 	godotenv.Load()
// }
