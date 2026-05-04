-- 学习平台数据库初始化脚本
-- 数据库: study_exam
-- 字符集: utf8mb4
-- 排序规则: utf8mb4_unicode_ci

CREATE DATABASE IF NOT EXISTS study_exam DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_unicode_ci;

USE study_exam;

-- 管理员表
CREATE TABLE IF NOT EXISTS admins (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    password VARCHAR(255) NOT NULL COMMENT '密码',
    nickname VARCHAR(50) COMMENT '昵称',
    avatar VARCHAR(255) COMMENT '头像',
    role VARCHAR(20) DEFAULT 'admin' COMMENT '角色',
    status TINYINT DEFAULT 1 COMMENT '状态 1:启用 0:禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理员表';

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    password VARCHAR(255) NOT NULL COMMENT '密码',
    nickname VARCHAR(50) COMMENT '昵称',
    avatar VARCHAR(255) COMMENT '头像',
    email VARCHAR(100) COMMENT '邮箱',
    phone VARCHAR(20) COMMENT '手机号',
    gender TINYINT DEFAULT 0 COMMENT '性别 0:未知 1:男 2:女',
    birthday DATE COMMENT '生日',
    status TINYINT DEFAULT 1 COMMENT '状态 1:启用 0:禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_username (username),
    INDEX idx_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 课程类型表
CREATE TABLE IF NOT EXISTS course_types (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL COMMENT '类型名称',
    description TEXT COMMENT '类型描述',
    sort INT DEFAULT 0 COMMENT '排序',
    status TINYINT DEFAULT 1 COMMENT '状态 1:启用 0:禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='课程类型表';

-- 课程表
CREATE TABLE IF NOT EXISTS courses (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(200) NOT NULL COMMENT '课程标题',
    type_id BIGINT UNSIGNED NOT NULL COMMENT '课程类型ID',
    description TEXT COMMENT '课程描述',
    cover VARCHAR(255) COMMENT '封面图片',
    author VARCHAR(100) COMMENT '作者/讲师',
    price DECIMAL(10,2) DEFAULT 0.00 COMMENT '价格',
    is_free TINYINT DEFAULT 1 COMMENT '是否免费 1:免费 0:付费',
    view_count INT DEFAULT 0 COMMENT '浏览次数',
    collect_count INT DEFAULT 0 COMMENT '收藏次数',
    status TINYINT DEFAULT 1 COMMENT '状态 1:上架 0:下架',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_type_id (type_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='课程表';

-- 章节表
CREATE TABLE IF NOT EXISTS chapters (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    course_id BIGINT UNSIGNED NOT NULL COMMENT '课程ID',
    title VARCHAR(200) NOT NULL COMMENT '章节标题',
    sort INT DEFAULT 0 COMMENT '排序',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_course_id (course_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='章节表';

-- 小节/视频表
CREATE TABLE IF NOT EXISTS sections (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    chapter_id BIGINT UNSIGNED NOT NULL COMMENT '章节ID',
    course_id BIGINT UNSIGNED NOT NULL COMMENT '课程ID',
    title VARCHAR(200) NOT NULL COMMENT '小节标题',
    content_type VARCHAR(20) DEFAULT 'video' COMMENT '内容类型 video:视频 book:书籍',
    video_url VARCHAR(255) COMMENT '视频地址',
    book_content TEXT COMMENT '书籍内容',
    duration INT DEFAULT 0 COMMENT '视频时长(秒)',
    sort INT DEFAULT 0 COMMENT '排序',
    view_count INT DEFAULT 0 COMMENT '浏览次数',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_chapter_id (chapter_id),
    INDEX idx_course_id (course_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='小节/视频表';

-- 收藏表
CREATE TABLE IF NOT EXISTS favorites (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    favorite_type VARCHAR(20) NOT NULL COMMENT '收藏类型 course:课程 section:视频',
    target_id BIGINT UNSIGNED NOT NULL COMMENT '目标ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_target (user_id, favorite_type, target_id),
    INDEX idx_user_id (user_id),
    INDEX idx_target (favorite_type, target_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='收藏表';

-- 评论表
CREATE TABLE IF NOT EXISTS comments (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    section_id BIGINT UNSIGNED NOT NULL COMMENT '视频/小节ID',
    course_id BIGINT UNSIGNED NOT NULL COMMENT '课程ID',
    content TEXT NOT NULL COMMENT '评论内容',
    parent_id BIGINT UNSIGNED DEFAULT 0 COMMENT '父评论ID',
    status TINYINT DEFAULT 1 COMMENT '状态 1:显示 0:隐藏',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_section_id (section_id),
    INDEX idx_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='评论表';

-- 用户行为表（用于协同过滤推荐）
CREATE TABLE IF NOT EXISTS user_behaviors (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    course_id BIGINT UNSIGNED NOT NULL COMMENT '课程ID',
    behavior_type VARCHAR(20) NOT NULL COMMENT '行为类型 view:浏览 collect:收藏 comment:评论 study:学习',
    rating INT DEFAULT 1 COMMENT '评分/权重 (1-5)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_course_id (course_id),
    INDEX idx_behavior (user_id, course_id, behavior_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户行为表';

-- 题库表
CREATE TABLE IF NOT EXISTS question_banks (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL COMMENT '题库名称',
    description TEXT COMMENT '题库描述',
    question_count INT DEFAULT 0 COMMENT '题目数量',
    status TINYINT DEFAULT 1 COMMENT '状态 1:启用 0:禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='题库表';

-- 试题表
CREATE TABLE IF NOT EXISTS questions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    bank_id BIGINT UNSIGNED NOT NULL COMMENT '题库ID',
    question_type VARCHAR(20) NOT NULL COMMENT '题目类型 single:单选 multiple:多选 judge:判断 fill:填空 essay:问答',
    content TEXT NOT NULL COMMENT '题目内容',
    options JSON COMMENT '选项 (JSON格式)',
    answer TEXT NOT NULL COMMENT '正确答案',
    analysis TEXT COMMENT '解析',
    score INT DEFAULT 2 COMMENT '分数',
    difficulty TINYINT DEFAULT 2 COMMENT '难度 1:简单 2:中等 3:困难',
    status TINYINT DEFAULT 1 COMMENT '状态 1:启用 0:禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_bank_id (bank_id),
    INDEX idx_question_type (question_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='试题表';

-- 试卷表
CREATE TABLE IF NOT EXISTS papers (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(200) NOT NULL COMMENT '试卷标题',
    description TEXT COMMENT '试卷描述',
    total_score INT DEFAULT 100 COMMENT '总分',
    pass_score INT DEFAULT 60 COMMENT '及格分',
    duration INT DEFAULT 60 COMMENT '考试时长(分钟)',
    question_count INT DEFAULT 0 COMMENT '题目数量',
    status TINYINT DEFAULT 1 COMMENT '状态 1:启用 0:禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='试卷表';

-- 试卷试题关联表
CREATE TABLE IF NOT EXISTS paper_questions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    paper_id BIGINT UNSIGNED NOT NULL COMMENT '试卷ID',
    question_id BIGINT UNSIGNED NOT NULL COMMENT '试题ID',
    sort INT DEFAULT 0 COMMENT '排序',
    score INT DEFAULT 2 COMMENT '该题分数',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_paper_id (paper_id),
    INDEX idx_question_id (question_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='试卷试题关联表';

-- 考试表
CREATE TABLE IF NOT EXISTS exams (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(200) NOT NULL COMMENT '考试标题',
    paper_id BIGINT UNSIGNED NOT NULL COMMENT '试卷ID',
    description TEXT COMMENT '考试描述',
    start_time TIMESTAMP NULL COMMENT '开始时间',
    end_time TIMESTAMP NULL COMMENT '结束时间',
    duration INT DEFAULT 60 COMMENT '考试时长(分钟)',
    total_score INT DEFAULT 100 COMMENT '总分',
    pass_score INT DEFAULT 60 COMMENT '及格分',
    status TINYINT DEFAULT 1 COMMENT '状态 1:进行中 0:已结束',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_paper_id (paper_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='考试表';

-- 考试记录表
CREATE TABLE IF NOT EXISTS exam_records (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    exam_id BIGINT UNSIGNED NOT NULL COMMENT '考试ID',
    paper_id BIGINT UNSIGNED NOT NULL COMMENT '试卷ID',
    score INT DEFAULT 0 COMMENT '得分',
    total_score INT DEFAULT 100 COMMENT '总分',
    is_pass TINYINT DEFAULT 0 COMMENT '是否及格 1:及格 0:不及格',
    start_time TIMESTAMP NULL COMMENT '开始答题时间',
    submit_time TIMESTAMP NULL COMMENT '提交时间',
    time_spent INT DEFAULT 0 COMMENT '用时(秒)',
    status VARCHAR(20) DEFAULT 'pending' COMMENT '状态 pending:进行中 submitted:已提交 corrected:已批改',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_exam_id (exam_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='考试记录表';

-- 用户答案表
CREATE TABLE IF NOT EXISTS user_answers (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    record_id BIGINT UNSIGNED NOT NULL COMMENT '考试记录ID',
    question_id BIGINT UNSIGNED NOT NULL COMMENT '试题ID',
    user_answer TEXT COMMENT '用户答案',
    is_correct TINYINT DEFAULT 0 COMMENT '是否正确 1:正确 0:错误',
    score INT DEFAULT 0 COMMENT '得分',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_record_id (record_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户答案表';

-- 论坛帖子表
CREATE TABLE IF NOT EXISTS forums (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    title VARCHAR(200) NOT NULL COMMENT '帖子标题',
    content TEXT NOT NULL COMMENT '帖子内容',
    view_count INT DEFAULT 0 COMMENT '浏览次数',
    reply_count INT DEFAULT 0 COMMENT '回复次数',
    status TINYINT DEFAULT 1 COMMENT '状态 1:显示 0:隐藏',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='论坛帖子表';

-- 论坛评论/回复表
CREATE TABLE IF NOT EXISTS forum_replies (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    forum_id BIGINT UNSIGNED NOT NULL COMMENT '帖子ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    content TEXT NOT NULL COMMENT '回复内容',
    parent_id BIGINT UNSIGNED DEFAULT 0 COMMENT '父回复ID',
    status TINYINT DEFAULT 1 COMMENT '状态 1:显示 0:隐藏',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_forum_id (forum_id),
    INDEX idx_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='论坛评论/回复表';

-- 公告表
CREATE TABLE IF NOT EXISTS announcements (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(200) NOT NULL COMMENT '公告标题',
    content TEXT NOT NULL COMMENT '公告内容',
    author VARCHAR(50) COMMENT '发布人',
    is_top TINYINT DEFAULT 0 COMMENT '是否置顶 1:置顶 0:否',
    sort INT DEFAULT 0 COMMENT '排序',
    status TINYINT DEFAULT 1 COMMENT '状态 1:显示 0:隐藏',
    view_count INT DEFAULT 0 COMMENT '浏览次数',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='公告表';

-- 轮播图表
CREATE TABLE IF NOT EXISTS banners (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(100) COMMENT '标题',
    image VARCHAR(255) NOT NULL COMMENT '图片地址',
    link VARCHAR(255) COMMENT '跳转链接',
    sort INT DEFAULT 0 COMMENT '排序',
    status TINYINT DEFAULT 1 COMMENT '状态 1:显示 0:隐藏',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='轮播图表';

-- 初始化管理员账号 (密码: admin123)
INSERT INTO admins (username, password, nickname, role, status) VALUES 
('admin', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt6Z5E', '超级管理员', 'super_admin', 1);

-- 初始化一些课程类型
INSERT INTO course_types (name, description, sort, status) VALUES 
('编程开发', '各类编程语言和框架学习', 1, 1),
('前端开发', 'HTML、CSS、JavaScript及前端框架', 2, 1),
('后端开发', 'Java、Python、Go等后端技术', 3, 1),
('移动开发', 'Android、iOS、Flutter等移动开发', 4, 1),
('人工智能', '机器学习、深度学习、AI应用', 5, 1),
('数据库', 'MySQL、Oracle、MongoDB等数据库技术', 6, 1);

-- 初始化一些课程
INSERT INTO courses (title, type_id, description, cover, author, price, is_free, view_count, collect_count, status) VALUES 
('Go语言从入门到精通', 1, '本课程全面讲解Go语言的基础语法、并发编程、Web开发等核心技术', 'https://example.com/go.jpg', '张老师', 0.00, 1, 1560, 320, 1),
('React实战开发', 2, '深入学习React 18的核心概念、Hooks、状态管理等技术', 'https://example.com/react.jpg', '李老师', 99.00, 0, 2340, 560, 1),
('Python数据分析', 3, '使用Python进行数据清洗、分析和可视化', 'https://example.com/python.jpg', '王老师', 0.00, 1, 1890, 420, 1),
('Vue3企业级开发', 2, 'Vue3 + TypeScript + Element Plus 实战开发', 'https://example.com/vue.jpg', '赵老师', 129.00, 0, 3200, 780, 1),
('机器学习入门', 5, '从数学基础到实战应用，全面掌握机器学习核心算法', 'https://example.com/ml.jpg', '刘老师', 199.00, 0, 1200, 280, 1);

-- 初始化章节
INSERT INTO chapters (course_id, title, sort) VALUES 
(1, 'Go语言基础', 1),
(1, 'Go语言进阶', 2),
(1, 'Go并发编程', 3),
(2, 'React基础', 1),
(2, 'React Hooks', 2),
(2, 'React状态管理', 3);

-- 初始化小节
INSERT INTO sections (chapter_id, course_id, title, content_type, video_url, duration, sort, view_count) VALUES 
(1, 1, 'Go语言简介与环境搭建', 'video', 'https://example.com/go1.mp4', 1800, 1, 890),
(1, 1, '变量与数据类型', 'video', 'https://example.com/go2.mp4', 2400, 2, 750),
(1, 1, '函数与方法', 'video', 'https://example.com/go3.mp4', 3000, 3, 680),
(2, 1, '结构体与接口', 'video', 'https://example.com/go4.mp4', 3600, 1, 560),
(2, 1, '错误处理', 'video', 'https://example.com/go5.mp4', 1800, 2, 520),
(4, 2, 'React简介与环境搭建', 'video', 'https://example.com/react1.mp4', 1800, 1, 1200),
(4, 2, 'JSX语法', 'video', 'https://example.com/react2.mp4', 2400, 2, 1100);

-- 初始化题库
INSERT INTO question_banks (name, description, question_count, status) VALUES 
('Go语言基础题库', '包含Go语言基础知识的各类题目', 50, 1),
('React开发题库', 'React框架相关的考试题目', 40, 1),
('Python基础题库', 'Python编程语言基础题目', 60, 1);

-- 初始化试题
INSERT INTO questions (bank_id, question_type, content, options, answer, analysis, score, difficulty, status) VALUES 
(1, 'single', 'Go语言中，哪个关键字用于定义常量？', '{"A": "var", "B": "const", "C": "let", "D": "static"}', 'B', 'Go语言使用const关键字定义常量，var用于定义变量。', 2, 1, 1),
(1, 'single', 'Go语言中，main函数所在的包名是？', '{"A": "main", "B": "app", "C": "root", "D": "index"}', 'A', 'Go语言中，main函数必须在main包中。', 2, 1, 1),
(1, 'multiple', '以下哪些是Go语言的数据类型？', '{"A": "int", "B": "string", "C": "boolean", "D": "float64"}', 'A,B,D', 'Go语言的数据类型包括int、string、float64等，布尔类型是bool不是boolean。', 4, 1, 1),
(2, 'single', 'React中用于管理组件内部状态的Hook是？', '{"A": "useEffect", "B": "useState", "C": "useContext", "D": "useRef"}', 'B', 'useState是React中用于管理组件状态的Hook。', 2, 1, 1),
(2, 'single', 'React组件的渲染方法返回什么？', '{"A": "DOM元素", "B": "React元素", "C": "字符串", "D": "对象"}', 'B', 'React组件的render方法或函数组件返回React元素。', 2, 2, 1);

-- 初始化试卷
INSERT INTO papers (title, description, total_score, pass_score, duration, question_count, status) VALUES 
('Go语言基础测试', '测试Go语言基础知识掌握情况', 100, 60, 60, 10, 1),
('React基础测试', '测试React框架基础知识', 100, 60, 45, 8, 1);

-- 初始化试卷试题关联
INSERT INTO paper_questions (paper_id, question_id, sort, score) VALUES 
(1, 1, 1, 10),
(1, 2, 2, 10),
(1, 3, 3, 20),
(2, 4, 1, 15),
(2, 5, 2, 15);

-- 初始化公告
INSERT INTO announcements (title, content, author, is_top, sort, status, view_count) VALUES 
('欢迎使用学习平台', '欢迎使用我们的在线学习平台！本平台提供丰富的课程资源和在线考试功能。', '管理员', 1, 1, 1, 2340),
('新功能上线：课程推荐', '我们的平台现已上线协同过滤推荐功能，根据您的学习行为为您推荐合适的课程。', '管理员', 0, 2, 1, 1560),
('在线考试功能升级', '在线考试功能已全面升级，支持更多题型和更完善的考试体验。', '管理员', 0, 3, 1, 1200);

-- 初始化轮播图
INSERT INTO banners (title, image, link, sort, status) VALUES 
('Go语言课程', 'https://example.com/banner1.jpg', '/course/1', 1, 1),
('React实战', 'https://example.com/banner2.jpg', '/course/2', 2, 1),
('Python数据分析', 'https://example.com/banner3.jpg', '/course/3', 3, 1);
