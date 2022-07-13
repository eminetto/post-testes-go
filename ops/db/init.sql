create database workshop;
grant all privileges on workshop.* to workshop@'%' identified by 'workshop';
use workshop;
create table if not exists person (id int AUTO_INCREMENT,first_name varchar(100), last_name varchar(100), created_at datetime, updated_at datetime, PRIMARY KEY (`id`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;
insert into person values (1, "Elton", "Minetto",now(), null);