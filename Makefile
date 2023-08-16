# Target: migrate-create
# Description: Create a new migration using goose
# Usage: make migrate-create name=<migration_name>
migrate-create: 
ifdef name
	goose -dir database/migrations create ${name} sql
else
	$(error "Usage: make migrate-create name=<migration_name>")
endif

# Target: migrate-up
# Description: Apply database migrations using goose
migrate-up: 
	goose -dir database/migrations postgres postgres://root:secret@localhost:5434/prototype?sslmode=disable up 

# Target: migrate-down
# Description: Roll back database migrations using goose
migrate-down: 
	goose -dir database/migrations postgres postgres://root:secret@localhost:5434/prototype?sslmode=disable down

# Target: migrate-redo
# Description: Roll back and reapply the latest database migration using goose
migrate-redo: 
	goose -dir database/migrations postgres postgres://root:secret@localhost:5434/prototype?sslmode=disable redo

# Target: migrate-status
# Description: Show the status of applied and pending migrations using goose
migrate-status: 
	goose -dir database/migrations postgres postgres://root:secret@localhost:5434/prototype?sslmode=disable status

	
.PHONY: migrate-create migrate-up migrate-down migrate-redo migrate-status
