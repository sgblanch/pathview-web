SOURCES = go.mod go.sum main.go $(wildcard cmd/*.go) $(wildcard internal/*/*.go)
SOURCES_NODE = $(wildcard resources/*.js) $(wildcard resources/*.json)
SOURCES_WEB = $(SOURCES_NODE) $(wildcard template/*.html) 


DOCKER_IMAGES = $(shell awk '/^FROM/ { if ($$2 != "scratch") { print $$2 }}' deployments/Dockerfile deployments/Dockerfile.runner | sort -u)
DOCKER_DEPS = $(SOURCES) $(SOURCES_WEB)

READ_YML = 'YAML::load(STDIN.read)["services"].each{|key, value|puts key if !value.key?("build")}'

COMPOSE_SERVICES = $(shell ruby -ryaml -e $(READ_YML) < deployments/docker-compose.yml)
COMPOSE_OPTS = --env-file .env -f deployments/docker-compose.yml
COMPOSE_DEPS = $(DOCKER_DEPS) .env deployments/docker-compose.yml deployments/pathview.env
COMPOSE_DEV_OPTS = $(COMPOSE_OPTS) -f deployments/docker-compose.dev.yml
COMPOSE_DEV_DEPS  = $(COMPOSE_DEPS) deployments/docker-compose.dev.yml deployments/etc/nginx.dev.conf
COMPOSE_PROD_DEPS = $(COMPOSE_DEPS) deployments/etc/nginx.prod.conf deployments/ssl/dhparams.pem

.DEFAULT_GOAL: dev
.PHONY: dev prod pull clean stamps

dev: stamps .stamps/web .stamps/runner .stamps/node $(COMPOSE_DEV_DEPS)
	docker compose $(COMPOSE_DEV_OPTS) up -d

prod: stamps .stamps/web .stamps/runner $(COMPOSE_PROD_DEPS)
	docker compose $(COMPOSE_OPTS) up -d

.stamps/node: .stamps/pull resources/package.json | .stamps
	docker compose $(COMPOSE_DEV_OPTS) run --entrypoint "npm install" node
	touch $@

.stamps/pull: | .stamps
	$(foreach image, $(DOCKER_IMAGES), docker image pull -q $(image); )
	docker compose $(COMPOSE_DEV_OPTS) pull $(COMPOSE_SERVICES)
	touch $@

.stamps/runner: .stamps/pull .stamps/web deployments/Dockerfile.runner $(DOCKER_DEPS) $(SOURCES) scripts/*.R | .stamps
	docker build -t pathview-runner:latest -f deployments/Dockerfile.runner .
	touch $@

.stamps/web: .stamps/pull deployments/Dockerfile $(DOCKER_DEPS) $(SOURCES) $(SOURCES_WEB) | .stamps
	docker build -t pathview-web:latest -f deployments/Dockerfile .
	touch $@

deployments/ssl:
	install -m 0700 -d $@

deployments/ssl/dhparams.pem: | deployments/ssl
	openssl dhparam -out $@ 2048

vendor: go.mod go.sum
# go mod vendor

.stamps:
	install -m 0700 -d $@

stamps: | .stamps
	find .stamps -mindepth 1 -maxdepth 1 -type f -mtime +5d -delete

clean:
	docker compose $(COMPOSE_DEV_OPTS) down --volumes
	rm -rf .stamps vendor web/static
