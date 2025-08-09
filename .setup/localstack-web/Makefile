STACK_NAME=localstack-web

start:
	docker-compose -p ${STACK_NAME} up -d --build --remove-orphans

stop:
	docker-compose -p ${STACK_NAME} stop

restart: stop start

clean:
	docker-compose -p ${STACK_NAME} down -v

logs:
	docker-compose -p ${STACK_NAME} logs -f

run:
	npm run dev

build:
	npm run build

preview:
	npm run preview

lint:
	npm run lint

preview:
	npm run preview
