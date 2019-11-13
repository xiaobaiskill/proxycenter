package initial

import (
	"proxycenter/pkg/models"
	"proxycenter/pkg/setting"
	clog "unknwon.dev/clog/v2"
)

func init() {
	// 加载配置文件
	ini_file := "conf/app.ini"
	ini, err := setting.NewContext(ini_file)
	if err != nil {
		panic("文件有误:" + err.Error())
	}

	// api host
	setting.AppName = ini.GetString("", "APP_NAME", "ProxyCenter")
	setting.AppAddr = ini.GetString("server", "HTTP_ADDR", "0.0.0.0")
	setting.AppPort = ini.GetString("server", "HTTP_PORT", "8082")
	// log
	setting.DebugMode = ini.GetBool("log", "DEBUG_MODE", false)
	setting.Log = ini.GetString("log", "LOG_FILE", "proxy.log")
	setting.LogLevel = ini.GetInt64("log", "LOG_LEVEL", 1)
	setting.LogRootPath = ini.GetString("log", "LOGROOTPATH", "")
	//db
	setting.DbType = ini.GetString("database", "DB_TYPE", "")
	setting.Host = ini.GetString("database", "HOST", "")
	setting.Name = ini.GetString("database", "NAME", "")
	setting.User = ini.GetString("database", "USER", "")
	setting.PassWd = ini.GetString("database", "PASSWD", "")
	setting.SslMode = ini.GetString("database", "SSL_MODE", "")
	setting.Path = ini.GetString("database", "PATH", "data/ProxyPool.db")
	setting.InstallLock = ini.GetBool("security", "INSTALL_LOCK", false)

	// log.xorm
	setting.Rotate = ini.GetBool("log.xorm", "ROTATE", true)
	setting.RotateDaily = ini.GetBool("log.xorm", "ROTATE_DAILY", true)
	setting.MaxSize = ini.GetInt64("log.xorm", "MAX_SIZE", 100)
	setting.MaxDays = ini.GetInt64("log.xorm", "MAX_DAYS", 3)

	setting.TimeOut = ini.GetInt64("request","TIMEOUT",1000)
}

func GlobalInit() {
	err := clog.NewConsole()
	if err != nil {
		panic("unable to create new logger: " + err.Error())
	}

	// 是否启动文件日志
	if !setting.DebugMode {
		err := clog.NewFile(clog.FileConfig{
			Level:    clog.Level(setting.LogLevel),
			Filename: setting.Log,
		})

		if err != nil {
			clog.Error("unable to create new logger: " + err.Error())
		}
	}

	models.LoadDatabaseInfo()

	if err := models.NewEngine(); err != nil {
		clog.Fatal("Fail to initialize ORM engine: %v", err)
	}
	models.HasEngine = true

}
