proof of concept for postgres RPC mechanism based on dblink.

unlike [postgresql-rpc](https://github.com/simon-engledew/postgresql-rpc), this should work on most major cloud providers.

useful when creating a unidirectional data flow that involves plpgsql.
