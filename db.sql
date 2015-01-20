drop table if exists stat cascade;
drop table if exists jobs;
drop table if exists arches;
drop table if exists diffs;

CREATE TABLE arches (id serial not null unique, arch text not null unique);
CREATE TABLE diffs (id serial not null unique, diffdata text not null);
CREATE TABLE stat (
       id serial not null unique,
       status text unique
);
CREATE TABLE jobs (
       id serial not null unique,
       created timestamp without time zone default now(),
       title text not null,
       descr text not null,
       port text not null,
       diff int not null references diffs (id) on delete cascade,
       status int not null default 1 references stat (id) on delete cascade,
       active bool default true
);


insert into arches (arch) values ('i386');
insert into arches (arch) values ('amd64');

insert into stat (status) values ('Grabable');
insert into stat (status) values ('Pending');
insert into stat (status) values ('Building');
insert into stat (status) values ('Failed');


