include .env

run:
	docker-compose -f docker-compose.yml -p $(PROJECT_NAME) up -d

migrate:
	docker-compose -f docker-compose-db-tools.yml run --rm flyway migrate
	docker-compose -f docker-compose.yml exec -e MYSQL_PWD=$(MYSQL_PASSWORD) db mysqldump -u$(MYSQL_USER) --no-tablespaces --skip-dump-date --no-data --ignore-table=$(MYSQL_DATABASE).flyway_schema_history $(MYSQL_DATABASE) > database/structure.sql

migrate-test:
	MYSQL_DATABASE=p2m_test MYSQL_USER=test MYSQL_PASSWORD=password DB_ENV=test docker-compose -f docker-compose-db-tools.yml run --rm flyway migrate

migrate-repair:
	docker-compose -f docker-compose-db-tools.yml run --rm flyway repair

codegen:
	# api
	mkdir -p "$(APP_API_DIR)"/gen/api
	docker-compose -f docker-compose-tools.yml run --rm oapi-codegen\
		-generate "types" -package api /spec/api_service.v1.yaml > "$(APP_API_DIR)"/gen/api/service.types.go
	docker-compose -f docker-compose-tools.yml run --rm oapi-codegen\
		-generate "gin-server,spec" -package api /spec/api_service.v1.yaml > "$(APP_API_DIR)"/gen/api/service.server.go
