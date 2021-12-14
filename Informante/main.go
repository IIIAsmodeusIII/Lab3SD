// Menu
// GetServerIp
//  Connect to Broker
//  Wait Broker
// Send Action
//  Connect to Fulcrum
//  Wait Fulcrum

package main

import (
	"fmt"
	"log"
    "strings"
	"context"
	"bufio"
	"os"
	"flag"

	pb "github.com/IIIAsmodeusIII/Tarea3/Proto"
    "google.golang.org/grpc"
)

type PlanetaryData struct {
    planet string
	city string
    server int32
	address string
    version []int32
}

var server_files []PlanetaryData
var broker_address string


// ================================ Aux Func ================================ //
func failOnError(err error, msg string) {
  if err != nil {
    log.Fatalf("%s: %s", msg, err)
  }
}

func FindCity(planet, city string) int {

    for i, data := range server_files {
       if data.planet == planet && data.city == city{
           return i
       }
    }

    return -1
}

func ShowFiles(){

	for _, file := range server_files {
		fmt.Printf("[Data] %v.\n", file)
	}
}

// ============================ CRUD ======================================== //
func AddCity(input, planet string, city string) int{

	// Check register existence
    index := FindCity(planet, city)

    // If not, try to create
    if index == -1 {

		if SendCommand(input, server_files[FindCity(planet, city)].version) == -1{
			return -1
		}

        var new_ver []int32
        new_ver = append(new_ver, 0)
		new_ver = append(new_ver, 0)
		new_ver = append(new_ver, 0)

        new_file    := PlanetaryData{
            planet: planet,
			city: city,
			server: 0,
			address: "",
            version: new_ver,
        }

        server_files = append(server_files, new_file)
    }

	return len(server_files) - 1
}

func UpdateName(input, planet, city, new_name string) int{
	index := FindCity(planet, city)

	if index == - 1{
		return -1
	}

	if SendCommand(input, server_files[FindCity(planet, city)].version) == -1{
		return -1
	}

	server_files[index].city = new_name
	return 0
}

func UpdateNumber(input, planet, city string) int{
	return SendCommand(input, server_files[FindCity(planet, city)].version)
}

func DestroyCity(input, planet, city string) int{

	index := FindCity(planet, city)

	if index == -1{
		return -1
	}

	if SendCommand(input, server_files[FindCity(planet, city)].version) == -1{
		return -1
	}

	server_files[index] = server_files[len(server_files)-1]
	server_files        = server_files[:len(server_files)-1]

	return 0
}



// ============================ SERVER ====================================== //
func SendCommand(command string, version []int32) int{

	// Connect to Broker
	for i:=0; i<15; i++ {
		fmt.Printf("[SendCommand Request] Broker_address:%v.\n", broker_address)
		conn, err := grpc.Dial(broker_address, grpc.WithInsecure(), grpc.WithBlock())
		failOnError(err, "Problema al conectar al servidor.")
		defer conn.Close()

		c := pb.NewBrokerClient(conn)
		ctx := context.Background()
		r, err := c.GetServer(ctx, &pb.GetServerReq{
			Command: command,
			Version: version,
		})

		failOnError(err, "No se pudede acceder al servicio")

		// Send command
		conn, err = grpc.Dial(r.Address, grpc.WithInsecure(), grpc.WithBlock())
		failOnError(err, "Problema al conectar al servidor.")
		defer conn.Close()

		c2 := pb.NewFulcrumClient(conn)
		ctx2 := context.Background()
		r2, err := c2.CRUD(ctx2, &pb.Command{
			Command: command,
			Version: version,
		})

		if r2.Code == int32(-1) {
			continue
		}else{
			return 0
		}
	}

	return -1
}



// ========================================================================== //
func Menu(){

	fmt.Println("[*] Type AddCity Planet CityName [Rebelds] to create a city.")
	fmt.Println("[*] Type UpdateNumber Planet CityName NewRebelds to update ammount of rebelds.")
	fmt.Println("[*] Type UpdateName Planet CityName NewCityName to update a city name.")
	fmt.Println("[*] Type DestroyCity Planet CityName to destroy a city.")
    fmt.Println("[*] Type 'exit()' to close. Insert command: ")
    for true {
        fmt.Print("-> ")

        scanner := bufio.NewScanner(os.Stdin)
        scanner.Scan()
        input   := scanner.Text()

        data    := strings.Split(input, " ")
		command := data[0]
		planet  := data[1]
		city    := data[2]

        if input == "exit()" {
            break
        }

		if command == "AddCity"{
			AddCity(input, planet, city)
		}else if(command == "UpdateName"){
			UpdateName(input, planet, city, data[3])
		}else if(command == "DestroyCity"){
			DestroyCity(input, planet, city)
		}else if(command == "Files()"){
			ShowFiles()
		}else if(command == "UpdateNumber"){
			UpdateNumber(input, planet, city)
		}else{
			fmt.Printf("[->] Comando '%v' desconocido.\n", input)
		}
    }
}

func main(){

	// Get params
	address := flag.String("ba", "localhost:50049", "Address of Broker")
	flag.Parse()

	// Set broker address
	broker_address = *address

	Menu()
}
