DATABASE_URL ?= postgres://root:secret@scheduler_db:5434/prototype?sslmode=disable

HAS_DOCKER := $(shell command -v docker 2> /dev/null)
ifndef HAS_DOCKER
    $(error "Docker is not installed. Please install Docker.")
endif

# Target: check-docker
# Description: Check if Docker is installed
check-docker:
ifndef HAS_DOCKER
    $(error "Docker is not installed. Please install Docker.")
endif

# Target: docker-up
# Description: Bring up Docker containers in detached mode
docker-up: check-docker
	docker compose up -d --remove-orphans

# Target: docker-down
# Description: Bring down Docker containers
docker-down: check-docker
	docker compose down

# Target: docker-rebuild
# Description: Rebuild and bring up Docker containers in detached mode
docker-rebuild: check-docker
	docker compose up -d --build

# Target: docker-shell
# Description: Open a shell inside the server container
docker-shell: check-docker
	docker compose exec server /bin/sh

# Target: migrate-create
# Description: Create a new migration using goose
# Usage: make migrate-create name=<migration_name>
migrate-create: check-docker
ifdef name
	docker compose exec scheduler_server goose -dir database/migrations create ${name} sql
else
	$(error "Usage: make migrate-create name=<migration_name>")
endif

# Target: migrate-up
# Description: Apply database migrations using goose
migrate-up: check-docker
	docker compose exec scheduler_server goose -dir database/migrations up

# Target: migrate-down
# Description: Roll back database migrations using goose
migrate-down: check-docker
	docker compose exec scheduler_server goose -dir database/migrations down

# Target: migrate-redo
# Description: Roll back and reapply the latest database migration using goose
migrate-redo: check-docker
	docker compose exec scheduler_server goose -dir database/migrations redo
	
# Target: migrate-reset
# Description: Roll back all migrations using goose
migrate-reset: check-docker
	docker compose exec scheduler_server goose -dir database/migrations reset

# Target: migrate-status
# Description: Show the status of applied and pending migrations using goose
migrate-status: check-docker
	docker compose exec scheduler_server goose -dir database/migrations status

	
.PHONY: check-docker docker-up docker-down docker-rebuild docker-shell migrate-create migrate-up migrate-down migrate-redo migrate-reset migrate-status
