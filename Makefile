SERVER_1_ADDRESS=localhost:50050
SERVER_2_ADDRESS=localhost:50051
SERVER_3_ADDRESS=localhost:50052

BROKER_PORT=:50049
BROKER_ADDRESS=localhost${BROKER_PORT}

BIN_NAME=program



## ========================================================================== ##
GRPC:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    Proto/services.proto

clean:
	rm ./Fulcrum/*.txt


## ========================================================================== ##
buildB:
	go build -o Broker/${BIN_NAME} main.go

runB:
	./Broker/${BIN_NAME} -s1=${SERVER_1_ADDRESS} -s2=${SERVER_2_ADDRESS} -s3=${SERVER_3_ADDRESS} -p=${BROKER_PORT}

startB:
	make buildB
	make runB

devB:
	go run Broker/main.go -s1=${SERVER_1_ADDRESS} -s2=${SERVER_2_ADDRESS} -s3=${SERVER_3_ADDRESS} -p=${BROKER_PORT}

# Requires:
# s1: Server address
# s2: Server address
# s3: Server address
# p : Process port



## ========================================================================== ##
buildF:
	go build -o Fulcrum/${BIN_NAME} main.go

runF:
	./Fulcrum/${BIN_NAME}

startF:
	make buildF
	make runF

devF:
	go run Fulcrum/main.go -sm=${SERVER_1_ADDRESS} -s1=${SERVER_2_ADDRESS} -s2=${SERVER_3_ADDRESS}

# Requires:
# p: Port to listen. 50050 default.
# sm: Server master address
# s1: Slave server 1 address
# s2: Slave server 2 address
# i: Server index (0, 1, 2)
# mt: Merge time. 120 default.



## ========================================================================== ##
buildI:
	go build -o Informante/${BIN_NAME} main.go

runI:
	./Informante/${BIN_NAME} -ba=${BROKER_ADDRESS}

startI:
	make buildI
	make runI

devI:
	go run Informante/main.go -ba=${BROKER_ADDRESS}

# Requires:
# ba: Broker address



## ========================================================================== ##
buildL:
	go build -o Leia/${BIN_NAME} main.go

runL:
	./Leia/${BIN_NAME} -ba=${BROKER_ADDRESS}

startL:
	make buildL
	make runL

devL:
	go run Leia/main.go -ba=${BROKER_ADDRESS}

# Requires:
# ba: Broker address



## ========================================================================== ##
