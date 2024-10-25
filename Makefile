docker:
	docker compose up -d mysql

.PHONY: seed
seed:
	cd seed && go run .

.PHONY: run
run:
	cd app && go run .
