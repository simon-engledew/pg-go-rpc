.PHONY: psql
psql:
	docker-compose exec postgresql psql -U admin main
