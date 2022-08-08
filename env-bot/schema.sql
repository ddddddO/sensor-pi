create table environment(
    id integer not null primary key autoincrement,
    date text not null,
    temperature real not null,
    pressure real not null,
    humidity real not null
);

create index date_idx on environment(date);

create table mh_z19(
    id integer not null primary key autoincrement,
    date text not null,
    co2 real not null
);

create index date_idx on mh_z19(date);
