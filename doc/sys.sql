/*
 Navicat Premium Data Transfer

 Source Server         : gokinsdb
 Source Server Type    : SQLite
 Source Server Version : 3030001
 Source Schema         : main

 Target Server Type    : SQLite
 Target Server Version : 3030001
 File Encoding         : 65001

 Date: 30/10/2020 11:08:09
*/

PRAGMA foreign_keys = false;

-- ----------------------------
-- Table structure for sys_param
-- ----------------------------
DROP TABLE IF EXISTS "sys_param";
CREATE TABLE "sys_param" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "key" varchar,
  "cont" blob,
  "times" datetime
);

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
-- Table structure for t_model
-- ----------------------------
DROP TABLE IF EXISTS "t_model";
CREATE TABLE "t_model" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "uid" varchar,
  "title" text,
  "desc" text,
  "times" datetime,
  "del" integer DEFAULT 0,
  "envs" text,
  "wrkdir" text,
  "clrdir" integer DEFAULT 0
);

-- ----------------------------
-- Table structure for t_model_run
-- ----------------------------
DROP TABLE IF EXISTS "t_model_run";
CREATE TABLE "t_model_run" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "uid" varchar,
  "tid" integer,
  "times" datetime,
  "timesd" datetime,
  "state" integer,
  "errs" text,
  "tgid" integer,
  "tgtyps" text
);

-- ----------------------------
-- Table structure for t_output
-- ----------------------------
DROP TABLE IF EXISTS "t_output";
CREATE TABLE "t_output" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "type" varchar(50),
  "tid" integer,
  "output" text,
  "times" datetime
);

-- ----------------------------
-- Table structure for t_plugin
-- ----------------------------
DROP TABLE IF EXISTS "t_plugin";
CREATE TABLE "t_plugin" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "tid" integer NOT NULL,
  "title" text,
  "type" integer DEFAULT 0,
  "para" text,
  "cont" text,
  "times" datetime,
  "sort" integer DEFAULT 100,
  "del" integer DEFAULT 0,
  "exend" integer DEFAULT 0
);

-- ----------------------------
-- Table structure for t_plugin_run
-- ----------------------------
DROP TABLE IF EXISTS "t_plugin_run";
CREATE TABLE "t_plugin_run" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "pid" integer,
  "mid" integer,
  "tid" integer,
  "times" datetime,
  "timesd" datetime,
  "state" integer,
  "excode" integer
);

-- ----------------------------
-- Table structure for t_trigger
-- ----------------------------
DROP TABLE IF EXISTS "t_trigger";
CREATE TABLE "t_trigger" (
  "id" integer NOT NULL,
  "uid" varchar,
  "types" varchar,
  "title" varchar,
  "desc" text,
  "times" date,
  "config" text,
  "del" integer DEFAULT 0,
  "enable" integer DEFAULT 0,
  "errs" text,
  "mid" integer,
  "meid" integer,
  "opt1" text,
  "opt2" text,
  "opt3" text,
  PRIMARY KEY ("id")
);

-- ----------------------------
-- Auto increment value for sys_param
-- ----------------------------

-- ----------------------------
-- Indexes structure for table sys_param
-- ----------------------------
CREATE INDEX "main"."key"
ON "sys_param" (
  "key" ASC
);

-- ----------------------------
-- Indexes structure for table sys_user
-- ----------------------------
CREATE INDEX "main"."IDX_sys_user_phone"
ON "sys_user" (
  "phone" ASC
);
CREATE INDEX "main"."name"
ON "sys_user" (
  "name" ASC
);
CREATE INDEX "main"."xid"
ON "sys_user" (
  "xid" ASC
);

-- ----------------------------
-- Auto increment value for t_model
-- ----------------------------

-- ----------------------------
-- Auto increment value for t_model_run
-- ----------------------------

-- ----------------------------
-- Auto increment value for t_output
-- ----------------------------

-- ----------------------------
-- Indexes structure for table t_output
-- ----------------------------
CREATE INDEX "main"."kv"
ON "t_output" (
  "type" ASC,
  "tid" ASC
);

-- ----------------------------
-- Auto increment value for t_plugin
-- ----------------------------

-- ----------------------------
-- Auto increment value for t_plugin_run
-- ----------------------------

PRAGMA foreign_keys = true;
