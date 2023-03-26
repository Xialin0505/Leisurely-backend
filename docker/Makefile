compose:
	docker compose -f ./docker-compose.yml up

createdb:
	docker exec -it backend-db-1 createdb --username=postgres --owner=postgres leisurely

dropdb:
	docker exec -it backend-db-1 dropdb --username=postgres leisurely

clean:
	docker container stop docker-backend-1 &
	docker container stop docker-db-1 &
	docker container rm docker-backend-1 &
	docker container rm docker-db-1 &
	docker image rm postgres &
	docker image rm docker_backend

.PHONY: compose createdb dropdb