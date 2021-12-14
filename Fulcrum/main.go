package main

import (
    "fmt"
	"net"
	"log"
    "os"
    "strconv"
    "math/rand"
    "io/ioutil"
    "io"
    "strings"
    "time"
    "path/filepath"
    "flag"

	"context"

	"google.golang.org/grpc"
	pb "github.com/IIIAsmodeusIII/Tarea3/Proto"
)

type server struct {
	pb.UnimplementedFulcrumServer
}

type PlanetaryData struct {
    name string
    version []int32
}

var server_files []PlanetaryData
var servers_ips [3]string
var server_index int32
var server_block bool

// ================================ Aux Func ================================ //
func failOnError(err error, msg string) {
  if err != nil {
    log.Fatalf("%s: %s", msg, err)
  }
}

func FindFile(file string) int {

    for i, data := range server_files {
       if data.name == file {
           return i
       }
    }

    return -1
}

func UpdateClock(file string) PlanetaryData{
    index := FindFile(file)
    if index == -1 {
        return PlanetaryData{}
    }

    server_files[index].version[int(server_index)] += int32(1)
    return server_files[index]
}

func Executer(commands []string) []string{

    var failed []string

    for _, command := range commands {

        if (command != "") {
            // Get command data
            fmt.Printf("[DEBUG] ExecuterComando: %v .\n", command)
            data := strings.Split(command, " ")
            comm := data[0]
            planet := data[1] + ".txt"
            city := data[2]

            var status int

            if (comm == "AddCity"){
                if len(data) == 4 {
                    status = AddCity(command, planet, city, data[3])
                }else{
                    status = AddCity(command, planet, city, "0")
                }
            }else if(comm == "UpdateName"){
                status = UpdateName(command, planet, city, data[3])
            }else if(comm == "UpdateNumber"){
                status = UpdateNumber(command, planet, city, data[3])
            }else if(comm == "DestroyCity"){
                status = DeleteCity(command, planet, city)
            }else{
                continue
            }

            if status == -1 {
                failed = append(failed, command)
            }
        }
    }
    return failed
}

func ShowFiles(){

	for _, file := range server_files {
		fmt.Printf("[Data] %v.\n", file)
	}
}

// =========================== Registro planetario ========================== //
func LineExist(file, city_name string) int{
    /*
    Function: Returns index of a line in certain file
    Input:
        file: file name, string. (e.g. example_file.txt)
        city_name: city_name to search (e.g. Mos_Pelgo)
    Output:
        if city exist, returns line_number. Else, returs -1.
    */
    index := FindFile(file)
    if index == -1 {
        return -1
    }

    f, err := os.OpenFile("./Fulcrum/" + file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
    f.Close()

    input, err := ioutil.ReadFile("./Fulcrum/" + file)
    failOnError(err, "No se pudo abrir el archivo: " + file)
    lines := strings.Split(string(input), "\n")

    for i, line := range lines {
        if line != ""{

            data := strings.Split(line, " ")
            city := data[1]

            if city == city_name {
                return i
            }
        }
    }

    return -1
}

func UpdateLine(file string, line int, new_data string){
    /*
    Function: Update certain line of a file with new data
    Input:
        file: file name, string. (e.g. example_file.txt)
        line: line number to update, int. (e.g. 3)
        new_data: data to replace line in file, string. (e.g. Tattoine Mos Palgo 100)
    */

    // Open file
    input, err := ioutil.ReadFile("./Fulcrum/" + file)
    failOnError(err, "No se pudo abrir el archivo: " + file)

    // Create new content in specific line
    lines := strings.Split(string(input), "\n")
    lines[line] = new_data

    // Write again updated content
    output := strings.Join(lines, "\n")
    err = ioutil.WriteFile("./Fulcrum/" + file, []byte(output), 0777)
    failOnError(err, "No se pudo escribir sobre el archivo: " + file)
}

func DeleteLine(file string, line int){
    // Open file
    input, err := ioutil.ReadFile("./Fulcrum/" + file)
    failOnError(err, "No se pudo abrir el archivo: " + "./Fulcrum/" + file)

    // Create new content in specific line
    lines := strings.Split(string(input), "\n")

    var new_lines []string
    for i, cur_line := range lines {
        if i != line{
            new_lines = append(new_lines, cur_line)
        }
    }

    // Write again updated content
    output := strings.Join(new_lines, "\n")
    err = ioutil.WriteFile("./Fulcrum/" + file, []byte(output), 0777)
    failOnError(err, "No se pudo escribir sobre el archivo: " + "./Fulcrum/" + file)
}

func WriteLine(file, data string){
    /*
    Function: Append a new line in file with data
    Input:
        file: file name, string. (e.g. example_file.txt)
        data: data to append, string. (e.g. Tattoine Mos Palgo 100)
    */

    // Open "Registro planetario"
    f, err := os.OpenFile("./Fulcrum/" + file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
    failOnError(err, "No se puede abrir el archivo: " + "./Fulcrum/" + file)

    // Write or append new data
    _, err = f.WriteString(data)
    failOnError(err, "No se puede escribir sobre el archivo:" + "./Fulcrum/" + file)

    f.Close()
}

func ReadLine(file string, line int) string{
    /*
    Function: Get a specific line from a file
    Input:
        file: file name, string. (e.g. example_file.txt)
        line: line number to read, int. (e.g. 3)
    */

    // Read file
    input, err := ioutil.ReadFile("./Fulcrum/" + file)
    failOnError(err, "No se pudo abrir el archivo: " + "./Fulcrum/" + file)

    // Return specific line
    lines := strings.Split(string(input), "\n")
    return lines[line]
}

func GetContent(file string) string {
    input, err := ioutil.ReadFile(file)
    failOnError(err, "No se pudo abrir el archivo: " + "./Fulcrum/" + file)

    return string(input)
}

func AppendLog(file, data string){
    WriteLine("log_" + file, data + "\n")
}

func DeleteLogs(){

    files, err := filepath.Glob("./Fulcrum/log_*.txt")
    failOnError(err, "Error al buscar logs.")

    for _, f := range files {
        err = os.Remove(f)
        failOnError(err, "Error al eliminar logs.")
    }
}


func AddCity(command, planet, city, ammount string) int{

    // Check register existence
    index := FindFile(planet)

    // If not, create
    if index == -1 {
        var new_ver []int32
        new_ver = append(new_ver, 0)
        new_ver = append(new_ver, 0)
        new_ver = append(new_ver, 0)

        new_file    := PlanetaryData{
            name: planet,
            version: new_ver,
        }

        server_files = append(server_files, new_file)

    // If exist, check if city exist
    }else{
        index := LineExist(planet, city)

        // If exist, return -1
        if index != -1{
            return -1
        }
    }

    // If city doesnt exist, create
    WriteLine(planet, fmt.Sprintf("%v %v %v\n", planet[:len(planet) - 4], city, ammount))
    return 0
}

func UpdateName(command, planet, city, data string) int{
    index := LineExist(planet, city)
    if index == -1 {
        return -1
    }

    ammount := strings.Split(ReadLine(planet, index), " ")[2]

    UpdateLine(planet, index, fmt.Sprintf("%v %v %v", planet[:len(planet) - 4], data, ammount))
    return 0
}

func UpdateNumber(command, planet, city, data string) int {
    index := LineExist(planet, city)
    if index == -1 {
        return -1
    }

    UpdateLine(planet, index, fmt.Sprintf("%v %v %v\n", planet[:len(planet) - 4], city, data))
    return 0
}

func DeleteCity(command, planet, city string) int {
    index :=  LineExist(planet, city)
    if index == -1 {
        return -1
    }

    DeleteLine(planet, index)
    return 0
}


func MergeMaster(){

    // Block servers while merge
    server_block = true
    fmt.Println("[Merge] Bloqueando servidores.")
    SwitchBlockMaster(servers_ips[1])
    SwitchBlockMaster(servers_ips[2])

    // Ask logs
    var QAddCity []string
    var QUpdateNumber []string
    var QUpdateName []string
    var QDestroyCity []string

    // Connect to first server
    conn, err := grpc.Dial(servers_ips[1], grpc.WithInsecure(), grpc.WithBlock())
    failOnError(err, "Problema al conectar al servidor.")

    defer conn.Close()
    c := pb.NewFulcrumClient(conn)

    stream, err := c.Merge(context.Background(), &pb.MergeReq{
        Code: int32(200),
    })
    failOnError(err, "Error al leer stream de datos durante merge.")

    // For each log
    for {
        server_log, err := stream.Recv()
        if err == io.EOF {
            break
        }
        failOnError(err, "Error durante lectura de stream.")

        // Get commands
        data    := strings.Split(server_log.Command, " ")
        command := data[0]
        file    := data[1] + ".txt"

        if command == "AddCity"{
            QAddCity = append(QAddCity, server_log.Command)
        }else if command == "UpdateNumber"{
            QUpdateNumber = append(QUpdateNumber, server_log.Command)
        }else if command == "UpdateName"{
            QUpdateName = append(QUpdateName, server_log.Command)
        }else if(command == "DestroyCity"){
            QDestroyCity = append(QDestroyCity, server_log.Command)
        }else{
            continue
        }

        // Update clock
        index := FindFile(file)
        if(index == -1){

            new_file    := PlanetaryData{
                name: file,
                version: server_log.Version,
            }

            server_files = append(server_files, new_file)
        }else{
            server_files[index].version[server_log.Server] = server_log.Version[server_log.Server]
        }
    }

    fmt.Println("[Merge] Logs server 1 leidos.")

    // Connect to second server
    conn, err = grpc.Dial(servers_ips[2], grpc.WithInsecure(), grpc.WithBlock())
    failOnError(err, "Problema al conectar al servidor.")

    defer conn.Close()
    c = pb.NewFulcrumClient(conn)

    stream, err = c.Merge(context.Background(), &pb.MergeReq{
        Code: int32(200),
    })
    failOnError(err, "Error al leer stream de datos durante merge.")

    // For each log
    for {
        server_log, err := stream.Recv()
        if err == io.EOF {
            break
        }
        failOnError(err, "Error durante lectura de stream.")

        // Get commands
        data    := strings.Split(server_log.Command, " ")
        command := data[0]
        file    := data[1] + ".txt"

        if command == "AddCity"{
            QAddCity = append(QAddCity, server_log.Command)
        }else if command == "UpdateNumber"{
            QUpdateNumber = append(QUpdateNumber, server_log.Command)
        }else if command == "UpdateName"{
            QUpdateName = append(QUpdateName, server_log.Command)
        }else if(command == "DestroyCity"){
            QDestroyCity = append(QDestroyCity, server_log.Command)
        }else{
            continue
        }

        // Update clock
        index := FindFile(file)
        if index == -1 {

            new_file    := PlanetaryData{
                name: file,
                version: server_log.Version,
            }

            server_files = append(server_files, new_file)
        }else{
            server_files[index].version[server_log.Server] = server_log.Version[server_log.Server]
        }
    }

    fmt.Println("[Merge] Logs server 2 leidos.")

    fmt.Printf("[Merge] CreateLogs:%v, UpdateLogs:%v, DelteLogs:%v.\n", len(QAddCity), len(QUpdateName) + len(QUpdateNumber), len(QDestroyCity))

    // Execute
    fmt.Println("[Merge] Ejecutando logs externos.")
    Executer(QAddCity)
    QUpdateNumber2 := Executer(QUpdateNumber)
    Executer(QUpdateName)
    Executer(QDestroyCity)
    Executer(QUpdateNumber2)

    // Delete logs
    fmt.Println("[Merge] Eliminando logs locales.")
    DeleteLogs()

    // Propagate files
    fmt.Println("[Merge] Propagando archivos a servidores.")
    files, err := filepath.Glob("./Fulcrum/*.txt")
    failOnError(err, "Error al buscar archivos.")
    fmt.Printf("[Merge] Archivos a propagar:%v.\n", len(files))

    Propagate(files, 1)
    Propagate(files, 2)

    // Release servers
    fmt.Println("[Merge] Liberando servidores.")
    SwitchBlockMaster(servers_ips[1])
    SwitchBlockMaster(servers_ips[2])
    server_block = false
}

func SwitchBlockMaster(address string){

    // Connect to Fulcrum
    conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
    failOnError(err, "Problema al conectar al servidor.")

    defer conn.Close()
    c := pb.NewFulcrumClient(conn)

    ctx := context.Background()
    _, err = c.SwitchBlock(ctx, &pb.BlockReq{
        Code: int32(200),
    })

    failOnError(err, "No se pudede acceder al servicio")
}

func Blocked(){
    /*
    Function: Wait until server finish merge
    */

    // if server is blocked
    if server_block{
        fmt.Println("[Informantes] Esperando proceso de merge...")
        // Wait
        for {
            if !server_block{
                break
            }
        }
        fmt.Println("[Informantes] Proceso merge finalizado.")
    }
}

func Propagate(files []string, server int){

    conn, err := grpc.Dial(servers_ips[server], grpc.WithInsecure(), grpc.WithBlock())
    failOnError(err, "Problema al conectar al servidor.")

    defer conn.Close()
    c := pb.NewFulcrumClient(conn)

    stream, err := c.File(context.Background())
    failOnError(err, "Error enviando archivos")

    for _, file := range files {

        new_file := &pb.FileSend{
            Name: file,
            File: GetContent(file),
    	}

        err = stream.Send(new_file)
        failOnError(err, "Error durante stream.")
    }

    _, err = stream.CloseAndRecv()
    failOnError(err, "Error al cerrar stream.")
}



// ================================= Server ================================= //

func (s *server) CRUD(ctx context.Context, req *pb.Command) (*pb.Data, error) {

    fmt.Printf("[Crud Request] Command: %v.\n", req.Command)
    Blocked()

    // Get command data
    data := strings.Split(req.Command, " ")
    comm := data[0]
    file := data[1] + ".txt"
    city := data[2]

    var status int

    if (comm == "AddCity"){
        if len(data) == 4 {
            status = AddCity(req.Command, file, city, data[3])
        }else{
            status = AddCity(req.Command, file, city, "0")
        }
    }else if(comm == "UpdateName"){
        status = UpdateName(req.Command, file, city, data[3])
    }else if(comm == "UpdateNumber"){
        status = UpdateNumber(req.Command, file, city, data[3])
    }else{
        status = DeleteCity(req.Command, file, city)
    }

    if status == -1 {
        fmt.Printf("[Crud Response] Error de consistencia. Sugiere probar con otro servidor, o en 2 minutos.\n")
    	return &pb.Data{
            Code: int32(-1),
    	}, nil
    }

    register := UpdateClock(file)
    AppendLog(file, req.Command)

    fmt.Printf("[Crud Response] Server: %v - %v.\n", server_index, register.version)
	return &pb.Data{
        Code: int32(0),
		Server: server_index,
        Version: register.version,
	}, nil
}

func (s *server) GetRebelds(ctx context.Context, req *pb.GetRebeldsReq) (*pb.GetRebeldsResp, error) {

    fmt.Printf("[GetRebelds Request] Planet: %v. City: %v.\n", req.Planet, req.City)
    Blocked()

    // Get command data
    file := req.Planet + ".txt"

    // Check planet exist
    planet := FindFile(file)
    if planet == -1 {

        new_version := []int32{int32(0),int32(0),int32(0),}
        fmt.Printf("[GetRebelds Response] Ammount: %v. Version: %v.\n", 0, new_version)

        return &pb.GetRebeldsResp{
            Ammount: int32(0),
            Server: int32(-1),
            Version: new_version,
    	}, nil
    }

    // Check city exist
    line := LineExist(file, req.City)
    fmt.Printf("[DEBUG] %v.\n", line)

    if line == -1 {

        new_version := []int32{int32(0),int32(0),int32(0),}
        fmt.Printf("[GetRebelds Response] Ammount: %v. Version: %v.\n", 0, new_version)

        return &pb.GetRebeldsResp{
            Ammount: int32(0),
            Server: int32(-1),
            Version: new_version,
    	}, nil
    }

    // Get data
    response, _ := strconv.Atoi(strings.Split(ReadLine(file, line), " ")[2])
    fmt.Printf("[GetRebelds Response] Ammount: %v. Version: %v.\n", response, server_files[planet].version)

	// Send response
	return &pb.GetRebeldsResp{
		Ammount: int32(response),
        Server: server_index,
        Version: server_files[planet].version,
	}, nil
}

func (s *server) Merge(req *pb.MergeReq, stream pb.Fulcrum_MergeServer) error {

    fmt.Println("[Merge] Iniciando envio de logs.")
    // Send logs
    logs, err := filepath.Glob("./Fulcrum/log_*.txt")
    failOnError(err, "Error al buscar archivos.")

    for _, log := range logs {

        // Open log
        input, err := ioutil.ReadFile(log)
        failOnError(err, "No se pudo abrir el archivo: " + log)

        // Create new content in specific line
        lines := strings.Split(string(input), "\n")

        for _, line := range lines {

            if line != ""{
                file := strings.Split(line, " ")[1]
                version := server_files[FindFile(file+".txt")].version

                new_log := &pb.MergeResp{
                    Command: line,
                    Server: server_index,
                    Version: version,
                }

                err := stream.Send(new_log)
                failOnError(err, "Error al enviar log line")
            }
        }
    }

    // Reset data
    fmt.Println("[Merge] Eliminando archivos obsoletos.")
    files, err := filepath.Glob("./Fulcrum/*.txt")
    failOnError(err, "Error al buscar archivos.")

    for _, f := range files {
        err = os.Remove(f)
        failOnError(err, "Error al eliminar logs.")
    }

    return nil
}

func (s *server) SwitchBlock(ctx context.Context, req *pb.BlockReq) (*pb.BlockResp, error) {

    if server_block{
        server_block = false
        fmt.Println("[Merge] Liberando servidor.")
    }else{
        server_block = true
        fmt.Println("[Merge] Bloqueando servidor.")
    }

	return &pb.BlockResp{
		Code: int32(200),
	}, nil
}

func (s *server) File(stream pb.Fulcrum_FileServer) error {

    fmt.Println("[Merge] Recibiendo datos de propagacion.")
    for {
        file, err := stream.Recv()
        if err == io.EOF {
          return stream.SendAndClose(&pb.FileResp{
            Code: int32(200),
          })
        }
        failOnError(err, "Error recibiendo archivo.")

        err = ioutil.WriteFile(file.Name, []byte(file.File), 0644)
        failOnError(err, "Error escribiendo archivo.")
    }

    fmt.Println("[Merge] Replica realizada.")

    return nil
}

func StartServer(port string){
    /*
    Function: Start a server with Fulcrum GRPC service.
    */

	// Set server
	lis, err := net.Listen("tcp", port)
    failOnError(err, "Error al crear listener.")

	s := grpc.NewServer()
	pb.RegisterFulcrumServer(s, &server{})

	log.Printf("Servidor escuchando en %v\n\n", lis.Addr())
	err = s.Serve(lis)
    failOnError(err, "Fallo al montar servidor.")
}



// ========================================================================== //
func main(){

    // Get params
    port              := flag.String("p", ":50050", "Port to listen")
    fulcrum_1_address := flag.String("sm", "localhost:50050", "Address of first server")
    fulcrum_2_address := flag.String("s1", "localhost:50051", "Address of second server")
    fulcrum_3_address := flag.String("s2", "localhost:50052", "Address of third server")
    index             := flag.Int("i", 0, "Server index")

    flag.Parse()

    // Set ip address
    servers_ips[0] = *fulcrum_1_address
    servers_ips[1] = *fulcrum_2_address
    servers_ips[2] = *fulcrum_3_address

    // Set server index
    server_index = int32(*index)

    // Set seed
    rand.Seed(time.Now().UnixNano())

    // Set server vars
    server_block = false

    // Merge every 'merge_time' seconds from master-node
    if server_index == 0 {
        go StartServer(*port)
        for range time.Tick(time.Second * 120) {
            MergeMaster()
        }
    }else{
        StartServer(*port)
    }
}
