create database if not exists testdata;
use testdata;
drop table if exists web;
create table web(
Id int not null primary key auto_increment,
GroupName varchar(64) not null default '' comment '����ϵ��������',
Total int not null default 0 comment '��������',
Concurrence int not null default 0 comment '������',
RequestData text not null default '' comment '������������',
ServerParameters text not null default '' comment '��������������',
ElapseTime float not null default 0 comment '����ʱ',
QPS float not null default 0 comment 'ƽ��ÿ�����������',
TPQ float not null default 0 comment 'ƽ��ÿ��������ʱ(ms)',
LogTime datetime not null default 0 comment '��¼ʱ��',
key IdxGroupName(GroupName),
key IdxLogTime(LogTime)
)engine=innodb charset=utf8;

drop table if exists db;
create table db(
Id int not null primary key auto_increment,
GroupName varchar(64) not null default '' comment '����ϵ��������',
Total int not null default 0 comment '��������',
Concurrence int not null default 0 comment '������',
ExistDataTotal int unsigned not null default 0 comment 'Ŀ���Ѵ�������������Ҫ���ڲ��Ի�����������һ�������',
RequestData text not null default '' comment '������������',
ServerParameters text not null default '' comment '��������������',
ElapseTime float not null default 0 comment '����ʱ',
QPS float not null default 0 comment 'ƽ��ÿ�����������',
TPQ float not null default 0 comment 'ƽ��ÿ��������ʱ(ms)',
LogTime datetime not null default 0 comment '��¼ʱ��',
key IdxGroupName(GroupName),
key IdxLogTime(LogTime)
)engine=innodb charset=utf8;