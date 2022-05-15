create table environment(
    id integer not null primary key autoincrement,
    date text not null,
    temperature real not null,
    pressure real not null,
    humidity real not null
);

create index date_idx on environment(date);
