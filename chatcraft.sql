/*
 Navicat Premium Data Transfer

 Source Server         : chat-craft
 Source Server Type    : MySQL
 Source Server Version : 80300
 Source Host           : localhost:3306
 Source Schema         : chatcraft

 Target Server Type    : MySQL
 Target Server Version : 80300
 File Encoding         : 65001

 Date: 27/01/2024 22:11:18
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for communities
-- ----------------------------
DROP TABLE IF EXISTS `communities`;
CREATE TABLE `communities`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `create_at` datetime(3) NULL DEFAULT NULL,
  `update_at` datetime(3) NULL DEFAULT NULL,
  `delete_at` datetime(3) NULL DEFAULT NULL,
  `name` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `owner_id` bigint UNSIGNED NULL DEFAULT NULL,
  `type` bigint NULL DEFAULT NULL,
  `image` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `desc` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_communities_delete_at`(`delete_at` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for relations
-- ----------------------------
DROP TABLE IF EXISTS `relations`;
CREATE TABLE `relations`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `create_at` datetime(3) NULL DEFAULT NULL,
  `update_at` datetime(3) NULL DEFAULT NULL,
  `delete_at` datetime(3) NULL DEFAULT NULL,
  `owner_id` bigint UNSIGNED NULL DEFAULT NULL,
  `target_id` bigint UNSIGNED NULL DEFAULT NULL,
  `type` bigint NULL DEFAULT NULL,
  `desc` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_relations_delete_at`(`delete_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 19 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user_basics
-- ----------------------------
DROP TABLE IF EXISTS `user_basics`;
CREATE TABLE `user_basics`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `create_at` datetime(3) NULL DEFAULT NULL,
  `update_at` datetime(3) NULL DEFAULT NULL,
  `delete_at` datetime(3) NULL DEFAULT NULL,
  `name` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `pass_word` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `avatar` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `gender` varchar(6) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT 'male',
  `phone` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `email` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `motto` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `identity` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `client_ip` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `client_port` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `salt` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `login_time` datetime(3) NULL DEFAULT NULL,
  `heart_beat_time` datetime(3) NULL DEFAULT NULL,
  `login_out_time` datetime(3) NULL DEFAULT NULL,
  `is_login_out` tinyint(1) NULL DEFAULT NULL,
  `device_info` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user_basics_delete_at`(`delete_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 16 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user_stories
-- ----------------------------
DROP TABLE IF EXISTS `user_stories`;
CREATE TABLE `user_stories`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `create_at` datetime(3) NULL DEFAULT NULL,
  `update_at` datetime(3) NULL DEFAULT NULL,
  `delete_at` datetime(3) NULL DEFAULT NULL,
  `owner_id` bigint UNSIGNED NULL DEFAULT NULL,
  `content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `media` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `type` bigint NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user_stories_delete_at`(`delete_at` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user_story_comments
-- ----------------------------
DROP TABLE IF EXISTS `user_story_comments`;
CREATE TABLE `user_story_comments`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `create_at` datetime(3) NULL DEFAULT NULL,
  `update_at` datetime(3) NULL DEFAULT NULL,
  `delete_at` datetime(3) NULL DEFAULT NULL,
  `user_story_id` bigint UNSIGNED NULL DEFAULT NULL,
  `comment_owner_id` bigint UNSIGNED NULL DEFAULT NULL,
  `comment_content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL,
  `type` bigint NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user_story_comments_delete_at`(`delete_at` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user_story_likes
-- ----------------------------
DROP TABLE IF EXISTS `user_story_likes`;
CREATE TABLE `user_story_likes`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `create_at` datetime(3) NULL DEFAULT NULL,
  `update_at` datetime(3) NULL DEFAULT NULL,
  `delete_at` datetime(3) NULL DEFAULT NULL,
  `user_story_id` bigint UNSIGNED NULL DEFAULT NULL,
  `like_owner_id` bigint UNSIGNED NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user_story_likes_delete_at`(`delete_at` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
