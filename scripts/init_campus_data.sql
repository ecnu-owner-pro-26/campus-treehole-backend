-- 校区和地点初始化数据
DELETE FROM locations;
DELETE FROM campuses;

-- 自动生成ID
DELETE FROM sqlite_sequence WHERE name IN ('campuses', 'locations');

-- 插入校区数据
INSERT INTO campuses (name, is_active, sort_order, created_at, updated_at) VALUES
('普陀校区', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('闵行校区', 1, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('滴水湖校区', 1, 3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- 插入地点数据
-- 普陀校区地点
INSERT INTO locations (campus_id, name, category, icon, is_active, sort_order, memory_count, created_at, updated_at) VALUES
(1, '图书馆', 'library', '/static/icons/putuo-library.png', 1, 1, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, '理科大楼', 'academic', '/static/icons/science building.png', 1, 2, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, '文史楼', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, '文附楼', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, '文科大楼', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, '田家炳教书院', 'academic', '/static/icons/ tianjiabing building.png', 1, 4, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, '干训楼', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, '办公楼', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, '药学院', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, '物理楼', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, '计算机楼', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, '自然地理馆', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, '科学会堂', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, '河西食堂', 'dining', '/static/icons/marker-library.png', 1, 5, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, '河东食堂', 'dining', '/static/icons/marker-library.png', 1, 6, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, '丽娃食堂', 'dining', '/static/icons/marker-library.png', 1, 7, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, '体育馆', 'sports', '/static/icons/marker-library.png', 1, 8, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, '共青场', 'sports', '/static/icons/marker-library.png', 1, 9, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, '丽娃操场', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, '篮球场', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, '共青运动区', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, '大草坪', 'outdoor', '/static/icons/marker-library.png', 1, 10, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, '爱之坪', 'outdoor', '/static/icons/marker-library.png', 1, 11, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, '大学生活动中心', 'activity', '/static/icons/marker-library.png', 1, 12, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, '化学实验楼', 'academic','/static/icons/marker-library.png', 1, 13, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, '办公楼', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, '校医院', 'service', '/static/icons/marker-library.png', 1, 14, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, '校史馆', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, '第三宿舍', 'dormitory', '/static/icons/third dormitory.png', 1, 15, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, '第五宿舍', 'dormitory', '/static/icons/fifth dormitory.png',1, 16, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, '第七宿舍', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, '第八宿舍', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, '第十宿舍', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(1, '菜鸟驿站', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

-- 闵行校区地点
(2, '图书馆（桶楼）', 'library', '/static/icons/marker-library.png', 1, 17, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, '图书馆（裙楼）', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, '第一教学楼', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, '第二教学楼', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, '第三教学楼', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, '第四教学楼', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, '物理楼', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, '信息技术楼', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, '数学楼', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, '法商楼', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, '外语学院', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, '历史学系', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, '中文系', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, '化学与分子工程学院', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, '生命科学学院', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, '资环楼', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, '美术系', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, '实验A楼', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, '实验B楼', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, '法学院', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, '体育与健康学院', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, '游泳馆', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, '学生之家', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, '剑川路学生公寓', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, '冬月厅', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, '秋实阁', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, '秋林阁', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, '夏雨厅', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
(2, '冬日餐厅', 'academic', '/static/icons/wenshi building.png', 1, 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

-- 临港校区地点

SELECT '初始化完成: ' || COUNT(*) || '个校区, ' ||
       (SELECT COUNT(*) FROM locations) || '个地点'
FROM campuses;