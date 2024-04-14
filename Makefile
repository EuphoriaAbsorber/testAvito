build:
	sudo docker compose up -d --build

build-it:
	sudo docker compose up --build

stop:docker-fix
	sudo docker-compose down

container-prune:
	sudo docker container prune -f

image-prune:
	sudo docker image prune -f

docker-postgres-bash:
	sudo docker exec -it postgres bash

docker-prune-all:
	sudo docker system prune -a

docker-fix:
	sudo killall containerd-shim