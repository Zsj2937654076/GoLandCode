-- Sample data for the student management system
USE student_management;

-- Make sure we have the admin user
INSERT IGNORE INTO users (username, password, email, role) VALUES 
('admin', '$2a$10$1Nj6JPJm5gKxCj7OfL5.a.jZAIqMMGqjKnUse/qG7pCzJKHnFNGc.', 'admin@example.com', 'admin');

-- Sample classes
INSERT INTO classes (name, description) VALUES
('一年级一班', '一年级一班是一个充满活力的班级，共有30名学生。'),
('一年级二班', '一年级二班是一个团结友爱的班级，共有28名学生。'),
('二年级一班', '二年级一班是一个积极向上的班级，共有32名学生。'),
('二年级二班', '二年级二班是一个勤奋好学的班级，共有29名学生。'),
('三年级一班', '三年级一班是一个充满创造力的班级，共有31名学生。'),
('三年级二班', '三年级二班是一个团结互助的班级，共有30名学生。');

-- Sample students for Class 1-A
INSERT INTO students (student_id, name, class_id, email, phone, address) VALUES
('2024001', '张三', 1, 'zhangsan@example.com', '13800138001', '北京市海淀区'),
('2024002', '李四', 1, 'lisi@example.com', '13800138002', '北京市朝阳区'),
('2024003', '王五', 1, 'wangwu@example.com', '13800138003', '北京市西城区'),
('2024004', '赵六', 2, 'zhaoliu@example.com', '13800138004', '北京市东城区'),
('2024005', '钱七', 2, 'qianqi@example.com', '13800138005', '北京市丰台区'),
('2024006', '孙八', 2, 'sunba@example.com', '13800138006', '北京市石景山区'),
('2024007', '周九', 3, 'zhoujiu@example.com', '13800138007', '北京市通州区'),
('2024008', '吴十', 3, 'wushi@example.com', '13800138008', '北京市昌平区'),
('2024009', '郑十一', 3, 'zheng11@example.com', '13800138009', '北京市大兴区'),
('2024010', '王十二', 4, 'wang12@example.com', '13800138010', '北京市顺义区'),
('2024011', '李十三', 4, 'li13@example.com', '13800138011', '北京市房山区'),
('2024012', '张十四', 4, 'zhang14@example.com', '13800138012', '北京市门头沟区'),
('2024013', '刘十五', 5, 'liu15@example.com', '13800138013', '北京市怀柔区'),
('2024014', '陈十六', 5, 'chen16@example.com', '13800138014', '北京市平谷区'),
('2024015', '杨十七', 5, 'yang17@example.com', '13800138015', '北京市密云区'),
('2024016', '黄十八', 6, 'huang18@example.com', '13800138016', '北京市延庆区'),
('2024017', '赵十九', 6, 'zhao19@example.com', '13800138017', '北京市海淀区'),
('2024018', '周二十', 6, 'zhou20@example.com', '13800138018', '北京市朝阳区'); 