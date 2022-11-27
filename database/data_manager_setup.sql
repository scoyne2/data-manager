CREATE TABLE feeds (
    id serial PRIMARY KEY,
    vendor TEXT,
    feed_name TEXT,
    feed_method TEXT
)
;

INSERT INTO feeds (vendor, feed_name, feed_method )
VALUES ('GoodRx', 'Claims', 'SFTP'),
('The Advisory Board', 'Physicians', 'S3')
;

CREATE TABLE feed_status (
    id serial PRIMARY KEY,
    feed_id INTEGER REFERENCES feeds (id),
    process_date TEXT,
    record_count INTEGER,
    error_count INTEGER,
    feed_status TEXT
)
;

INSERT INTO feed_status (feed_id, process_date, record_count, error_count, feed_status)
VALUES (1, '09/27/2018', 1958623, 0, 'success'),
(1, '09/23/2016', 2389635, 61, 'errors'),
(2, '10/15/2017', 0, 500, 'failed'),
(2, '03/24/2018', 13076, 0, 'success')
;