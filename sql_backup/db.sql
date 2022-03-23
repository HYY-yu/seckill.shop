create table goods
(
	id int auto_increment
		primary key,
	name varchar(255) not null,
	`desc` varchar(500) default '' not null,
	count int default 0 not null,
	create_time int not null,
	delete_time int default 0 not null,
	constraint goods_index_0
		unique (name)
)
comment '商品表';

create table `order`
(
	id int auto_increment
		primary key,
	sec_id int not null,
	shop_id int not null,
	user_id int not null,
	create_time int not null
)
comment '订单表';

create table sec_kill
(
	id int auto_increment
		primary key,
	shop_id int not null,
	start_time int not null,
	end_time int default 0 not null,
	status tinyint not null,
	create_time int not null
)
comment '秒杀表';

create table users
(
	id int auto_increment
		primary key,
	user_name varchar(100) not null,
	constraint users_index_1
		unique (user_name)
)
comment '用户表';