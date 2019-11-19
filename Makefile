run: build up

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down --remove-orphans

reup:
	docker-compose down --remove-orphans ;\
	docker-compose build ;\
	docker-compose up -d ;\

rmi:
	docker rmi $(docker images -a -q)

rm:
	docker rm $(docker ps -a -f status=exited -q)

test:
	docker-compose -f docker-compose.test.yml up --build -d ;\
	test_status=0 ;\
	docker-compose -f docker-compose.test.yml run integration_tests go test -v ./... || test_status=$$? ;\
	docker-compose -f docker-compose.test.yml down ;\
	echo "status="$$test_status;exit $$test_status ;\
