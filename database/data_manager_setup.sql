CREATE TABLE feeds (
    id serial PRIMARY KEY,
    vendor TEXT,
    feed_name TEXT,
    feed_method TEXT
);

CREATE TABLE feed_status (
    id serial PRIMARY KEY,
    feed_id INTEGER REFERENCES feeds (id),
    process_date TIMESTAMP,
    record_count INTEGER,
    error_count INTEGER,
    feed_status TEXT
);

INSERT INTO feeds (vendor, feed_name, feed_method )
VALUES ('Coyne Enterprises', 'Orders', 'SFTP'),
('The Mochi Dog Board', 'Dog Breeds', 'S3');

INSERT INTO feed_status (feed_id, process_date, record_count, error_count, feed_status)
VALUES (1, '2021-03-31 22:30:20','YYYY-MM-DD HH:MI:SS', 1958623, 0, 'success'),
(1, '2017-03-31 14:30:20','YYYY-MM-DD HH:MI:SS', 2389635, 61, 'errors'),
(2, '2022-03-31 08:30:20','YYYY-MM-DD HH:MI:SS', 0, 500, 'failed'),
(2, '2022-03-31 11:30:20','YYYY-MM-DD HH:MI:SS', 13076, 0, 'success');

