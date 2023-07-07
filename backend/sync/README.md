## Getting Started

1. [Install Go here](https://go.dev/doc/install)

2. [Install Docker](https://docs.docker.com/get-docker/)

3. Install and run Temporal:

```
git clone https://github.com/temporalio/docker-compose.git
cd docker-compose
docker compose up
```

4. Build and run the sync worker

```
make
./bin/worker
```