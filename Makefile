build:
	docker build . -t lazhevskyi/hash-sv:latest

run:
	docker run -p 8080:80 -p 8081:81 -e HASH_TTL=5m lazhevskyi/hash-sv:latest