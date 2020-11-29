/*
 Navicat MySQL Data Transfer

 Source Server         : 127.0.0.1
 Source Server Type    : MySQL
 Source Server Version : 50725
 Source Host           : localhost:3306
 Source Schema         : nonsense

 Target Server Type    : MySQL
 Target Server Version : 50725
 File Encoding         : 65001

 Date: 27/11/2020 15:27:48
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for app
-- ----------------------------
DROP TABLE IF EXISTS `app`;
CREATE TABLE `app`  (
                        `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
                        `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'app 名称',
                        `private_key` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '私钥',
                        `create_time` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                        `update_time` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '更新时间',
                        PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '应用配置表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for device
-- ----------------------------
DROP TABLE IF EXISTS `device`;
CREATE TABLE `device` (
                          `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                          `device_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '设备id',
                          `app_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT 'app_id',
                          `user_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '账户id',
                          `type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '设备类型,1:Android；2：IOS；3：Windows; 4：MacOS；5：Web',
                          `brand` varchar(20) NOT NULL DEFAULT '' COMMENT '手机厂商',
                          `model` varchar(20) NOT NULL DEFAULT '' COMMENT '机型',
                          `system_version` varchar(10) NOT NULL DEFAULT '' COMMENT '系统版本',
                          `sdk_version` varchar(10) NOT NULL DEFAULT '' COMMENT 'app版本',
                          `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '在线状态，0：离线；1：在线',
                          `conn_id` varchar(25) NOT NULL DEFAULT '' COMMENT '连接层服务器id',
                          `user_ip` varchar(30) NOT NULL DEFAULT '' COMMENT '用户ip',
                          `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                          `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                          PRIMARY KEY (`id`) USING BTREE,
                          UNIQUE KEY `uk_device_id` (`device_id`) USING BTREE,
                          KEY `idx_app_id_user_id` (`app_id`,`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='设备注册表';

-- ----------------------------
-- Table structure for group
-- ----------------------------
DROP TABLE IF EXISTS `group`;
CREATE TABLE `group`  (
                          `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                          `app_id` bigint(20) NOT NULL DEFAULT 0 COMMENT 'app_id',
                          `group_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '群组id',
                          `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '群组名称',
                          `introduction` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '群组简介',
                          `user_num` int(11) NOT NULL DEFAULT 0 COMMENT '群组人数',
                          `type` tinyint(4) NOT NULL DEFAULT 0 COMMENT '群组类型，1：小群；2：大群',
                          `extra` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '附加属性',
                          `create_time` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                          `update_time` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '更新时间',
                          PRIMARY KEY (`id`) USING BTREE,
                          UNIQUE INDEX `uk_app_id_group_id`(`app_id`, `group_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '群组信息表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for group_user
-- ----------------------------
DROP TABLE IF EXISTS `group_user`;
CREATE TABLE `group_user`  (
                               `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                               `app_id` bigint(20) NOT NULL DEFAULT 0 COMMENT 'app_id',
                               `group_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '组id',
                               `user_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户id',
                               `label` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户在群组的昵称',
                               `extra` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '附加属性',
                               `create_time` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                               `update_time` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '更新时间',
                               PRIMARY KEY (`id`) USING BTREE,
                               UNIQUE INDEX `uk_app_id_group_id_user_id`(`app_id`, `group_id`, `user_id`) USING BTREE,
                               INDEX `idx_app_id_user_id`(`app_id`, `user_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '群组成员表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for message
-- ----------------------------
DROP TABLE IF EXISTS `message`;
CREATE TABLE `message` (
                           `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                           `app_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT 'app_id',
                           `message_id` varchar(64) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '消息id,后续会使用全局id，暂时未用',
                           `receiver_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '接收者id私聊为user_id，群组消息为group_id',
                           `seq` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '消息序列号',
                           `sender_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '发送者id',
                           `receiver_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '接收者类型,1:个人；2：群组',
                           `type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '消息类型',
                           `object_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '所属类型的id',
                           `object_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '所属类型，1：用户；2：群组',
                           `sender_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '发送者类型',
                           `to_user_ids` varchar(255) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '需要@的用户id列表，多个用户用，隔开',
                           `sender_device_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '发送设备id',
                           `content` varchar(4094) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '消息内容',
                           `send_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '消息发送时间',
                           `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '消息状态，0：未处理1：消息撤回',
                           `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                           `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                           PRIMARY KEY (`id`,`sender_type`),
                           UNIQUE KEY `idx_app_id_recid_seq` (`app_id`,`seq`,`receiver_id`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC COMMENT='消息';

-- ----------------------------
-- Table structure for uid
-- ----------------------------
DROP TABLE IF EXISTS `uid`;
CREATE TABLE `uid` (
                       `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                       `app_id` varchar(128) NOT NULL DEFAULT '' COMMENT '业务id',
                       `max_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '最大id',
                       `step` int(10) unsigned NOT NULL DEFAULT '1000' COMMENT '步长',
                       `description` varchar(255) NOT NULL DEFAULT '' COMMENT '描述',
                       `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                       `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                       PRIMARY KEY (`id`) USING BTREE,
                       UNIQUE KEY `uk_business_id` (`app_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='分布式自增主键';

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
                        `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                        `app_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT 'app_id',
                        `user_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '用户id',
                        `passwd` varchar(100) DEFAULT '' COMMENT '用户密码',
                        `nickname` varchar(20) NOT NULL DEFAULT '' COMMENT '昵称',
                        `sex` tinyint(4) NOT NULL DEFAULT '1' COMMENT '性别，0:未知；1:男；2:女',
                        `birthday` varchar(20) NOT NULL DEFAULT '0',
                        `mobile` varchar(20) NOT NULL DEFAULT '',
                        `avatar_url` varchar(50) NOT NULL DEFAULT '' COMMENT '用户头像链接',
                        `email` varchar(100) DEFAULT '' COMMENT '用户邮箱',
                        `extra` varchar(1024) NOT NULL DEFAULT '' COMMENT '附加属性',
                        `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                        `last_login_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP  COMMENT '最后登录时间',
                        `last_login_ip` varchar(255) NOT NULL DEFAULT '' COMMENT '最后登录ip',
                        `register_ip` varchar(255) NOT NULL DEFAULT '' COMMENT '注册时ip',
                        `weixin_openid` varchar(50) NOT NULL DEFAULT '' COMMENT '微信开放id',
                        PRIMARY KEY (`id`) USING BTREE,
                        UNIQUE KEY `uk_app_id_user_id` (`app_id`,`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='用户信息表';

-- ----------
CREATE TABLE `user_seq` (
                            `id` bigint(20) NOT NULL AUTO_INCREMENT,
                            `app_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '应用id',
                            `group_id` bigint(20) DEFAULT '0' COMMENT '超大群id',
                            `user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '用户id',
                            `read_seq` bigint(20) DEFAULT '1' COMMENT '已读消息偏移',
                            `receive_seq` bigint(20) DEFAULT '0' COMMENT '收到消息总偏移',
                            PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------

SET FOREIGN_KEY_CHECKS = 1;
