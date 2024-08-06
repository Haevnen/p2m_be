include .env

run:
	docker-compose -f docker-compose.yml -p $(PROJECT_NAME) up -d

down:
	docker-compose -f docker-compose.yml -p $(PROJECT_NAME) down --remove-orphans -v

migrate:
	docker-compose -f docker-compose-db-tools.yml run --rm flyway migrate
	docker-compose -f docker-compose.yml exec -e MYSQL_PWD=$(MYSQL_PASSWORD) db mysqldump -u$(MYSQL_USER) --no-tablespaces --skip-dump-date --no-data --ignore-table=$(MYSQL_DATABASE).flyway_schema_history $(MYSQL_DATABASE) > database/structure.sql
migrate-test:
	MYSQL_DATABASE=p2m_test MYSQL_USER=test MYSQL_PASSWORD=password DB_ENV=test docker-compose -f docker-compose-db-tools.yml run --rm flyway migrate

migrate-repair:
	docker-compose -f docker-compose-db-tools.yml run --rm flyway repair

UNIFIED_DIR := gen/unified-openapi
codegen-unify:
	docker-compose -f docker-compose-tools.yml run --rm openapi-generator-cli generate -g openapi-yaml -i /api/api_service.v1.yaml -o /api/$(UNIFIED_DIR)/api_service

codegen: codegen-unify
	# api
	mkdir -p "$(APP_API_DIR)"/gen/api
	docker-compose -f docker-compose-tools.yml run --rm oapi-codegen\
		-generate "types" -package api /api/$(UNIFIED_DIR)/api_service/openapi/openapi.yaml > "$(APP_API_DIR)"/gen/api/service.types.go
	docker-compose -f docker-compose-tools.yml run --rm oapi-codegen\
		-generate "gin-server,spec" -package api /api/$(UNIFIED_DIR)/api_service/openapi/openapi.yaml > "$(APP_API_DIR)"/gen/api/service.server.go