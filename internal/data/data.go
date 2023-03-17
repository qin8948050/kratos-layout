package data

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/qin8948050/kratos-layout/internal/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGormDB, NewGreeterRepo)

// Data .
type Data struct {
	gormDB *gorm.DB
}

// NewData .
func NewData(logger log.Logger, db *gorm.DB) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{gormDB: db}, cleanup, nil
}

func NewGormDB(c *conf.Data, logger log.Logger) (*gorm.DB, error) {
	dsn := c.Database.Source
	helper := log.NewHelper(logger)
	db, err := gorm.Open(mysql.Open(dsn), gormConfig(helper))
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(50)
	sqlDB.SetMaxOpenConns(150)
	sqlDB.SetConnMaxLifetime(time.Second * 25)
	return db, err
}

type config struct {
	SlowThreshold time.Duration
	LogLevel      logger.LogLevel
}

func gormConfig(log *log.Helper) *gorm.Config {
	config := &gorm.Config{
		Logger: encoder(config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      logger.Silent,
		}, log),
	}
	return config
}

type customLogger struct {
	config
	log                                 *log.Helper
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
	ctx                                 context.Context
}

func encoder(config config, log *log.Helper) logger.Interface {
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)
	return &customLogger{
		log:          log,
		config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

// LogMode log mode
func (c *customLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *c
	newLogger.LogLevel = level
	return &newLogger
}

// Info print info
func (c *customLogger) Info(ctx context.Context, message string, data ...interface{}) {
	if c.LogLevel >= logger.Info {
		c.ctx = ctx
		c.Printf(c.infoStr+message, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (c *customLogger) Warn(ctx context.Context, message string, data ...interface{}) {
	if c.LogLevel >= logger.Warn {
		c.ctx = ctx
		c.Printf(c.warnStr+message, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (c *customLogger) Error(ctx context.Context, message string, data ...interface{}) {
	if c.LogLevel >= logger.Error {
		c.ctx = ctx
		c.Printf(c.errStr+message, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
func (c *customLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	c.ctx = ctx
	if c.LogLevel > 0 {
		elapsed := time.Since(begin)
		switch {
		case err != nil && c.LogLevel >= logger.Error:
			sql, rows := fc()
			if rows == -1 {
				c.Printf(c.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				c.Printf(c.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case elapsed > c.SlowThreshold && c.SlowThreshold != 0 && c.LogLevel >= logger.Warn:
			sql, rows := fc()
			slowLog := fmt.Sprintf("SLOW SQL >= %v", c.SlowThreshold)
			if rows == -1 {
				c.Printf(c.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				c.Printf(c.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case c.LogLevel >= logger.Info:
			sql, rows := fc()
			if rows == -1 {
				c.Printf(c.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				c.Printf(c.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		}
	}
}

func (c *customLogger) Printf(message string, data ...interface{}) {
	switch len(data) {
	case 0:
		c.log.WithContext(c.ctx).Info(message)
	case 1:
		c.log.WithContext(c.ctx).Infof("gorm,src:%v", data[0])
	case 2:
		c.log.WithContext(c.ctx).Infof("gorm,src:%v,duration:%v", data[0], data[1])
	case 3:
		c.log.WithContext(c.ctx).Infof("gorm,src:%v,duration:%v,row:%v", data[0], data[1], data[2])
	case 4:
		c.log.WithContext(c.ctx).Infof("gorm,src:%v,duration:%v,row:%v,sql:%v", data[0], data[1], data[2], data[3])
	}
	return
}
