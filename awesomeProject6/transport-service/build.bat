set CGO_ENABLED=0
set GOOS=linux
go build -a -installsuffix cgo -o main .
docker build -t transport-service -f dockerfile.scratch .