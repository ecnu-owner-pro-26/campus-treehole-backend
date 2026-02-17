package infra

import (
	"campus-memory/infra/model"
	"fmt"
	"os"
)

// InitDatabase 实现初始化数据库
func InitDatabase(dbPath string) (*Database, error) {

	// 自动创建data目录（如果不存在）
	if err := os.MkdirAll("data", 0755); err != nil {
		return nil, fmt.Errorf("fail to create dir %s: %w", "data", err)
	}

	// 创建数据库连接
	db, err := NewDatabase(dbPath)
	if err != nil {
		return nil, fmt.Errorf("fail to connect database: %w", err)
	}

	// 自动迁移所有模型
	if err := migrateAllModels(db); err != nil {
		return nil, fmt.Errorf("fail to migrate models: %w", err)
	}

	// 返回Database实例
	return db, nil
}

// migrateAllModels 迁移所有7个模型
func migrateAllModels(db *Database) error {

	// 定义所有需要迁移的模型
	models := []interface{}{
		&model.CampusModel{},
		&model.CommentModel{},
		&model.ImageModel{},
		&model.LikeModel{},
		&model.LocationModel{},
		&model.MemoryModel{},
		&model.UserModel{},
	}

	// 执行迁移
	if err := db.AutoMigrate(models...); err != nil {
		return err
	}

	return nil
}
