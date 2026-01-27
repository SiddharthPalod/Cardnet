CREATE KEYSPACE IF NOT EXISTS cardnet_ledger 
WITH REPLICATION = { 
   'class' : 'SimpleStrategy', 
   'replication_factor' : 1 
};

USE cardnet_ledger;

CREATE TABLE IF NOT EXISTS auth_events_by_merchant (
    merchant_id text,
    event_time timestamp,
    auth_id text,
    event_id uuid,
    event_type text,
    decision text,
    payload blob,
    PRIMARY KEY ((merchant_id), event_time, auth_id)
) WITH CLUSTERING ORDER BY (event_time DESC);

CREATE TABLE IF NOT EXISTS auth_events_by_card (
    card_hash text,
    event_time timestamp,
    auth_id text,
    merchant_id text,
    event_id uuid,
    event_type text,
    decision text,
    payload blob,
    PRIMARY KEY ((card_hash), event_time, auth_id)
) WITH CLUSTERING ORDER BY (event_time DESC);

CREATE TABLE IF NOT EXISTS auth_events_by_time (
    day_bucket timestamp,
    event_time timestamp,
    auth_id text,
    merchant_id text,
    card_hash text,
    event_type text,
    decision text,
    payload blob,
    event_id uuid,
    PRIMARY KEY ((day_bucket), event_time, auth_id)
) WITH CLUSTERING ORDER BY (event_time DESC);
