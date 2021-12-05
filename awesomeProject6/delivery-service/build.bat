set CGO_ENABLED=0
set GOOS=linux
go build -a -installsuffix cgo -o main .
docker build -t delivery-service -f dockerfile.scratch .