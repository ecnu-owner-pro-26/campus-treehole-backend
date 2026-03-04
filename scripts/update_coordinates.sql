-- 为现有地点数据添加坐标的脚本
-- 执行方式: sqlite3 data/campus_memory.db < scripts/update_coordinates.sql

-- 普陀校区地点坐标更新（华东师范大学普陀校区，中山北路校区）
-- 中心坐标约: 31.2304, 121.4245

UPDATE locations SET latitude = 31.2304, longitude = 121.4245 WHERE name = '图书馆' AND campus_id = 1;
UPDATE locations SET latitude = 31.2308, longitude = 121.4248 WHERE name = '三馆' AND campus_id = 1;
UPDATE locations SET latitude = 31.2310, longitude = 121.4250 WHERE name = '博学楼' AND campus_id = 1;
UPDATE locations SET latitude = 31.2306, longitude = 121.4242 WHERE name = '干训楼' AND campus_id = 1;
UPDATE locations SET latitude = 31.2312, longitude = 121.4252 WHERE name = '化学馆' AND campus_id = 1;
UPDATE locations SET latitude = 31.2302, longitude = 121.4240 WHERE name = '俊秀楼' AND campus_id = 1;
UPDATE locations SET latitude = 31.2314, longitude = 121.4254 WHERE name = '科教楼（脑所）' AND campus_id = 1;
UPDATE locations SET latitude = 31.2300, longitude = 121.4238 WHERE name = '数学馆' AND campus_id = 1;
UPDATE locations SET latitude = 31.2298, longitude = 121.4236 WHERE name = '体育楼' AND campus_id = 1;
UPDATE locations SET latitude = 31.2316, longitude = 121.4256 WHERE name = '微波所' AND campus_id = 1;
UPDATE locations SET latitude = 31.2296, longitude = 121.4234 WHERE name = '文附楼' AND campus_id = 1;
UPDATE locations SET latitude = 31.2294, longitude = 121.4232 WHERE name = '文史楼（群贤堂）' AND campus_id = 1;
UPDATE locations SET latitude = 31.2318, longitude = 121.4258 WHERE name = '物理楼' AND campus_id = 1;
UPDATE locations SET latitude = 31.2292, longitude = 121.4230 WHERE name = '小教楼' AND campus_id = 1;
UPDATE locations SET latitude = 31.2320, longitude = 121.4260 WHERE name = '专家楼' AND campus_id = 1;

-- 闵行校区地点坐标更新（华东师范大学闵行校区）
-- 中心坐标约: 31.0234, 121.4420

UPDATE locations SET latitude = 31.0234, longitude = 121.4420 WHERE name = '图书馆（桶楼）' AND campus_id = 2;
UPDATE locations SET latitude = 31.0236, longitude = 121.4422 WHERE name = '图书馆（裙楼）' AND campus_id = 2;
UPDATE locations SET latitude = 31.0238, longitude = 121.4424 WHERE name = '办公楼' AND campus_id = 2;
UPDATE locations SET latitude = 31.0240, longitude = 121.4426 WHERE name = '第一教学楼' AND campus_id = 2;
UPDATE locations SET latitude = 31.0242, longitude = 121.4428 WHERE name = '第二教学楼' AND campus_id = 2;
UPDATE locations SET latitude = 31.0244, longitude = 121.4430 WHERE name = '第三教学楼' AND campus_id = 2;
UPDATE locations SET latitude = 31.0246, longitude = 121.4432 WHERE name = '第四教学楼' AND campus_id = 2;

-- 为所有未设置坐标的地点设置默认坐标（按校区）
UPDATE locations SET latitude = 31.2304, longitude = 121.4245 WHERE campus_id = 1 AND (latitude = 0.0 OR latitude IS NULL);
UPDATE locations SET latitude = 31.0234, longitude = 121.4420 WHERE campus_id = 2 AND (latitude = 0.0 OR latitude IS NULL);
UPDATE locations SET latitude = 30.9, longitude = 121.9 WHERE campus_id = 3 AND (latitude = 0.0 OR latitude IS NULL);

SELECT '坐标更新完成！' AS status;
SELECT COUNT(*) || ' 个地点已更新坐标' AS info FROM locations WHERE latitude != 0.0;
