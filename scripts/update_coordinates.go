package main

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 连接数据库
	db, err := gorm.Open(sqlite.Open("data/campus_memory.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	// 更新普陀校区地点坐标（华东师范大学普陀校区）
	// 中心坐标约: 31.2304, 121.4245
	updates := []struct {
		name      string
		campusID  int64
		latitude  float64
		longitude float64
	}{
		{"图书馆", 1, 31.2304, 121.4245},
		{"三馆", 1, 31.2308, 121.4248},
		{"博学楼", 1, 31.2310, 121.4250},
		{"干训楼", 1, 31.2306, 121.4242},
		{"化学馆", 1, 31.2312, 121.4252},
		{"河西食堂", 1, 31.2303, 121.4243},
		{"河东食堂", 1, 31.2301, 121.4241},
		{"丽娃食堂", 1, 31.2299, 121.4239},
		{"大草坪", 1, 31.2327, 121.4267},
		{"体育馆", 1, 31.2331, 121.4271},
		
		// 闵行校区地点（华东师范大学闵行校区）
		// 中心坐标约: 31.0234, 121.4420
		{"图书馆（桶楼）", 2, 31.0234, 121.4420},
		{"图书馆（裙楼）", 2, 31.0236, 121.4422},
		{"办公楼", 2, 31.0238, 121.4424},
		{"第一教学楼", 2, 31.0240, 121.4426},
		{"第二教学楼", 2, 31.0242, 121.4428},
	}

	// 执行更新
	count := 0
	for _, u := range updates {
		result := db.Exec("UPDATE locations SET latitude = ?, longitude = ? WHERE name = ? AND campus_id = ?",
			u.latitude, u.longitude, u.name, u.campusID)
		if result.Error != nil {
			log.Printf("更新 %s 失败: %v", u.name, result.Error)
		} else if result.RowsAffected > 0 {
			count++
			fmt.Printf("✓ 更新 %s 坐标: (%.4f, %.4f)\n", u.name, u.latitude, u.longitude)
		}
	}

	// 为所有未设置坐标的地点设置默认坐标
	db.Exec("UPDATE locations SET latitude = 31.2304, longitude = 121.4245 WHERE campus_id = 1 AND (latitude = 0.0 OR latitude IS NULL)")
	db.Exec("UPDATE locations SET latitude = 31.0234, longitude = 121.4420 WHERE campus_id = 2 AND (latitude = 0.0 OR latitude IS NULL)")
	db.Exec("UPDATE locations SET latitude = 30.9, longitude = 121.9 WHERE campus_id = 3 AND (latitude = 0.0 OR latitude IS NULL)")

	fmt.Printf("\n✓ 迁移完成! 共更新 %d 个地点的坐标\n", count)
	
	// 统计
	var total int64
	db.Raw("SELECT COUNT(*) FROM locations WHERE latitude != 0.0").Scan(&total)
	fmt.Printf("✓ 数据库中共有 %d 个地点已设置坐标\n", total)
}
