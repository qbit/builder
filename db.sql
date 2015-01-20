DROP TABLE if exists status;
drop table if exists jobs;
drop table if exists arches;
drop table if exists diffs;

CREATE TABLE arches (id serial not null unique, arch text not null unique);
CREATE TABLE diffs (id serial not null unique, diffdata text not null);

CREATE TABLE jobs (
       id serial not null unique,
       created timestamp without time zone default now(),
       title text not null,
       descr text not null,
       port text not null,
       diff int not null references diffs (id) on delete cascade,
       active bool default true
);
CREATE TABLE status (
       id serial not null unique,
       created timestamp without time zone default now(),
       jid int not null references jobs (id) on delete cascade,
       status bool not null default false
);


insert into arches (arch) values ('i386');
insert into arches (arch) values ('amd64');

