PROJECT?=github.com/donskova1ex/1cServices
API_NAME?=api
API_VERSION?=0.0.1
API_CONTAINER_NAME?=docker.io/donskova1ex/${API_NAME}

# for local uses
api_local_build:
	go build -o bin/${API_NAME} cmd/${API_NAME}/${API_NAME}.go

# docker
api_docker_build:
	docker build -t ${API_CONTAINER_NAME}:${API_VERSION} -t ${API_CONTAINER_NAME}:latest -f Dockerfile.api .
api_docker_push:
	docker push ${API_CONTAINER_NAME}:${API_VERSION}
	docker push ${API_CONTAINER_NAME}:latest
api_docker_update: api_docker_build api_docker_push