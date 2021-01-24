CREATE DATABASE IF NOT EXISTS goapi;
use goapi;

CREATE TABLE IF NOT EXISTS WEATHER
(
dat char(8) NOT NULL,
weather int NOT NULL,
location_id int NOT NULL,
comment varchar(255),
PRIMARY KEY(dat, location_id)
);

CREATE TABLE IF NOT EXISTS LOCATION
(
id int PRIMARY KEY,
city varchar(20) NOT NULL
);

INSERT INTO LOCATION VALUES
(1, "新宿");
