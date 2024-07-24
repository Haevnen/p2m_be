include .env

run:
	docker-compose -f docker-compose.yml -p $(PROJECT_NAME) up -d# note: call scripts from /scripts
