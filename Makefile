default: test

test: build-test
	docker-compose run test

build-test:
	docker-compose build test