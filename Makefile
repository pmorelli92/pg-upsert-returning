run-db:
	docker run --name db -p 5432:5432 -e POSTGRES_HOST_AUTH_METHOD=trust postgres:latest

run-app:
	docker build -t pg-upsert-app:local .
	docker run -p 8080:8080 --link db --name app -i pg-upsert-app:local

run-bench:
	docker build -t benchmark:local -f benchmark/Dockerfile benchmark/.
	docker run --link app -t benchmark:local
