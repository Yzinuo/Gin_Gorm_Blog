/*
 Navicat Premium Data Transfer

 Source Server         : 172.18.45.12
 Source Server Type    : MySQL
 Source Server Version : 80031
 Source Host           : 172.18.45.12:3306
 Source Schema         : ginblog

 Target Server Type    : MySQL
 Target Server Version : 80031
 File Encoding         : 65001

 Date: 29/12/2023 23:17:19
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

CREATE DATABASE IF NOT EXISTS`gvb` DEFAULT CHARACTER SET utf8mb4;
USE `gvb`;

-- ----------------------------
-- Table structure for article
-- ----------------------------
DROP TABLE IF EXISTS `article`;
CREATE TABLE `article`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `title` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `desc` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `img` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `type` tinyint NULL DEFAULT NULL COMMENT '类型(1-原创 2-转载 3-翻译)',
  `status` tinyint NULL DEFAULT NULL COMMENT '状态(1-公开 2-私密)',
  `is_top` tinyint(1) NULL DEFAULT NULL,
  `is_delete` tinyint(1) NULL DEFAULT NULL,
  `original_url` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `category_id` bigint NULL DEFAULT NULL,
  `user_id` bigint NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of article
-- ----------------------------
INSERT INTO `article` VALUES (2, '2023-12-27 22:47:47.513', '2023-12-27 22:48:58.872', '学习有捷径', '', '学习有捷径。学习的捷径之一就是多看看别人是怎么理解这些知识的。\n\n举两个例子。\n\n如果你喜欢《水浒》，千万不要只读原著当故事看，一定要读一读各代名家的批注和点评，看他们是如何理解的。之前学 C# 时，看《CLR via C#》晦涩难懂，但是我又通过看《你必须知道的.net》而更加了解了。因为后者就是中国一个 80 后写的，我通过他对 C# 的了解，作为桥梁和阶梯，再去窥探比较高达上的书籍和知识。\n\n最后，真诚的希望你能借助别人的力量来提高自己。我也一直在这样要求我自己。\n\n$$\n1/2 + 3/4 + 5/6 + 7^{99} = 999\n$$', 'https://gvbresource.oss-cn-hongkong.aliyuncs.com/%E9%BB%98%E8%AE%A4%E5%B0%81%E9%9D%A2.jpg', 1, 1, 0, 0, '', 4, 1);


-- ----------------------------
-- Table structure for article_tag
-- ----------------------------
DROP TABLE IF EXISTS `article_tag`;
CREATE TABLE `article_tag`  (
  `tag_id` bigint NOT NULL,
  `article_id` bigint NOT NULL,
  PRIMARY KEY (`tag_id`, `article_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of article_tag
-- ----------------------------
INSERT INTO `article_tag` VALUES (1, 1);
INSERT INTO `article_tag` VALUES (1, 3);
INSERT INTO `article_tag` VALUES (2, 1);
INSERT INTO `article_tag` VALUES (2, 3);
INSERT INTO `article_tag` VALUES (3, 2);

-- ----------------------------
-- Table structure for category
-- ----------------------------
DROP TABLE IF EXISTS `category`;
CREATE TABLE `category`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `name`(`name` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of category
-- ----------------------------
INSERT INTO `category` VALUES (1, '2023-12-27 22:45:09.369', '2023-12-27 22:45:09.369', '后端');
INSERT INTO `category` VALUES (2, '2023-12-27 22:45:15.006', '2023-12-27 22:45:15.006', '前端');
INSERT INTO `category` VALUES (3, '2023-12-27 22:46:36.057', '2023-12-27 22:46:36.057', '项目');
INSERT INTO `category` VALUES (4, '2023-12-27 22:47:47.501', '2023-12-27 22:47:47.501', '学习');

-- ----------------------------
-- Table structure for comment
-- ----------------------------
DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `user_id` bigint NULL DEFAULT NULL,
  `reply_user_id` bigint NULL DEFAULT NULL,
  `topic_id` bigint NULL DEFAULT NULL,
  `parent_id` bigint NULL DEFAULT NULL,
  `content` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `type` tinyint(1) NOT NULL COMMENT '评论类型(1.文章 2.友链 3.说说)',
  `is_review` tinyint(1) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of comment
-- ----------------------------

-- ----------------------------
-- Table structure for config
-- ----------------------------
DROP TABLE IF EXISTS `config`;
CREATE TABLE `config`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `config` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `key` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `value` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `desc` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `key`(`key` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 17 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of config
-- ----------------------------
INSERT INTO `config` VALUES (1, '2023-12-27 22:40:22.813', '2023-12-27 23:01:35.029', '', 'website_avatar', 'https://gvbresource.oss-cn-hongkong.aliyuncs.com/%E7%BD%91%E7%AB%99%E5%A4%B4%E5%83%8F.jpg', '网站头像');
INSERT INTO `config` VALUES (2, '2023-12-27 22:40:22.813', '2023-12-27 23:01:35.033', '', 'website_name', '子诺个人博客', '网站名称');
INSERT INTO `config` VALUES (3, '2023-12-27 22:40:22.813', '2023-12-27 23:01:35.027', '', 'website_author', '子诺', '网站作者');
INSERT INTO `config` VALUES (4, '2023-12-27 22:40:22.813', '2023-12-27 23:01:35.023', '', 'website_intro', '永远强过昨的的自己', '网站介绍');
INSERT INTO `config` VALUES (5, '2023-12-27 22:40:22.813', '2023-12-27 23:01:35.038', '', 'website_notice', '我的处女作,比较简陋,just for fun', '网站公告');
INSERT INTO `config` VALUES (6, '2023-12-27 22:40:22.813', '2023-12-27 23:01:35.031', '', 'website_createtime', '2023-12-27 22:40:22', '网站创建日期');
INSERT INTO `config` VALUES (7, '2023-12-27 22:40:22.813', '2023-12-27 23:01:35.011', '', 'website_record', '粤ICP备2021032312号', '网站备案号');
INSERT INTO `config` VALUES (8, '2023-12-27 22:40:22.813', '2023-12-27 23:01:35.008', '', 'qq', '123456789', 'QQ');
INSERT INTO `config` VALUES (9, '2023-12-27 22:40:22.813', '2023-12-27 23:01:35.015', '', 'github', 'https://github.com/Yzinuo', 'github');
INSERT INTO `config` VALUES (10, '2023-12-27 22:40:22.813', '2023-12-27 23:01:35.025', '', 'gitee', 'https://github.com/Yzinuo', 'gitee');
INSERT INTO `config` VALUES (11, '2023-12-27 22:40:22.813', '2023-12-27 23:01:35.019', '', 'tourist_avatar', 'https://cdn.hahacode.cn/config/tourist_avatar.png', '默认游客头像');
INSERT INTO `config` VALUES (12, '2023-12-27 22:40:22.813', '2023-12-27 23:01:35.021', '', 'user_avatar', 'https://cdn.hahacode.cn/config/user_avatar.png', '默认用户头像');
INSERT INTO `config` VALUES (13, '2023-12-27 22:40:22.813', '2023-12-27 23:01:35.013', '', 'article_cover', 'https://gvbresource.oss-cn-hongkong.aliyuncs.com/%E9%BB%98%E8%AE%A4%E5%B0%81%E9%9D%A2.jpg', '默认文章封面');
INSERT INTO `config` VALUES (14, '2023-12-27 22:40:22.813', '2023-12-27 23:01:35.039', '', 'is_comment_review', 'true', '评论默认审核');
INSERT INTO `config` VALUES (15, '2023-12-27 22:40:22.813', '2023-12-27 23:01:35.017', '', 'is_message_review', 'true', '留言默认审核');
INSERT INTO `config` VALUES (16, '2023-12-27 22:59:20.110', '2023-12-27 23:01:35.035', '', 'about', '```javascript\nconsole.log(\"Hello World\")\n```\n\n我就是我，不一样的烟火！', '');

-- ----------------------------
-- Table structure for friend_link
-- ----------------------------
DROP TABLE IF EXISTS `friend_link`;
CREATE TABLE `friend_link`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `intro` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of friend_link
-- ----------------------------

-- ----------------------------
-- Table structure for menu
-- ----------------------------
DROP TABLE IF EXISTS `menu`;
CREATE TABLE `menu`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `parent_id` bigint NULL DEFAULT NULL,
  `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `path` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `component` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `icon` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `order_num` tinyint NULL DEFAULT NULL,
  `redirect` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `catalogue` tinyint(1) NULL DEFAULT NULL,
  `hidden` tinyint(1) NULL DEFAULT NULL,
  `keep_alive` tinyint(1) NULL DEFAULT NULL,
  `external` tinyint(1) NULL DEFAULT NULL,
  `external_link` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 49 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of menu
-- ----------------------------
INSERT INTO `menu` VALUES (2, '2022-10-31 09:41:03.000', '2023-12-27 23:26:43.807', 0, '文章管理', '/article', 'Layout', 'ic:twotone-article', 1, '/article/list', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (3, '2022-10-31 09:41:03.000', '2023-12-24 23:33:34.013', 0, '消息管理', '/message', 'Layout', 'ic:twotone-email', 2, '/message/comment	', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (4, '2022-10-31 09:41:03.000', '2023-12-24 23:32:35.177', 0, '用户管理', '/user', 'Layout', 'ph:user-list-bold', 4, '/user/list', 0, 0, 0, 0, NULL);
INSERT INTO `menu` VALUES (5, '2022-10-31 09:41:03.000', '2023-12-24 23:32:34.788', 0, '系统管理', '/setting', 'Layout', 'ion:md-settings', 5, '/setting/website', 0, 0, 0, 0, NULL);
INSERT INTO `menu` VALUES (6, '2022-10-31 09:41:03.000', '2023-12-24 23:22:29.519', 2, '发布文章', 'write', '/article/write', 'icon-park-outline:write', 1, '', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (8, '2022-10-31 09:41:03.000', '2023-12-21 20:58:29.873', 2, '文章列表', 'list', '/article/list', 'material-symbols:format-list-bulleted', 2, '', 0, 0, 0, 0, NULL);
INSERT INTO `menu` VALUES (9, '2022-10-31 09:41:03.000', '2022-11-01 01:18:30.931', 2, '分类管理', 'category', '/article/category', 'tabler:category', 3, '', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (10, '2022-10-31 09:41:03.000', '2022-11-01 01:18:35.502', 2, '标签管理', 'tag', '/article/tag', 'tabler:tag', 4, '', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (16, '2022-10-31 09:41:03.000', '2022-11-01 10:11:23.195', 0, '权限管理', '/auth', 'Layout', 'cib:adguard', 3, '/auth/menu', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (17, '2022-10-31 09:41:03.000', NULL, 16, '菜单管理', 'menu', '/auth/menu', 'ic:twotone-menu-book', 1, NULL, 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (23, '2022-10-31 09:41:03.000', NULL, 16, '接口管理', 'resource', '/auth/resource', 'mdi:api', 2, NULL, 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (24, '2022-10-31 09:41:03.000', '2022-10-31 10:09:18.913', 16, '角色管理', 'role', '/auth/role', 'carbon:user-role', 3, NULL, 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (25, '2022-10-31 10:11:09.232', '2022-11-01 01:29:48.520', 3, '评论管理', 'comment', '/message/comment', 'ic:twotone-comment', 1, '', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (26, '2022-10-31 10:12:01.546', '2022-11-01 01:29:54.130', 3, '留言管理', 'leave-msg', '/message/leave-msg', 'ic:twotone-message', 2, '', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (27, '2022-10-31 10:54:03.201', '2022-11-01 01:30:06.901', 4, '用户列表', 'list', '/user/list', 'mdi:account', 1, '', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (28, '2022-10-31 10:54:34.167', '2022-11-01 01:30:13.400', 4, '在线用户', 'online', '/user/online', 'ic:outline-online-prediction', 2, '', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (29, '2022-10-31 10:59:33.255', '2022-11-01 01:30:20.688', 5, '网站管理', 'website', '/setting/website', 'el:website', 1, '', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (30, '2022-10-31 11:00:09.997', '2022-11-01 01:30:24.097', 5, '页面管理', 'page', '/setting/page', 'iconoir:journal-page', 2, '', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (31, '2022-10-31 11:00:33.543', '2022-11-01 01:30:28.497', 5, '友链管理', 'link', '/setting/link', 'mdi:telegram', 3, '', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (32, '2022-10-31 11:01:00.444', '2022-11-01 01:30:33.186', 5, '关于我', 'about', '/setting/about', 'cib:about-me', 4, '', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (33, '2022-11-01 01:43:10.142', '2023-12-27 23:26:41.553', 0, '首页', '/home', '/home', 'ic:sharp-home', 0, '', 1, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (34, '2022-11-01 09:54:36.252', '2022-11-01 10:07:00.254', 2, '修改文章', 'write/:id', '/article/write', 'icon-park-outline:write', 1, '', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (36, '2022-11-04 15:50:45.993', '2023-12-24 23:32:33.538', 0, '日志管理', '/log', 'Layout', 'material-symbols:receipt-long-outline-rounded', 6, '/log/operation', 0, 0, 0, 0, NULL);
INSERT INTO `menu` VALUES (37, '2022-11-04 15:53:00.251', '2023-12-24 23:15:22.034', 36, '操作日志', 'operation', '/log/operation', 'mdi:book-open-page-variant-outline', 1, '', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (38, '2022-11-04 16:02:42.306', '2022-11-04 16:05:35.761', 36, '登录日志', 'login', '/log/login', 'material-symbols:login', 2, '', 0, 0, 1, 0, NULL);
INSERT INTO `menu` VALUES (39, '2022-12-07 20:47:08.349', '2023-12-24 23:33:35.701', 0, '个人中心', '/profile', '/profile', 'mdi:account', 7, '', 1, 0, 0, 0, NULL);
INSERT INTO `menu` VALUES (47, '2023-12-24 20:26:14.173', '2023-12-24 23:33:36.247', 0, '测试一级菜单', '/testone', 'Layout', '', 88, '', 0, 0, 0, 1, NULL);
INSERT INTO `menu` VALUES (48, '2023-12-24 23:26:19.441', '2023-12-24 23:26:27.704', 0, '测试外链', 'https://www.baidu.com', 'Layout', 'mdi-fan-speed-3', 66, '', 1, 0, 0, 1, '');

-- ----------------------------
-- Table structure for message
-- ----------------------------
DROP TABLE IF EXISTS `message`;
CREATE TABLE `message`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '昵称',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '头像地址',
  `content` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '留言内容',
  `ip_address` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'IP 地址',
  `ip_source` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'IP 来源',
  `speed` tinyint(1) NULL DEFAULT NULL COMMENT '弹幕速度',
  `is_review` tinyint(1) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of message
-- ----------------------------

-- ----------------------------
-- Table structure for operation_log
-- ----------------------------
DROP TABLE IF EXISTS `operation_log`;
CREATE TABLE `operation_log`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `opt_module` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '操作模块',
  `opt_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '操作类型',
  `opt_method` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '操作方法',
  `opt_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '操作URL',
  `opt_desc` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '操作描述',
  `request_param` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '请求参数',
  `request_method` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '请求方法',
  `response_data` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '响应数据',
  `user_id` bigint NULL DEFAULT NULL COMMENT '用户ID',
  `nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '用户昵称',
  `ip_address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '操作IP',
  `ip_source` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '操作地址',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of operation_log
-- ----------------------------

-- ----------------------------
-- Table structure for page
-- ----------------------------
DROP TABLE IF EXISTS `page`;
CREATE TABLE `page`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `label` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `cover` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `name`(`name` ASC) USING BTREE,
  UNIQUE INDEX `label`(`label` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 20 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of page
-- ----------------------------
INSERT INTO `page` VALUES (1, '2022-12-08 13:09:58.500', '2023-12-28 16:31:43.682', '首页', 'home', 'https://cdn.hahacode.cn/page/home.jpg');
INSERT INTO `page` VALUES (2, '2022-12-08 13:51:49.474', '2023-12-28 14:55:58.704', '归档', 'archive', 'https://cdn.hahacode.cn/page/tag.png');
INSERT INTO `page` VALUES (3, '2022-12-08 13:52:18.084', '2023-12-28 16:31:30.137', '分类', 'category', 'https://cdn.hahacode.cn/page/category.png');
INSERT INTO `page` VALUES (4, '2022-12-08 13:52:31.364', '2023-12-28 14:55:45.058', '标签', 'tag', 'https://cdn.hahacode.cn/page/tag.png');
INSERT INTO `page` VALUES (5, '2022-12-08 13:52:52.389', '2023-12-28 15:02:21.859', '友链', 'link', 'https://cdn.hahacode.cn/page/link.jpg');
INSERT INTO `page` VALUES (6, '2022-12-08 13:53:04.159', '2023-12-28 16:30:03.928', '关于', 'about', 'https://cdn.hahacode.cn/page/about.jpg');
INSERT INTO `page` VALUES (7, '2022-12-08 13:53:17.707', '2023-12-28 16:27:13.418', '留言', 'message', 'https://gvbresource.oss-cn-hongkong.aliyuncs.com/111.jpg');
INSERT INTO `page` VALUES (8, '2022-12-08 13:53:30.187', '2023-12-28 14:55:25.724', '个人中心', 'user', 'https://gvbresource.oss-cn-hongkong.aliyuncs.com/222.jpg');
INSERT INTO `page` VALUES (9, '2022-12-16 23:54:52.650', '2023-12-28 14:54:42.341', '相册', 'album', 'https://cdn.hahacode.cn/page/album.png');
INSERT INTO `page` VALUES (10, '2022-12-16 23:55:36.059', '2023-12-28 14:55:09.345', '错误页面', '404', 'https://cdn.hahacode.cn/page/404.jpg');
INSERT INTO `page` VALUES (11, '2022-12-16 23:56:17.917', '2023-12-28 16:33:16.644', '文章列表', 'article_list', 'https://cdn.hahacode.cn/page/article_list.jpg');

-- ----------------------------
-- Table structure for resource
-- ----------------------------
DROP TABLE IF EXISTS `resource`;
CREATE TABLE `resource`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `parent_id` bigint NULL DEFAULT NULL,
  `url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `method` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `anonymous` tinyint(1) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 117 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of resource
-- ----------------------------
INSERT INTO `resource` VALUES (3, '2022-10-20 22:42:00.664', '2022-10-20 22:42:00.664', 0, '', '', '文章模块', 0);
INSERT INTO `resource` VALUES (6, '2022-10-20 22:42:23.349', '2022-10-20 22:42:23.349', 0, '', '', '留言模块', 0);
INSERT INTO `resource` VALUES (7, '2022-10-20 22:42:28.550', '2022-10-20 22:42:28.550', 0, '', '', '菜单模块', 0);
INSERT INTO `resource` VALUES (8, '2022-10-20 22:42:31.623', '2022-10-20 22:42:31.623', 0, '', '', '角色模块', 0);
INSERT INTO `resource` VALUES (9, '2022-10-20 22:42:36.262', '2022-10-20 22:42:36.262', 0, '', '', '评论模块', 0);
INSERT INTO `resource` VALUES (10, '2022-10-20 22:42:40.700', '2022-10-20 22:42:40.700', 0, '', '', '资源模块', 0);
INSERT INTO `resource` VALUES (11, '2022-10-20 22:42:51.023', '2022-10-20 22:42:51.023', 0, '', '', '博客信息模块', 0);
INSERT INTO `resource` VALUES (23, '2022-10-22 22:13:12.455', '2022-10-26 11:15:23.546', 10, '/resource/anonymous', 'PUT', '修改资源匿名访问', 0);
INSERT INTO `resource` VALUES (34, '2022-10-31 17:14:11.708', '2022-10-31 17:14:11.708', 10, '/resource', 'POST', '新增/编辑资源', 0);
INSERT INTO `resource` VALUES (35, '2022-10-31 17:14:42.320', '2022-10-31 17:18:52.508', 10, '/resource/list', 'GET', '资源列表', 0);
INSERT INTO `resource` VALUES (36, '2022-10-31 17:15:14.999', '2022-10-31 17:19:01.460', 10, '/resource/option', 'GET', '资源选项列表(树形)', 0);
INSERT INTO `resource` VALUES (37, '2022-10-31 17:16:56.830', '2022-10-31 17:16:56.830', 10, '/resource/:id', 'DELETE', '删除资源', 0);
INSERT INTO `resource` VALUES (38, '2022-10-31 17:19:28.905', '2022-10-31 17:19:28.905', 7, '/menu/list', 'GET', '菜单列表', 0);
INSERT INTO `resource` VALUES (39, '2022-10-31 18:46:33.051', '2022-10-31 18:46:33.051', 7, '/menu', 'POST', '新增/编辑菜单', 0);
INSERT INTO `resource` VALUES (40, '2022-10-31 18:46:53.804', '2022-10-31 18:46:53.804', 7, '/menu/:id', 'DELETE', '删除菜单', 0);
INSERT INTO `resource` VALUES (41, '2022-10-31 18:47:17.272', '2022-10-31 18:47:28.130', 7, '/menu/option', 'GET', '菜单选项列表(树形)', 0);
INSERT INTO `resource` VALUES (42, '2022-10-31 18:48:04.780', '2022-10-31 18:48:04.780', 7, '/menu/user/list', 'GET', '获取当前用户菜单', 0);
INSERT INTO `resource` VALUES (43, '2022-10-31 19:20:35.427', '2023-12-27 23:21:22.669', 3, '/article/list', 'GET', '文章列表', 0);
INSERT INTO `resource` VALUES (44, '2022-10-31 19:21:02.096', '2023-12-27 22:07:57.702', 3, '/article/:id', 'GET', '文章详情', 0);
INSERT INTO `resource` VALUES (45, '2022-10-31 19:26:04.763', '2022-10-31 19:26:09.709', 3, '/article', 'POST', '新增/编辑文章', 0);
INSERT INTO `resource` VALUES (46, '2022-10-31 19:26:36.453', '2022-10-31 19:26:36.453', 3, '/article/soft-delete', 'PUT', '软删除文章', 0);
INSERT INTO `resource` VALUES (47, '2022-10-31 19:26:52.344', '2022-10-31 19:26:52.344', 3, '/article', 'DELETE', '删除文章', 0);
INSERT INTO `resource` VALUES (48, '2022-10-31 19:27:07.731', '2022-10-31 19:27:07.731', 3, '/article/top', 'PUT', '修改文章置顶', 0);
INSERT INTO `resource` VALUES (49, '2022-10-31 20:19:55.588', '2022-10-31 20:19:55.588', 0, '', '', '分类模块', 0);
INSERT INTO `resource` VALUES (50, '2022-10-31 20:20:03.400', '2022-10-31 20:20:03.400', 0, '', '', '标签模块', 0);
INSERT INTO `resource` VALUES (51, '2022-10-31 20:22:03.799', '2022-10-31 20:22:03.799', 49, '/category/list', 'GET', '分类列表', 0);
INSERT INTO `resource` VALUES (52, '2022-10-31 20:22:28.840', '2022-10-31 20:22:28.840', 49, '/category', 'POST', '新增/编辑分类', 0);
INSERT INTO `resource` VALUES (53, '2022-10-31 20:31:04.577', '2022-10-31 20:31:04.577', 49, '/category', 'DELETE', '删除分类', 0);
INSERT INTO `resource` VALUES (54, '2022-10-31 20:31:36.612', '2022-10-31 20:31:36.612', 49, '/category/option', 'GET', '分类选项列表', 0);
INSERT INTO `resource` VALUES (55, '2022-10-31 20:32:57.112', '2022-10-31 20:33:13.261', 50, '/tag/list', 'GET', '标签列表', 0);
INSERT INTO `resource` VALUES (56, '2022-10-31 20:33:29.080', '2022-10-31 20:33:29.080', 50, '/tag', 'POST', '新增/编辑标签', 0);
INSERT INTO `resource` VALUES (57, '2022-10-31 20:33:39.992', '2022-10-31 20:33:39.992', 50, '/tag', 'DELETE', '删除标签', 0);
INSERT INTO `resource` VALUES (58, '2022-10-31 20:33:53.962', '2022-10-31 20:33:53.962', 50, '/tag/option', 'GET', '标签选项列表', 0);
INSERT INTO `resource` VALUES (59, '2022-10-31 20:35:05.647', '2022-10-31 20:35:05.647', 6, '/message/list', 'GET', '留言列表', 0);
INSERT INTO `resource` VALUES (60, '2022-10-31 20:35:25.551', '2022-10-31 20:35:25.551', 6, '/message', 'DELETE', '删除留言', 0);
INSERT INTO `resource` VALUES (61, '2022-10-31 20:36:20.587', '2022-10-31 20:36:20.587', 6, '/message/review', 'PUT', '修改留言审核', 0);
INSERT INTO `resource` VALUES (62, '2022-10-31 20:37:04.637', '2022-10-31 20:37:04.637', 9, '/comment/list', 'GET', '评论列表', 0);
INSERT INTO `resource` VALUES (63, '2022-10-31 20:37:29.779', '2022-10-31 20:37:29.779', 9, '/comment', 'DELETE', '删除评论', 0);
INSERT INTO `resource` VALUES (64, '2022-10-31 20:37:40.317', '2022-10-31 20:37:40.317', 9, '/comment/review', 'PUT', '修改评论审核', 0);
INSERT INTO `resource` VALUES (65, '2022-10-31 20:38:30.506', '2022-10-31 20:38:30.506', 8, '/role/list', 'GET', '角色列表', 0);
INSERT INTO `resource` VALUES (66, '2022-10-31 20:38:50.606', '2022-10-31 20:38:50.606', 8, '/role', 'POST', '新增/编辑角色', 0);
INSERT INTO `resource` VALUES (67, '2022-10-31 20:39:03.752', '2022-10-31 20:39:03.752', 8, '/role', 'DELETE', '删除角色', 0);
INSERT INTO `resource` VALUES (68, '2022-10-31 20:39:28.232', '2022-10-31 20:39:28.232', 8, '/role/option', 'GET', '角色选项', 0);
INSERT INTO `resource` VALUES (69, '2022-10-31 20:44:22.622', '2022-10-31 20:44:22.622', 0, '', '', '友链模块', 0);
INSERT INTO `resource` VALUES (70, '2022-10-31 20:44:41.334', '2022-10-31 20:44:41.334', 69, '/link/list', 'GET', '友链列表', 0);
INSERT INTO `resource` VALUES (71, '2022-10-31 20:45:01.150', '2022-10-31 20:45:01.150', 69, '/link', 'POST', '新增/编辑友链', 0);
INSERT INTO `resource` VALUES (72, '2022-10-31 20:45:12.406', '2022-10-31 20:45:12.406', 69, '/link', 'DELETE', '删除友链', 0);
INSERT INTO `resource` VALUES (74, '2022-10-31 20:46:48.330', '2022-10-31 20:47:01.505', 0, '', '', '用户信息模块', 0);
INSERT INTO `resource` VALUES (78, '2022-10-31 20:51:15.607', '2022-10-31 20:51:15.607', 74, '/user/list', 'GET', '用户列表', 0);
INSERT INTO `resource` VALUES (79, '2022-10-31 20:55:15.496', '2022-10-31 20:55:15.496', 11, '/setting/blog-config', 'GET', '获取博客设置', 0);
INSERT INTO `resource` VALUES (80, '2022-10-31 20:55:48.257', '2022-10-31 20:55:48.257', 11, '/setting/about', 'GET', '获取关于我', 0);
INSERT INTO `resource` VALUES (81, '2022-10-31 20:56:21.722', '2022-10-31 20:56:21.722', 11, '/setting/blog-config', 'PUT', '修改博客设置', 0);
INSERT INTO `resource` VALUES (82, '2022-10-31 21:57:30.021', '2022-10-31 21:57:30.021', 74, '/user/info', 'GET', '获取当前用户信息', 0);
INSERT INTO `resource` VALUES (84, '2022-10-31 22:06:18.348', '2022-10-31 22:07:38.004', 74, '/user', 'PUT', '修改用户信息', 0);
INSERT INTO `resource` VALUES (85, '2022-11-02 11:55:05.395', '2022-11-02 11:55:05.395', 11, '/setting/about', 'PUT', '修改关于我', 0);
INSERT INTO `resource` VALUES (86, '2022-11-02 13:20:09.485', '2022-11-02 13:20:09.485', 74, '/user/online', 'GET', '获取在线用户列表', 0);
INSERT INTO `resource` VALUES (91, '2022-11-03 16:42:31.558', '2022-11-03 16:42:31.558', 0, '', '', '操作日志模块', 0);
INSERT INTO `resource` VALUES (92, '2022-11-03 16:42:49.681', '2022-11-03 16:42:49.681', 91, '/operation/log/list', 'GET', '获取操作日志列表', 0);
INSERT INTO `resource` VALUES (93, '2022-11-03 16:43:04.906', '2022-11-03 16:43:04.906', 91, '/operation/log', 'DELETE', '删除操作日志', 0);
INSERT INTO `resource` VALUES (95, '2022-11-05 14:22:48.240', '2022-11-05 14:22:48.240', 11, '/home', 'GET', '获取后台首页信息', 0);
INSERT INTO `resource` VALUES (98, '2022-11-29 23:35:42.865', '2022-11-29 23:35:42.865', 74, '/user/offline', 'DELETE', '强制离线用户', 0);
INSERT INTO `resource` VALUES (99, '2022-12-07 20:48:05.939', '2022-12-07 20:48:05.939', 74, '/user/current/password', 'PUT', '修改当前用户密码', 0);
INSERT INTO `resource` VALUES (100, '2022-12-07 20:48:35.511', '2022-12-07 20:48:35.511', 74, '/user/current', 'PUT', '修改当前用户信息', 0);
INSERT INTO `resource` VALUES (101, '2022-12-07 20:55:08.271', '2022-12-07 20:55:08.271', 74, '/user/disable', 'PUT', '修改用户禁用', 0);
INSERT INTO `resource` VALUES (102, '2022-12-08 15:43:15.421', '2022-12-08 15:43:15.421', 0, '', '', '页面模块', 0);
INSERT INTO `resource` VALUES (103, '2022-12-08 15:43:26.009', '2022-12-08 15:43:26.009', 102, '/page/list', 'GET', '页面列表', 0);
INSERT INTO `resource` VALUES (104, '2022-12-08 15:43:38.570', '2022-12-08 15:43:38.570', 102, '/page', 'POST', '新增/编辑页面', 0);
INSERT INTO `resource` VALUES (105, '2022-12-08 15:43:50.879', '2022-12-08 15:43:50.879', 102, '/page', 'DELETE', '删除页面', 0);
INSERT INTO `resource` VALUES (106, '2022-12-16 11:53:57.989', '2022-12-16 11:53:57.989', 0, '', '', '文件模块', 0);
INSERT INTO `resource` VALUES (107, '2022-12-16 11:54:20.891', '2022-12-16 11:54:20.891', 106, '/upload', 'POST', '文件上传', 0);
INSERT INTO `resource` VALUES (108, '2022-12-18 01:34:47.800', '2022-12-18 01:34:47.800', 3, '/article/export', 'POST', '导出文章', 0);
INSERT INTO `resource` VALUES (109, '2022-12-18 01:34:59.255', '2022-12-18 01:34:59.255', 3, '/article/import', 'POST', '导入文章', 0);

-- ----------------------------
-- Table structure for role
-- ----------------------------
DROP TABLE IF EXISTS `role`;
CREATE TABLE `role`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `label` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `is_disable` tinyint(1) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `name`(`name` ASC) USING BTREE,
  UNIQUE INDEX `label`(`label` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of role
-- ----------------------------
INSERT INTO `role` VALUES (1, '2023-12-27 23:16:38.105', '2023-12-27 23:34:10.830', '管理员', 'admin', 0);
INSERT INTO `role` VALUES (2, '2023-12-27 23:16:50.687', '2023-12-29 23:13:46.460', '普通用户', 'user', 0);
INSERT INTO `role` VALUES (3, '2023-12-27 23:17:00.263', '2023-12-27 23:38:15.697', 'test', '测试用户', 0);

-- ----------------------------
-- Table structure for role_menu
-- ----------------------------
DROP TABLE IF EXISTS `role_menu`;
CREATE TABLE `role_menu`  (
  `menu_id` bigint NOT NULL,
  `role_id` bigint NOT NULL,
  PRIMARY KEY (`menu_id`, `role_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of role_menu
-- ----------------------------
INSERT INTO `role_menu` VALUES (2, 1);
INSERT INTO `role_menu` VALUES (2, 2);
INSERT INTO `role_menu` VALUES (2, 3);
INSERT INTO `role_menu` VALUES (3, 1);
INSERT INTO `role_menu` VALUES (3, 2);
INSERT INTO `role_menu` VALUES (3, 3);
INSERT INTO `role_menu` VALUES (4, 1);
INSERT INTO `role_menu` VALUES (4, 2);
INSERT INTO `role_menu` VALUES (4, 3);
INSERT INTO `role_menu` VALUES (5, 1);
INSERT INTO `role_menu` VALUES (5, 2);
INSERT INTO `role_menu` VALUES (5, 3);
INSERT INTO `role_menu` VALUES (6, 1);
INSERT INTO `role_menu` VALUES (6, 2);
INSERT INTO `role_menu` VALUES (6, 3);
INSERT INTO `role_menu` VALUES (8, 1);
INSERT INTO `role_menu` VALUES (8, 2);
INSERT INTO `role_menu` VALUES (8, 3);
INSERT INTO `role_menu` VALUES (9, 1);
INSERT INTO `role_menu` VALUES (9, 2);
INSERT INTO `role_menu` VALUES (9, 3);
INSERT INTO `role_menu` VALUES (10, 1);
INSERT INTO `role_menu` VALUES (10, 2);
INSERT INTO `role_menu` VALUES (10, 3);
INSERT INTO `role_menu` VALUES (16, 1);
INSERT INTO `role_menu` VALUES (16, 2);
INSERT INTO `role_menu` VALUES (16, 3);
INSERT INTO `role_menu` VALUES (17, 1);
INSERT INTO `role_menu` VALUES (17, 2);
INSERT INTO `role_menu` VALUES (17, 3);
INSERT INTO `role_menu` VALUES (23, 1);
INSERT INTO `role_menu` VALUES (23, 2);
INSERT INTO `role_menu` VALUES (23, 3);
INSERT INTO `role_menu` VALUES (24, 1);
INSERT INTO `role_menu` VALUES (24, 2);
INSERT INTO `role_menu` VALUES (24, 3);
INSERT INTO `role_menu` VALUES (25, 1);
INSERT INTO `role_menu` VALUES (25, 2);
INSERT INTO `role_menu` VALUES (25, 3);
INSERT INTO `role_menu` VALUES (26, 1);
INSERT INTO `role_menu` VALUES (26, 2);
INSERT INTO `role_menu` VALUES (26, 3);
INSERT INTO `role_menu` VALUES (27, 1);
INSERT INTO `role_menu` VALUES (27, 2);
INSERT INTO `role_menu` VALUES (27, 3);
INSERT INTO `role_menu` VALUES (28, 1);
INSERT INTO `role_menu` VALUES (28, 2);
INSERT INTO `role_menu` VALUES (28, 3);
INSERT INTO `role_menu` VALUES (29, 1);
INSERT INTO `role_menu` VALUES (29, 2);
INSERT INTO `role_menu` VALUES (29, 3);
INSERT INTO `role_menu` VALUES (30, 1);
INSERT INTO `role_menu` VALUES (30, 2);
INSERT INTO `role_menu` VALUES (30, 3);
INSERT INTO `role_menu` VALUES (31, 1);
INSERT INTO `role_menu` VALUES (31, 2);
INSERT INTO `role_menu` VALUES (31, 3);
INSERT INTO `role_menu` VALUES (32, 1);
INSERT INTO `role_menu` VALUES (32, 2);
INSERT INTO `role_menu` VALUES (32, 3);
INSERT INTO `role_menu` VALUES (33, 1);
INSERT INTO `role_menu` VALUES (33, 2);
INSERT INTO `role_menu` VALUES (33, 3);
INSERT INTO `role_menu` VALUES (34, 1);
INSERT INTO `role_menu` VALUES (34, 2);
INSERT INTO `role_menu` VALUES (34, 3);
INSERT INTO `role_menu` VALUES (36, 1);
INSERT INTO `role_menu` VALUES (36, 2);
INSERT INTO `role_menu` VALUES (36, 3);
INSERT INTO `role_menu` VALUES (37, 1);
INSERT INTO `role_menu` VALUES (37, 2);
INSERT INTO `role_menu` VALUES (37, 3);
INSERT INTO `role_menu` VALUES (38, 1);
INSERT INTO `role_menu` VALUES (38, 2);
INSERT INTO `role_menu` VALUES (38, 3);
INSERT INTO `role_menu` VALUES (39, 1);
INSERT INTO `role_menu` VALUES (39, 2);
INSERT INTO `role_menu` VALUES (39, 3);
INSERT INTO `role_menu` VALUES (47, 1);
INSERT INTO `role_menu` VALUES (48, 1);

-- ----------------------------
-- Table structure for role_resource
-- ----------------------------
DROP TABLE IF EXISTS `role_resource`;
CREATE TABLE `role_resource`  (
  `resource_id` bigint NOT NULL,
  `role_id` bigint NOT NULL,
  PRIMARY KEY (`resource_id`, `role_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of role_resource
-- ----------------------------
INSERT INTO `role_resource` VALUES (3, 1);
INSERT INTO `role_resource` VALUES (6, 1);
INSERT INTO `role_resource` VALUES (7, 1);
INSERT INTO `role_resource` VALUES (8, 1);
INSERT INTO `role_resource` VALUES (9, 1);
INSERT INTO `role_resource` VALUES (10, 1);
INSERT INTO `role_resource` VALUES (11, 1);
INSERT INTO `role_resource` VALUES (23, 1);
INSERT INTO `role_resource` VALUES (34, 1);
INSERT INTO `role_resource` VALUES (35, 1);
INSERT INTO `role_resource` VALUES (35, 2);
INSERT INTO `role_resource` VALUES (35, 3);
INSERT INTO `role_resource` VALUES (36, 1);
INSERT INTO `role_resource` VALUES (36, 2);
INSERT INTO `role_resource` VALUES (36, 3);
INSERT INTO `role_resource` VALUES (37, 1);
INSERT INTO `role_resource` VALUES (38, 1);
INSERT INTO `role_resource` VALUES (38, 2);
INSERT INTO `role_resource` VALUES (38, 3);
INSERT INTO `role_resource` VALUES (39, 1);
INSERT INTO `role_resource` VALUES (40, 1);
INSERT INTO `role_resource` VALUES (41, 1);
INSERT INTO `role_resource` VALUES (41, 2);
INSERT INTO `role_resource` VALUES (41, 3);
INSERT INTO `role_resource` VALUES (42, 1);
INSERT INTO `role_resource` VALUES (42, 2);
INSERT INTO `role_resource` VALUES (42, 3);
INSERT INTO `role_resource` VALUES (43, 1);
INSERT INTO `role_resource` VALUES (43, 2);
INSERT INTO `role_resource` VALUES (43, 3);
INSERT INTO `role_resource` VALUES (44, 1);
INSERT INTO `role_resource` VALUES (44, 2);
INSERT INTO `role_resource` VALUES (44, 3);
INSERT INTO `role_resource` VALUES (45, 1);
INSERT INTO `role_resource` VALUES (46, 1);
INSERT INTO `role_resource` VALUES (47, 1);
INSERT INTO `role_resource` VALUES (48, 1);
INSERT INTO `role_resource` VALUES (48, 2);
INSERT INTO `role_resource` VALUES (48, 3);
INSERT INTO `role_resource` VALUES (49, 1);
INSERT INTO `role_resource` VALUES (50, 1);
INSERT INTO `role_resource` VALUES (51, 1);
INSERT INTO `role_resource` VALUES (51, 2);
INSERT INTO `role_resource` VALUES (51, 3);
INSERT INTO `role_resource` VALUES (52, 1);
INSERT INTO `role_resource` VALUES (53, 1);
INSERT INTO `role_resource` VALUES (54, 1);
INSERT INTO `role_resource` VALUES (54, 2);
INSERT INTO `role_resource` VALUES (54, 3);
INSERT INTO `role_resource` VALUES (55, 1);
INSERT INTO `role_resource` VALUES (55, 2);
INSERT INTO `role_resource` VALUES (55, 3);
INSERT INTO `role_resource` VALUES (56, 1);
INSERT INTO `role_resource` VALUES (57, 1);
INSERT INTO `role_resource` VALUES (58, 1);
INSERT INTO `role_resource` VALUES (58, 2);
INSERT INTO `role_resource` VALUES (58, 3);
INSERT INTO `role_resource` VALUES (59, 1);
INSERT INTO `role_resource` VALUES (59, 2);
INSERT INTO `role_resource` VALUES (59, 3);
INSERT INTO `role_resource` VALUES (60, 1);
INSERT INTO `role_resource` VALUES (61, 1);
INSERT INTO `role_resource` VALUES (61, 2);
INSERT INTO `role_resource` VALUES (62, 1);
INSERT INTO `role_resource` VALUES (62, 2);
INSERT INTO `role_resource` VALUES (62, 3);
INSERT INTO `role_resource` VALUES (63, 1);
INSERT INTO `role_resource` VALUES (64, 1);
INSERT INTO `role_resource` VALUES (64, 2);
INSERT INTO `role_resource` VALUES (65, 1);
INSERT INTO `role_resource` VALUES (65, 2);
INSERT INTO `role_resource` VALUES (65, 3);
INSERT INTO `role_resource` VALUES (66, 1);
INSERT INTO `role_resource` VALUES (67, 1);
INSERT INTO `role_resource` VALUES (68, 1);
INSERT INTO `role_resource` VALUES (68, 2);
INSERT INTO `role_resource` VALUES (68, 3);
INSERT INTO `role_resource` VALUES (69, 1);
INSERT INTO `role_resource` VALUES (70, 1);
INSERT INTO `role_resource` VALUES (70, 2);
INSERT INTO `role_resource` VALUES (70, 3);
INSERT INTO `role_resource` VALUES (71, 1);
INSERT INTO `role_resource` VALUES (72, 1);
INSERT INTO `role_resource` VALUES (74, 1);
INSERT INTO `role_resource` VALUES (78, 1);
INSERT INTO `role_resource` VALUES (78, 2);
INSERT INTO `role_resource` VALUES (78, 3);
INSERT INTO `role_resource` VALUES (79, 1);
INSERT INTO `role_resource` VALUES (79, 2);
INSERT INTO `role_resource` VALUES (79, 3);
INSERT INTO `role_resource` VALUES (80, 1);
INSERT INTO `role_resource` VALUES (80, 2);
INSERT INTO `role_resource` VALUES (80, 3);
INSERT INTO `role_resource` VALUES (81, 1);
INSERT INTO `role_resource` VALUES (82, 1);
INSERT INTO `role_resource` VALUES (82, 2);
INSERT INTO `role_resource` VALUES (82, 3);
INSERT INTO `role_resource` VALUES (84, 1);
INSERT INTO `role_resource` VALUES (85, 1);
INSERT INTO `role_resource` VALUES (86, 1);
INSERT INTO `role_resource` VALUES (86, 2);
INSERT INTO `role_resource` VALUES (86, 3);
INSERT INTO `role_resource` VALUES (91, 1);
INSERT INTO `role_resource` VALUES (92, 1);
INSERT INTO `role_resource` VALUES (92, 2);
INSERT INTO `role_resource` VALUES (92, 3);
INSERT INTO `role_resource` VALUES (93, 1);
INSERT INTO `role_resource` VALUES (95, 1);
INSERT INTO `role_resource` VALUES (95, 2);
INSERT INTO `role_resource` VALUES (95, 3);
INSERT INTO `role_resource` VALUES (98, 1);
INSERT INTO `role_resource` VALUES (99, 1);
INSERT INTO `role_resource` VALUES (100, 1);
INSERT INTO `role_resource` VALUES (101, 1);
INSERT INTO `role_resource` VALUES (102, 1);
INSERT INTO `role_resource` VALUES (103, 1);
INSERT INTO `role_resource` VALUES (103, 2);
INSERT INTO `role_resource` VALUES (103, 3);
INSERT INTO `role_resource` VALUES (104, 1);
INSERT INTO `role_resource` VALUES (105, 1);
INSERT INTO `role_resource` VALUES (106, 1);
INSERT INTO `role_resource` VALUES (107, 1);
INSERT INTO `role_resource` VALUES (108, 1);
INSERT INTO `role_resource` VALUES (108, 2);
INSERT INTO `role_resource` VALUES (108, 3);
INSERT INTO `role_resource` VALUES (109, 1);

-- ----------------------------
-- Table structure for tag
-- ----------------------------
DROP TABLE IF EXISTS `tag`;
CREATE TABLE `tag`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `name`(`name` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of tag
-- ----------------------------
INSERT INTO `tag` VALUES (1, '2023-12-27 22:45:40.731', '2023-12-27 22:45:40.731', 'Golang');
INSERT INTO `tag` VALUES (2, '2023-12-27 22:46:36.082', '2023-12-27 22:46:36.082', 'Vue');
INSERT INTO `tag` VALUES (3, '2023-12-27 22:47:47.530', '2023-12-27 22:47:47.530', '感悟');

-- ----------------------------
-- Table structure for user_auth
-- ----------------------------
DROP TABLE IF EXISTS `user_auth`;
CREATE TABLE `user_auth`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `password` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `login_type` tinyint(1) NULL DEFAULT NULL COMMENT '登录类型',
  `ip_address` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '登录IP地址',
  `ip_source` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'IP来源',
  `last_login_time` datetime(3) NULL DEFAULT NULL,
  `is_disable` tinyint(1) NULL DEFAULT NULL,
  `is_super` tinyint(1) NULL DEFAULT NULL,
  `user_info_id` bigint NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `username`(`username` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 6 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_auth
-- ----------------------------
INSERT INTO `user_auth` VALUES (1, '2023-12-27 22:40:42.565', '2023-12-29 23:13:11.500', 'superadmin', '$2a$10$R1kus4SbUJ5QzAgcUuxrbOhv10j.CaVtUdmRbZ17C2552frAj7opW', 1, '172.18.45.12', '内网IP', '2023-12-29 23:13:11.500', 0, 1, 1);
INSERT INTO `user_auth` VALUES (2, '2022-10-31 21:54:11.040', '2023-12-27 23:44:06.416', 'admin', '$2a$10$urGRaFQoLoblBUUdgi1NCei3oGnCHJwtHFmVcIfC94135KdNffy4.', 1, '172.18.45.12', '内网IP', '2023-12-27 23:44:06.416', 0, 0, 2);
INSERT INTO `user_auth` VALUES (3, '2022-11-01 10:41:13.300', '2023-12-29 23:04:48.284', 'test@qq.com', '$2a$10$FmU4jxwDlibSL9pdt.AsuODkbB4gLp3IyyXeoMmW/XALtT/HdwTsi', 1, '172.18.45.12', '内网IP', '2023-12-29 23:04:48.284', 0, 0, 3);
INSERT INTO `user_auth` VALUES (4, '2022-10-19 22:31:26.805', '2023-12-26 21:10:35.334', 'user', '$2a$10$9vHpoeT7sF4j9beiZfPsOe0jJ67gOceO2WKJzJtHRZCjNJajl7Fhq', 1, '172.12.0.6:48716', '', '2022-12-24 12:13:52.494', 0, 0, 4);

-- ----------------------------
-- Table structure for user_auth_role
-- ----------------------------
DROP TABLE IF EXISTS `user_auth_role`;
CREATE TABLE `user_auth_role`  (
  `user_auth_id` bigint NOT NULL,
  `role_id` bigint NOT NULL,
  PRIMARY KEY (`user_auth_id`, `role_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_auth_role
-- ----------------------------
INSERT INTO `user_auth_role` VALUES (2, 1);
INSERT INTO `user_auth_role` VALUES (3, 2);
INSERT INTO `user_auth_role` VALUES (4, 3);

-- ----------------------------
-- Table structure for user_info
-- ----------------------------
DROP TABLE IF EXISTS `user_info`;
CREATE TABLE `user_info`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `email` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `nickname` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `avatar` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `intro` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `website` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `nickname`(`nickname` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_info
-- ----------------------------
INSERT INTO `user_info` VALUES (1, '2023-12-27 22:40:42.495', '2023-12-28 16:34:24.836', '', 'superadmin', 'public/uploaded/4c50eef3bdaf0b4164ce179e576f2b2d_20231228163423.gif', '这个人很懒，什么都没有留下', 'https://www.hahacode.cn');
INSERT INTO `user_info` VALUES (2, '2022-10-31 21:54:10.935', '2023-12-27 23:44:01.402', '', '管理员', 'https://www.bing.com/rp/ar_9isCNU2Q-VG1yEDDHnx8HAFQ.png', '我就是我，不一样的烟火', 'https://www.hahacode.cn');
INSERT INTO `user_info` VALUES (3, '2022-10-19 22:31:26.734', '2023-12-27 23:31:39.169', 'user@qq.com', '普通用户', 'https://www.bing.com/rp/ar_9isCNU2Q-VG1yEDDHnx8HAFQ.png', '我是个普通用户！', 'https://www.hahacode.cn');
INSERT INTO `user_info` VALUES (4, '2022-11-01 10:41:13.234', '2023-12-27 23:31:42.587', 'test@qq.com', '测试用户', 'https://www.bing.com/rp/ar_9isCNU2Q-VG1yEDDHnx8HAFQ.png', '我是测试用的！', 'https://www.hahacode.cn');

SET FOREIGN_KEY_CHECKS = 1;
