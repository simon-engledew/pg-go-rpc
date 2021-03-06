proof of concept for postgres RPC mechanism based on dblink.

unlike [postgresql-rpc](https://github.com/simon-engledew/postgresql-rpc), this should work on most major cloud providers.

useful when creating a unidirectional data flow that involves plpgsql.

golang pretends to be a postgresql-compatible server using a modified example from [pgproto3](github.com/jackc/pgproto3/v2), you can then send it RPC messages and get JSON responses back:

```
CREATE EXTENSION IF NOT EXISTS dblink;

SELECT dblink_connect('rpc', 'host=service port=15432 user=admin dbname=main');

CREATE OR REPLACE FUNCTION rpc(value JSON) RETURNS JSON AS $$
SELECT response FROM dblink('rpc', 'RPC ' || value) AS t1(response JSON)
$$ LANGUAGE SQL VOLATILE;

SELECT rpc(json_build_object('moose', 'goose'));
```

To see it in action, run `make up` to build and start the services using Docker. Then run `make test` to fire a test RPC call.