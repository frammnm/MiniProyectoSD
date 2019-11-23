/*
PROPOSITO:
		- Tipo abstracto para realizar la simulacion de una (sub)RdP.
HISTORIA DE CAMBIOS:
COMENTARIOS:
		- El resultado de una simulacion local sera un slice dinamico de
		componentes, de forma que cada una de ella sera una structura estatica de
		dos enteros, el primero de ellos sera el codigo de la transicion
		disparada y el segundo sera el valor del reloj local para el que se
		disparo.
-----------------------------------------------------------------
*/
package centralsim

import (
	"fmt"
	"time"
	"net"
	"encoding/gob"
)

// TypeClock defines integer size for holding time.
type TypeClock int64

// ResultadoTransition holds fired transition id and time of firing
type ResultadoTransition struct {
	CodTransition     IndLocalTrans
	ValorRelojDisparo TypeClock
}

// SimulationEngine is the basic data type for simulation execution
type SimulationEngine struct {
	il_mislefs    	Lefs                  // Estructura de datos del simulador
	ii_relojlocal 	TypeClock             // Valor de mi reloj local
	iv_results    	[]ResultadoTransition // slice dinamico con los resultados
	se_addr 				string 								// Direccion de escucha de mensajes
	se_port 				string 								// Puerto de escucha de mensaje
	se_lookout_done bool 									// Booleano para indicar que llegaron todos los mensajes LookAhead
}

/*
-----------------------------------------------------------------
   METODO: NewMotorSimulation
   RECIBE: EStructura datos Lefs
   DEVUELVE: Nada
   PROPOSITO: Construir que recibe la estructura de datos con la que
	   simular, inicializa variables...
   HISTORIA DE CAMBIOS:
COMENTARIOS:
-----------------------------------------------------------------
*/
func MakeMotorSimulation(alLaLef Lefs) SimulationEngine {
	m := SimulationEngine{}
	m.il_mislefs = alLaLef
	return m
}

/*
-----------------------------------------------------------------
   METODO: disparar_transiciones_sensibilizadas
   RECIBE: Valor del reloj local
   DEVUELVE: Nada
   PROPOSITO: Accede a la lista de transiciones sensibilizadas y procede con su
	   disparo, lo que generara nuevos eventos y modificara el marcado de la
		transicion disparada. Igualmente anotara en los resultados el disparo de
		cada transicion para el reloj actual dado
   HISTORIA DE CAMBIOS:
COMENTARIOS:
-----------------------------------------------------------------
*/
func (self *SimulationEngine) fireEnabledTransitions(aiLocalClock TypeClock) {
	for self.il_mislefs.hay_sensibilizadas() { //while
		liCodTrans := self.il_mislefs.get_sensibilizada()
		// fmt.Println("obtuve sensibilizada:", liCodTrans)
		self.il_mislefs.disparar(liCodTrans)

		// Anotar el Resultado que disparo la liCodTrans en tiempoaiLocalClock
		self.iv_results = append(self.iv_results,
			ResultadoTransition{liCodTrans, aiLocalClock})
	}
}

/*
-----------------------------------------------------------------
   METODO: tratar_eventos
   RECIBE: Tiempo para el que trataremos los eventos
   DEVUELVE: Nada
   PROPOSITO: Accede a la lista de eventos y trata todos aquellos con tiempo
	   igual al recibido. Al tratar los eventos se modificaran los valores de
		las funciones de sensibilizacion de algunas transiciones, por lo que puede
		que tengamos nuevas transiciones sensibilizadas.
   HISTORIA DE CAMBIOS:
COMENTARIOS:
-----------------------------------------------------------------
*/
func (self *SimulationEngine) tratar_eventos(ai_tiempo TypeClock) {
	var le_evento Event

	for self.il_mislefs.hay_eventos(ai_tiempo) {
		le_evento = self.il_mislefs.get_primer_evento()

		// Si el valor de la transicion es negativo,indica que pertenece
		// a otra subred y el codigo global de la transicion es pasarlo
		// a positivo y restarle 1
		// ej: -3 -> transicion -(-3) -1 = 2
		// fmt.Println("le_evento.Ii_transicion:", le_evento.Ii_transicion)
		if le_evento.Ii_transicion >= 0 {
			// Establecer nuevo valor de la funcion
			self.il_mislefs.updateFuncValue(le_evento.Ii_transicion,
				le_evento.Ii_cte)
			// Establecer nuevo valor del tiempo
			self.il_mislefs.actualiza_tiempo(le_evento.Ii_transicion,
				le_evento.Ii_tiempo)
		} else {
			fmt.Println("Entre en un evento con transicion remota")
			// Transformar a indice global y enviar mensajea subred remota
			le_evento.Ii_transicion = (-1) * le_evento.Ii_transicion 
			// Buscar direccion de la subred dado indice de transicion global
			addr := self.il_mislefs.Il_pos[le_evento.Ii_transicion]
			var msg MsgI
			msg = MsgEvent{le_evento}
			self.send_message(msg, addr)
		}
	}
}


//Enviar mensaje a traves de la red de forma codificada
func (self *SimulationEngine) send_message(msg MsgI, addr string) {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			fmt.Println("Dial error, retrying..:", err.Error())
		} else {
			// Encode and send data
			encoder := gob.NewEncoder(conn)
			err = encoder.Encode(&msg)
			fmt.Println("Mensaje: ", msg," enviado a: ", addr)
			// Close connection
			conn.Close()
		}
}
/*
-----------------------------------------------------------------
   METODO: esperar_agentes
   RECIBE: Nada
   DEVUELVE: Nada
   PROPOSITO: Espera a que lleguen todos los agentes que hemos enviado
	   anteriormente, para recibir nuevos eventos o el mensaje "No voy
		a generar nada hasta T"
   HISTORIA DE CAMBIOS:
COMENTARIOS:
-----------------------------------------------------------------
*/
func (self *SimulationEngine) esperar_agentes() {
	fmt.Println("Enviando agentes en busqueda de LookAheads..")
	// Buscar direccion de las subredes
	addrs := self.il_mislefs.Il_pre
	self.se_lookout_done = false
	for idGlobal, addr := range(addrs) {
		var msg MsgI
		msg = MsgNull{idGlobal, self.se_addr}
		//Inicializo el mapa de lookouts para esa direccion
		self.il_mislefs.Il_lookOuts[addr] = -1
		go self.send_message(msg, addr)
	}
	self.esperar_lookaheads()
	return
}

func (self *SimulationEngine) esperar_lookaheads() {

	time.Sleep(1000 * time.Millisecond)

	for !self.se_lookout_done {}

	fmt.Println("Recibidos todos los lookAhead")
	return
}

/*
-----------------------------------------------------------------
   METODO: avanzar_tiempo
   RECIBE: Nada
   DEVUELVE: Nada
   PROPOSITO: Modifica el reloj local con el minimo tiempo de entre los
	   recibidos por los agentes o del primer evento encolado en la lista
		de eventos
   HISTORIA DE CAMBIOS:
COMENTARIOS:
-----------------------------------------------------------------
*/
func (self *SimulationEngine) avanzar_tiempo() TypeClock {

	nextTime := self.il_mislefs.tiempo_primer_evento()

	// Operar tiempo basado en self.il_mislefs.Il_lookOuts vs self.il_mislefs.tiempo_primer_evento()
	var currentLookAhead TypeClock
	var nextLookAhead TypeClock
	//Valores de max int 32
	currentLookAhead = 2147483647
	nextLookAhead = 2147483647
	for _, addr := range (self.il_mislefs.Il_pre) {
		nextLookAhead = self.il_mislefs.Il_lookOuts[addr]
		currentLookAhead = min(currentLookAhead, nextLookAhead)
	}

	// fmt.Printf("\nSE OPERO EL TIEMPO CON NEXTTIME:%v   currentLookAhead:%v  lookouts: %v \n\n\n", nextTime, currentLookAhead, self.il_mislefs.Il_lookOuts)

	//Reinicializo el mapa de lookouts despues de calcular
	self.il_mislefs.Il_lookOuts = make(map[string]TypeClock)

	if nextTime <= currentLookAhead && nextTime != -1 {
		fmt.Println("NEXT CLOCK...... : ", nextTime)
		return nextTime
	} else {
		fmt.Println("NEXT CLOCK...... : ", currentLookAhead)
		return currentLookAhead
	}
}

/*
-----------------------------------------------------------------
   METODO: devolver_resultados
   RECIBE: Nada
   DEVUELVE: Nada
   PROPOSITO: Mostrar los resultados de la simulacion
   HISTORIA DE CAMBIOS:

COMENTARIOS:
-----------------------------------------------------------------
*/
func (self SimulationEngine) devolver_resultados() string {
	resultados := "----------------------------------------\n"
	resultados += "Resultados del simulador local\n"
	resultados += "----------------------------------------\n"
	if len(self.iv_results) == 0 {
		resultados += "No esperes ningun resultado...\n"
	}

	for _, li_result := range self.iv_results {
		resultados +=
			"TIEMPO: " + fmt.Sprintf("%v", li_result.ValorRelojDisparo) +
				" -> TRANSICION: " + fmt.Sprintf("%v", li_result.CodTransition) + "\n"
	}

	fmt.Println(resultados)
	return resultados
}

/*
-----------------------------------------------------------------
   METODO: simular
   RECIBE: Ciclo con el que partimos (por si el marcado recibido no
				se corresponde al inicial sino a uno obtenido tras simular
				ai_cicloinicial ciclos)
			Ciclo con el que terminamos
   DEVUELVE: Nada
   PROPOSITO: Simular una RdP
   HISTORIA DE CAMBIOS:
COMENTARIOS:
-----------------------------------------------------------------
*/
func (self *SimulationEngine) Simular(ai_cicloinicial, ai_nciclos TypeClock) {

	//Iniciamos escucha de mensajes
	go self.listen_subnets()
	time.Sleep(time.Duration(10) * time.Second)
	
	ld_ini := time.Now()
	// Inicializamos el reloj local
	// ------------------------------------------------------------------
	self.ii_relojlocal = ai_cicloinicial

	// Inicializamos las transiciones sensibilizadas, es decir, ver si con el
	// marcado inicial tenemos transiciones sensibilizadas
	// ------------------------------------------------------------------
	self.il_mislefs.actualiza_sensibilizadas(self.ii_relojlocal)

	for self.ii_relojlocal <= ai_nciclos {
		// self.il_mislefs.Imprime() //DEPURACION
		fmt.Println("RELOJ LOCAL !!!  = ", self.ii_relojlocal)

		// Si existen transiciones sensibilizadas para reloj local las disparamos
		// ------------------------------------------------------------------
		if self.il_mislefs.hay_sensibilizadas() {
			self.fireEnabledTransitions(self.ii_relojlocal)
		}

		//self.il_mislefs.il_eventos.Imprime()

		// Si existen eventos para el reloj local los tratamos
		// ------------------------------------------------------------------
		if self.il_mislefs.hay_eventos(self.ii_relojlocal) {
			self.tratar_eventos(self.ii_relojlocal)
		}

		// Los nuevos eventos han podido sensibilizar nuevas transiciones
		// ------------------------------------------------------------------
		self.il_mislefs.actualiza_sensibilizadas(self.ii_relojlocal)

		// Tras tratar todos los eventos, si no nos quedan transiciones
		// sensibilizadas no podemos simular nada mas, luego esperamos a
		// los agentes y si no nos generan nuevos eventos procedemos a avanzar
		// el reloj local
		// ------------------------------------------------------------------
		if !self.il_mislefs.hay_sensibilizadas() {
			self.esperar_agentes()
			if !self.il_mislefs.hay_eventos(self.ii_relojlocal) {
				self.ii_relojlocal = self.avanzar_tiempo()

				if self.ii_relojlocal == -1 {
					self.ii_relojlocal = ai_nciclos + 1
				}
			}
		}
	}

	elapsedTime := time.Since(ld_ini)

	// Devolver los resultados de la simulacion
	self.devolver_resultados()
	result := "\n---------------------"
	result += "\nNUMERO DE TRANSICIONES DISPARADAS " +
		fmt.Sprintf("%d", len(self.iv_results)) + "\n"
	result += "TIEMPO SIMULADO en ciclos: " +
		fmt.Sprintf("%d", ai_nciclos-ai_cicloinicial) + "\n"
	result += "COSTE REAL SIMULACION: " +
		fmt.Sprintf("%v", elapsedTime.String()) + "\n"
	fmt.Println(result)
}


func (self *SimulationEngine) listen_subnets() {
	
	//  Preparing to receive conncection
	ln, err := net.Listen("tcp", ":" + self.se_port)

	if err != nil {
		fmt.Println("Error al escuchar en el puerto: ", err.Error())
	} else {
		fmt.Println("***************** Escuchando en la direccion: ", self.se_addr," *****************")
	}

	for {
		
		// Accept incoming connection
		conn, err := ln.Accept()
		// fmt.Printf("\nConexion aceptada desde: [%v]  ---->  en: %v\n", conn.RemoteAddr().String(), self.se_addr)

		if err != nil {
			fmt.Println("Error aceptando conexion:", err.Error())
		}
		
		// Decode data
		var msg MsgI
		decoder := gob.NewDecoder(conn)
		err = decoder.Decode(&msg)

		if err != nil {
			fmt.Printf("Error decodificando mensaje desde: [%s]\n", conn.RemoteAddr().String())
			fmt.Println(err.Error())
		} else {			
			switch val := msg.(type) {
				case *MsgEvent:
					idGlobal := val.Value.Ii_transicion
					for _, t := range(self.il_mislefs.Subnet) {
						if t.IdGlobal == idGlobal {
							val.Value.Ii_transicion = t.IdLocal
						}
					}
					fmt.Printf("\nMensaje de evento agregado a la cola: **%v**\n", val)
					self.il_mislefs.agnade_evento(val.Value)
				case *MsgLookAhead:
					// fmt.Printf("\nMensaje LookAhead recibido y agregado al mapa: **%v**\n", val)
					self.il_mislefs.Il_lookOuts[val.From] = val.Value
					
					done := true
					for _, lookahead := range(self.il_mislefs.Il_lookOuts) {
						if lookahead < 0 {
							done = false 
						}
					}

					if done {
						// fmt.Println("DICCIONARIO DE MIS LOOKAHEADS COMPLETOS: ", self.il_mislefs.Il_lookOuts)
						self.se_lookout_done = true
					}
				case *MsgNull:
					fmt.Printf("\nMensaje Null recibido: **%v**\n", val)
					//Enviar LookAhead
					value := self.calculate_la(val.Value)
					var msg MsgI
					msg = MsgLookAhead{value, self.se_addr}
					self.send_message(msg, val.From)
					// fmt.Printf("\nMensaje LookAhead enviado: %v\n", msg)
				default:
					fmt.Printf("No se que tipo es %T!\n", val)
			}
		}

		// Shut down the connection.
		conn.Close()

	}
}

func (self *SimulationEngine) calculate_la(idGlobal IndLocalTrans) TypeClock {
	var durTime TypeClock  
	for _, t := range(self.il_mislefs.Subnet) {
		if t.IdGlobal == idGlobal {
			durTime = t.Ii_duracion_disparo
			break
		}
	}
	return self.ii_relojlocal + durTime
}
