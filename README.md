# Tarea 3 Sistemas Distribuidos 2021-2

## Integrantes

1. Tamara Letelier     201773564-3
1. Francisca Rubio     201773555-4
1. Bruno Vega          201854051-k

##Supuestos
- Si un planeta o ciudad no existe, se asumen 0 rebeldes y Version [0,0,0]. De forma que cualquier update de parte de los informantes, genere consistencia al proceso Leia.

- Se asumen el index del server fulcrum desde 0. (Server0, Server1, Server2)

- Dadas las libertades que otorga el enunciado, se hicieron en Leia y en Informantes de manera distinta algunas cosas:
    1. En Leia, Broker siempre selecciona un server aleatorio. Si la respuesta de este posee un reloj menos actualizado que el solicitado, busca respuesta en otro servidor. De esta manera, Leia solamente realiza una consulta que siempre le retornara un valor igual o mas actualizado.
    1. En Informantes, Broker siempre selecciona un server aleatorio. Si el Informante al conectarse a ese servidor, no encuentra un registro suficientemente actualizado a sus writes, espera a que ocurra un proceso de merge. Se podria re intentar la conexion y realizar un proceso similar a Leia, pero queriamos cubrir diferentes formas de resolver el problema.

##Merge

Se busca una manera de reducir la informacion perdida en cada merge. Si bien se busca consistencia eventual, tambien se reconoce la importancia de realizar un mejor esfuerzo al solucionar errores de consistencia. Entonces, en afan de una secuencia logica de pasos y entendiendo que las secuencias sin sentido: DestroyCity X Y antes de AddCity X Y no son aceptadas y por ende no quedan en logs, se asume que todos los logs poseen llamados que debieron tener sentido al estado del sistema.

Por esto, se propone lo siguiente:

- Se asume que en el master-node ya se aplicaron sus cambios locales. Se solicitan los logs de los slave-nodes via GRPC.
- Se separan los logs en comandos de AddCity, comandos de UpdateNumber, comandos de UpdateName y comandos de DestroyCity.
- Se ejecutan sobre el master-node todos los comandos de AddCity,
- Se ejecutan sobre el master-node los comandos de UpdateNumber: En este paso, algunos comandos no se realizan pues la ciudad que actualizan solamente existe una vez se actualiza su nombre; Estos comandos son guardados para una futura ejecucion.
- Se ejecutan sobre el master-node todos los comandos de UpdateName.
- Se ejecutan sobre el master-node los restantes comandos de UpdateNumer.
- Se ejecutan sobre el master-node todos los comandos DestroyCity.
- Se eliminan los logs tanto en master-node como en los slave-nodes. Tambien se borran los archivos en los slave-nodes.
- Se envian los archivos locales para mantener consistencia en slave-nodes.

##Instrucciones

Las instrucciones de ejecucion son similares en cada maquina. La diferencia radica en el nombre del comando: Para compilar y ejecutar un servidor Fulcrum en su maquina: 'make startF' para iniciar un Informante: 'make startI'. Los codigos para cada ejecucion se dejan a continuacion.

Se recuerda se debe tener el fireware desactivado para que los procesos en maquinas distintas puedan comunicarse.

### Codigos

I - Informantes
L - Leia
B - Broker
F - Fulcrum

### Pasos para compilar y ejecutar

1. make start{CODIGO}

### Pasos para ejecutar sin compilar

1. make dev{CODIGO}
