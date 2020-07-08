/*
 Navicat Premium Data Transfer

 Source Server         : gokins
 Source Server Type    : SQLite
 Source Server Version : 3030001
 Source Schema         : main

 Target Server Type    : SQLite
 Target Server Version : 3030001
 File Encoding         : 65001

 Date: 08/07/2020 15:23:00
*/

PRAGMA foreign_keys = false;

-- ----------------------------
-- Table structure for sys_user
-- ----------------------------
DROP TABLE IF EXISTS "sys_user";
CREATE TABLE "sys_user" (
  "id" integer,
  "xid" text NOT NULL,
  "name" text NOT NULL,
  "pass" text,
  "nick" text,
  "phone" text,
  "times" datetime,
  "logintm" datetime,
  "fwtm" datetime,
  "avat" text,
  PRIMARY KEY ("id")
);

-- ----------------------------
-- Records of sys_user
-- ----------------------------
INSERT INTO "sys_user" VALUES (1, 'admin', 'root', NULL, '超级管理员', NULL, 'CURRENT_TIMESTAMP', NULL, NULL, NULL);

-- ----------------------------
-- Indexes structure for table sys_user
-- ----------------------------
CREATE INDEX "IDX_sys_user_phone"
ON "sys_user" (
  "phone" ASC
);
CREATE INDEX "name"
ON "sys_user" (
  "name" ASC
);
CREATE INDEX "xid"
ON "sys_user" (
  "xid" ASC
);

PRAGMA foreign_keys = true;
