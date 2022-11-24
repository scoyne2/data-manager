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