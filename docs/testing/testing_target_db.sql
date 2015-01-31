create database if not exists test_target_db;
use test_target_db;
drop table if exists target;
create table target(
Id int not null primary key auto_increment,
Num int not null default 0 comment 'int',
String varchar(64) not null default 0 comment 'string',
LogTime datetime not null default 0 comment 'log time',
key IdxLogTime(LogTime)
)engine=innodb charset=utf8;
