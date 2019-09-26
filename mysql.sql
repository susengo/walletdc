
create database baas_api;

use baas_api;

-- auto-generated definition
create table user
(
  id       int auto_increment primary key,
  account  varchar(30)  not null,
  password varchar(100) not null,
  avatar   varchar(200) null,
  name     varchar(20)  not null,
  created      bigint not null,
  constraint user_account_uindex
  unique (account)
)  ENGINE=InnoDB  DEFAULT CHARSET=utf8 comment '用户表';

-- auto-generated definition
create table role
(
  rkey        varchar(20)  not null primary key,
  name        varchar(40)  not null,
  description varchar(200) null
)ENGINE=InnoDB  DEFAULT CHARSET=utf8 comment '角色表';

-- admin 123456
INSERT INTO baas_api.user (id, account, password, avatar, name, created) VALUES (1, 'admin', 'pbkdf2_sha256$180000$JEavgdkTBzU3$3pIgoygm1QBtgbfEHeWZ7H4O2rEIgkgLxYV48mE+J4M=', 'https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif', 'admin', 1557977199);

INSERT INTO baas_api.role (rkey, name, description) VALUES ('admin', '管理员', '超级管理员,拥有所有权限');
INSERT INTO baas_api.role (rkey, name, description) VALUES ('user', '用户', '普通用户');

INSERT INTO baas_api.user_role (user_id, role_key) VALUES (1, 'admin');