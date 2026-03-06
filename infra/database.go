package infra

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database 数据库基础设施
type Database struct {
	DB     *gorm.DB
	DBPath string
}

// NewDatabase 创建数据库连接
func NewDatabase(dbPath string) (*Database, error) {

	// 配置GORM日志器
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // 标准输出
		logger.Config{
			SlowThreshold: time.Second,  // 慢SQL阈值
			LogLevel:      logger.Error, // 日志级别：Error级别显示重要信息
			Colorful:      true,         // 彩色输出
		},
	)

	// 配置GORM框架
	config := &gorm.Config{
		Logger:                                   gormLogger,
		DisableForeignKeyConstraintWhenMigrating: true, // SQLite迁移时禁用外键约束
	}

	// 打开数据库连接
	db, err := gorm.Open(sqlite.Open(dbPath), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最长存活时间

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return &Database{DB: db, DBPath: dbPath}, nil
}

// Close 关闭数据库连接
func (d *Database) Close() error {

	if d.DB == nil {
		return nil
	}

	// 获取底层连接
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB for closing: %w", err)
	}

	// 执行关闭数据库连接
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	return nil
}

// AutoMigrate 自动迁移数据库表
func (d *Database) AutoMigrate(models ...interface{}) error {

	if d.DB == nil {
		return nil
	}

	return d.DB.AutoMigrate(models...)
}
