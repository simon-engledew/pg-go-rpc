CREATE EXTENSION IF NOT EXISTS dblink;

-- see docker-compose.yml #service
SELECT dblink_connect('rpc', 'host=service port=15432 user=admin dbname=main');

CREATE OR REPLACE FUNCTION rpc(value JSON) RETURNS JSON AS $$
SELECT response FROM dblink('rpc', 'RPC ' || value) AS t1(response JSON)
$$ LANGUAGE SQL VOLATILE;

SELECT rpc(json_build_object('moose', 'goose'));
