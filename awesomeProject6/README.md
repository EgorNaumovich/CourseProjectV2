# HOWTO

## Build

build.bat:

```cmd
set CGO_ENABLED=0
set GOOS=linux
go build -a -installsuffix cgo -o main .
docker build -t transport-service -f dockerfile.scratch .
```

Run it:
```cmd
cd delivery-service
.\build.bat
cd ../transport-service
.\build.bat
```

## Network

Create if doesn`t exist:

```cmd
docker network create delivery-network
```

## Run services

Run Transport:

```cmd
docker run --rm --net delivery-network --name transport-service -d transport-service
```

Run delivery service:

```cmd
docker run --rm -p 8080:8080 -d --net delivery-network --name delivery-service delivery-service
```

Run database service:

```cmd
cd database
docker-compose up
```

Inspect network:
```cmd
docker network inspect delivery-network
```

If the `database_mysql_1` container is not listed in `"Containers"` section, run the command:
```cmd
docker network connect delivery-network <sql-container-id>
```

Run cli service:

```cmd
cd delivery-cli
go run main.go
```

## Connect to Adminer

Open Chrome web browser and connect `localhost:8081`.  Use the following credentials:

* `Имя пользователя=testuser`
* `Пароль=testpassword`
* `База данных=Deliveries`

Consider selecting `SQL-запрос` option and typing `select * from Deliveries;` statement to view db records.