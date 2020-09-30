package comm

import (
	"bytes"
	"gokins/model"

	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

func InitDb() error {
	db, err := xorm.NewEngine("sqlite3", Dir+"/db.dat")
	if err != nil {
		return err
	}
	Db = db
	isext, err := Db.IsTableExist(model.SysUser{})
	if err == nil && !isext {
		_, err := Db.Import(bytes.NewBufferString(sqls))
		if err != nil {
			println("Db.Import err:" + err.Error())
		}
		//e:=&models.SysUser{}
		//e.Times=time.Now()
		//db.Cols("times").Where("xid=?","admin").Update(e)
	}
	return nil
}

const sqls = `
/*
 Navicat Premium Data Transfer

 Source Server         : gokinsdb
 Source Server Type    : SQLite
 Source Server Version : 3030001
 Source Schema         : main

 Target Server Type    : SQLite
 Target Server Version : 3030001
 File Encoding         : 65001

 Date: 30/09/2020 15:39:55
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
BEGIN;
INSERT INTO "sys_user" VALUES (1, 'admin', 'root', NULL, '超级管理员', NULL, '2020-07-08 07:25:53', NULL, NULL, NULL);
COMMIT;

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
  "del" integer DEFAULT 0
);

-- ----------------------------
-- Records of t_model
-- ----------------------------
BEGIN;
COMMIT;

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
  "state" integer
);

-- ----------------------------
-- Records of t_model_run
-- ----------------------------
BEGIN;
COMMIT;

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
-- Records of t_output
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for t_plungin
-- ----------------------------
DROP TABLE IF EXISTS "t_plungin";
CREATE TABLE "t_plungin" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "tid" integer NOT NULL,
  "title" text,
  "type" integer DEFAULT 0,
  "para" text,
  "cont" text,
  "times" datetime
);

-- ----------------------------
-- Records of t_plungin
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for t_plungin_run
-- ----------------------------
DROP TABLE IF EXISTS "t_plungin_run";
CREATE TABLE "t_plungin_run" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "pid" integer,
  "tid" integer,
  "times" datetime,
  "timesd" datetime,
  "state" integer,
  "excode" integer,
  "output" text
);

-- ----------------------------
-- Records of t_plungin_run
-- ----------------------------
BEGIN;
COMMIT;

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
-- Auto increment value for t_plungin
-- ----------------------------

-- ----------------------------
-- Auto increment value for t_plungin_run
-- ----------------------------

PRAGMA foreign_keys = true;

`
