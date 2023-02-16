package repository

import (
	"DouYIn/config"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// 定义全局的db对象，我们执行数据库操作主要通过他实现。
// 不用担心协程并发使用同样的db对象会共用同一个连接，db对象在调用他的方法的时候会从数据库连接池中获取新的连接
var Db *gorm.DB

var newLogger = logger.New(
	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
	logger.Config{
		SlowThreshold:             time.Duration(config.LOGGER_CONFIG.SlowThreshold), // 慢 SQL 阈值
		LogLevel:                  config.LOGGER_CONFIG.LogLevel,                     // 日志级别
		IgnoreRecordNotFoundError: config.LOGGER_CONFIG.IgnoreRecordNotFoundError,    // 忽略ErrRecordNotFound（记录未找到）错误
		Colorful:                  config.LOGGER_CONFIG.Colorful,                     // 禁用彩色打印
	},
)

// 包初始化函数，golang特性，每个包初始化的时候会自动执行init函数，这里用来初始化gorm。
func init() {
	//拼接下dsn参数, dsn格式可以参考上面的语法，这里使用Sprintf动态拼接dsn参数，因为一般数据库连接参数，我们都是保存在配置文件里面，需要从配置文件加载参数，然后拼接dsn。
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%ds", config.MYSQL_CONFIG.Username, config.MYSQL_CONFIG.Password, config.MYSQL_CONFIG.Host, config.MYSQL_CONFIG.Port, config.MYSQL_CONFIG.Dbname, config.MYSQL_CONFIG.Timeout)
	log.Println("dsn: ", dsn)
	//连接MYSQL, 获得DB类型实例，用于后面的数据库读写操作。
	// gorm - v2
	// gorm.Config 参考 https://gorm.io/docs/gorm_config.html
	_db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: config.MYSQL_CONFIG.SingularTable, // 如果设置为true,则`User`的默认表名为`user`,使用`TableName`设置的表名不受影响
		},
		Logger: newLogger,
	})
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
	Db = _db
	sqlDB, _ := Db.DB()

	//设置数据库连接池参数
	sqlDB.SetMaxOpenConns(config.MYSQL_CONFIG.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MYSQL_CONFIG.MaxIdleConns)
}
