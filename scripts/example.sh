#!/bin/bash


ADDR=""


docker-compose exec -T postgresql psql -U admin main <<'__EOF__'
CREATE EXTENSION IF NOT EXISTS dblink;

SELECT dblink_connect('rpc', 'host=192.168.0.101 port=15432 user=admin dbname=main');

CREATE OR REPLACE FUNCTION rpc(value JSON) RETURNS JSON AS $$
SELECT response FROM dblink('rpc', 'RPC ' || value) AS t1(response JSON)
$$ LANGUAGE SQL VOLATILE;

SELECT rpc(json_build_object('moose', 'goose'));
__EOF__