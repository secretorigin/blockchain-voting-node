clear:
	docker rmi vb-node-img || true
	
clear-logs:
	rm -rf ./log
	rm -rf ./tests/log

build:
	docker build -t vb-node-img .

run-docker: build
	docker-compose up

stop-docker:
	docker-compose down