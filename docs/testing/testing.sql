create database if not exists testdata;
use testdata;
drop table if exists web;
create table web(
Id int not null primary key auto_increment,
GroupName varchar(64) not null default '' comment '测试系列组名称',
Total int not null default 0 comment '总请求量',
Concurrence int not null default 0 comment '并发量',
ElapseTime float not null default 0 comment '总用时',
QPS float not null default 0 comment '平均每秒完成请求数',
TPQ float not null default 0 comment '平均每次请求用时(ms)',
LogTime datetime not null default 0 comment '记录时间',
RequestData text not null default '' comment '测试请求数据',
ServerParameters text not null default '' comment '服务器参数设置',
key IdxGroupName(GroupName),
key IdxLogTime(LogTime)
)engine=innodb charset=utf8;

drop table if exists db_group;
create table db_group(
Id int not null primary key auto_increment,
Name varchar(128) not null default '' comment '测试项名称',
SettingParameters varchar(1024) not null default '' comment '参数设置',
LogTime datetime not null default 0 comment '记录时间',
key IdxName(Name),
key IdxLogTime(LogTime)
)engine=innodb charset=utf8;

drop table if exists db;
create table db(
Id int not null primary key auto_increment,
GroupId int not null default 0,
Name varchar(128) not null default '' comment '测试项名称',
Total int not null default 0 comment '总请求量',
Concurrence int not null default 0 comment '并发量',
ElapseTime float not null default 0 comment '总用时',
QPS float not null default 0 comment '平均每秒完成请求数',
TPQ float not null default 0 comment '平均每次请求用时(ms)',
LogTime datetime not null default 0 comment '记录时间',
key IdxName(Name),
key IdxLogTime(LogTime)
)engine=innodb charset=utf8;