buildGRPC:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    Proto/services.proto


buildB:
	go build -o Broker/bin/main main.go

runB:
	go run Broker/main.go

startB:
	make buildB
	make runB


buildF:
	go build -o Fulcrum/bin/main main.go

runF:
	go run Fulcrum/main.go

startF:
	make buildF
	make runF


buildI:
	go build -o Informante/bin/main main.go

runI:
	go run Informante/main.go

startI:
	make buildI
	make runI


buildL:
	go build -o Leia/bin/main main.go

runL:
	go run Leia/main.go

startL:
	make buildL
	make runL

start:
	xterm -e startB
	xterm -e startL
