package models

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"proxycenter/pkg/setting"
	clog "unknwon.dev/clog/v2"
	"xorm.io/core"
	"xorm.io/xorm"
)

// Engine represents a XORM engine or session.
type Engine interface {
	Delete(interface{}) (int64, error)
	Exec(string, ...interface{}) (sql.Result, error)
	Exist(...interface{}) (bool, error)
	Find(interface{}, ...interface{}) error
	Get(interface{}) (bool, error)
	Id(interface{}) *xorm.Session
	In(string, ...interface{}) *xorm.Session
	Insert(...interface{}) (int64, error)
	InsertOne(interface{}) (int64, error)
	Iterate(interface{}, xorm.IterFunc) error
	Query(string, ...interface{}) (sql.Result, error)
	Sql(string, ...interface{}) *xorm.Session
	Table(interface{}) *xorm.Session
	Where(interface{}, ...interface{}) *xorm.Session
}

var (
	x      *xorm.Engine
	tables []interface{}
	// HasEngine .
	HasEngine bool

	DbCfg struct {
		Type, Host, Name, User, Passwd, Path, SSLMode string
	}
	EnableSQLite3 bool
)

func init() {
	tables = append(tables,
		new(IP))
	gonicNames := []string{"SSL"}
	for _, name := range gonicNames {
		core.LintGonicMapper[name] = true
	}

}

// LoadDatabaseInfo .
func LoadDatabaseInfo() {
	DbCfg.Type = setting.DbType
	switch DbCfg.Type {
	case "sqlite3":
		setting.UseSQLite3 = true
		EnableSQLite3 = true
	case "mysql":
		setting.UseMySQL = true
	case "postgres":
		setting.UsePostgreSQL = true
	case "mssql":
		setting.UseMSSQL = true
	}
	DbCfg.Host = setting.Host
	DbCfg.Name = setting.Name
	DbCfg.User = setting.User
	if len(DbCfg.Passwd) == 0 {
		DbCfg.Passwd = setting.PassWd
	}
	DbCfg.SSLMode = setting.SslMode
	DbCfg.Path = setting.Path
}

// parsePostgreSQLHostPort parses given input in various forms defined in
// https://www.postgresql.org/docs/current/static/libpq-connect.html#LIBPQ-CONNSTRING
// and returns proper host and port number.
func parsePostgreSQLHostPort(info string) (string, string) {
	host, port := "127.0.0.1", "5432"
	if strings.Contains(info, ":") && !strings.HasSuffix(info, "]") {
		idx := strings.LastIndex(info, ":")
		host = info[:idx]
		port = info[idx+1:]
	} else if len(info) > 0 {
		host = info
	}
	return host, port
}

func parseMSSQLHostPort(info string) (string, string) {
	host, port := "127.0.0.1", "1433"
	if strings.Contains(info, ":") {
		host = strings.Split(info, ":")[0]
		port = strings.Split(info, ":")[1]
	} else if strings.Contains(info, ",") {
		host = strings.Split(info, ",")[0]
		port = strings.TrimSpace(strings.Split(info, ",")[1])
	} else if len(info) > 0 {
		host = info
	}
	return host, port
}

func getEngine() (*xorm.Engine, error) {
	connStr := ""
	Param := "?"
	if strings.Contains(DbCfg.Name, Param) {
		Param = "&"
	}
	//fmt.Println(DbCfg.Type)
	switch DbCfg.Type {
	case "mysql":
		if DbCfg.Host[0] == '/' { // looks like a unix socket
			connStr = fmt.Sprintf("%s:%s@unix(%s)/%s%scharset=utf8&parseTime=true",
				DbCfg.User, DbCfg.Passwd, DbCfg.Host, DbCfg.Name, Param)
		} else {
			connStr = fmt.Sprintf("%s:%s@tcp(%s)/%s%scharset=utf8&parseTime=true",
				DbCfg.User, DbCfg.Passwd, DbCfg.Host, DbCfg.Name, Param)
		}
	case "postgres":
		host, port := parsePostgreSQLHostPort(DbCfg.Host)
		if host[0] == '/' { // looks like a unix socket
			connStr = fmt.Sprintf("postgres://%s:%s@:%s/%s%ssslmode=%s&host=%s",
				url.QueryEscape(DbCfg.User), url.QueryEscape(DbCfg.Passwd), port, DbCfg.Name, Param, DbCfg.SSLMode, host)
		} else {
			connStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s%ssslmode=%s",
				url.QueryEscape(DbCfg.User), url.QueryEscape(DbCfg.Passwd), host, port, DbCfg.Name, Param, DbCfg.SSLMode)
		}
	case "mssql":
		host, port := parseMSSQLHostPort(DbCfg.Host)
		connStr = fmt.Sprintf("server=%s; port=%s; database=%s; user id=%s; password=%s;", host, port, DbCfg.Name, DbCfg.User, DbCfg.Passwd)
	case "sqlite3":
		if !EnableSQLite3 {
			return nil, errors.New("this binary version does not build support for SQLite3")
		}
		if err := os.MkdirAll(path.Dir(DbCfg.Path), os.ModePerm); err != nil {
			return nil, fmt.Errorf("Fail to create directories: %v", err)
		}
		connStr = "file:" + DbCfg.Path + "?cache=shared&mode=rwc"
	default:
		return nil, fmt.Errorf("Unknown database type: %s", DbCfg.Type)
	}
	return xorm.NewEngine(DbCfg.Type, connStr)
}

// NewTestEngine .
func NewTestEngine(x *xorm.Engine) (err error) {
	x, err = getEngine()
	if err != nil {
		return fmt.Errorf("Connect to database: %v", err)
	}

	x.SetMapper(core.GonicMapper{})
	return x.StoreEngine("InnoDB").Sync2(tables...)
}

// SetEngine .
func SetEngine() (err error) {
	x, err = getEngine()
	if err != nil {
		return fmt.Errorf("Fail to connect to database: %v", err)
	}

	x.SetMapper(core.GonicMapper{})

	// WARNING: for serv command, MUST remove the output to os.stdout,
	// so use log file to instead print to stdout.
	logger, err := clog.NewFileWriter(path.Join(setting.LogRootPath, "xorm.log"),
		clog.FileRotationConfig{
			Rotate:  setting.Rotate,
			Daily:   setting.RotateDaily,
			MaxSize: setting.MaxSize * 1024 * 1024,
			MaxDays: setting.MaxDays,
		})
	if err != nil {
		return fmt.Errorf("Fail to create 'xorm.log': %v", err)
	}

	if !setting.DebugMode {
		x.SetLogger(xorm.NewSimpleLogger2(logger, xorm.DEFAULT_LOG_PREFIX, xorm.DEFAULT_LOG_FLAG))
	} else {
		x.SetLogger(xorm.NewSimpleLogger(logger))
	}
	x.ShowSQL(true)
	return nil
}

// NewEngine .
func NewEngine() (err error) {
	if err = SetEngine(); err != nil {
		fmt.Println("链接失败了")
		return err
	}
	if err = x.StoreEngine("InnoDB").Sync2(tables...); err != nil {
		fmt.Println("数据库同步失败了")
		return fmt.Errorf("sync database struct error: %v", err)
	}

	return nil
}
