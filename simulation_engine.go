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
	il_mislefs    Lefs                  // Estructura de datos del simulador
	ii_relojlocal TypeClock             // Valor de mi reloj local
	iv_results    []ResultadoTransition // slice dinamico con los resultados
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
		if le_evento.Ii_transicion >= 0 {
			// Establecer nuevo valor de la funcion
			self.il_mislefs.updateFuncValue(le_evento.Ii_transicion,
				le_evento.Ii_cte)
			// Establecer nuevo valor del tiempo
			self.il_mislefs.actualiza_tiempo(le_evento.Ii_transicion,
				le_evento.Ii_tiempo)
		} else {
			fmt.Println("Entre en un evento con transicion remota")
		}
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
	fmt.Println("Aun sin agentes")
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
	fmt.Println("NEXT CLOCK...... : ", nextTime)
	return nextTime
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
	ld_ini := time.Now()

	// Inicializamos el reloj local
	// ------------------------------------------------------------------
	self.ii_relojlocal = ai_cicloinicial

	// Inicializamos las transiciones sensibilizadas, es decir, ver si con el
	// marcado inicial tenemos transiciones sensibilizadas
	// ------------------------------------------------------------------
	//self.il_mislefs.actualiza_sensibilizadas(self.ii_relojlocal)

	for self.ii_relojlocal <= ai_nciclos {

				// Los nuevos eventos han podido sensibilizar nuevas transiciones
		// ------------------------------------------------------------------
		self.il_mislefs.actualiza_sensibilizadas(self.ii_relojlocal)

		// Tras tratar todos los eventos, si no nos quedan transiciones
		// sensibilizadas no podemos simular nada mas, luego esperamos a
		// los agentes y si no nos generan nuevos eventos procedemos a avanzar
		// el reloj local
		// ------------------------------------------------------------------
		if !self.il_mislefs.hay_sensibilizadas() {
			// self.esperar_agentes()
			if !self.il_mislefs.hay_eventos(self.ii_relojlocal) {
				self.ii_relojlocal = self.avanzar_tiempo()

				if self.ii_relojlocal == -1 {
					self.ii_relojlocal = ai_nciclos + 1
				}
			}
		}

		self.il_mislefs.Imprime() //DEPURACION
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


	}

	elapsedTime := time.Since(ld_ini)

	// Devolver los resultados de la simulacion
	self.devolver_resultados()
	result := "\n---------------------"
	result += "NUMERO DE TRANSICIONES DISPARADAS " +
		fmt.Sprintf("%d", len(self.iv_results)) + "\n"
	result += "TIEMPO SIMULADO en ciclos: " +
		fmt.Sprintf("%d", ai_nciclos-ai_cicloinicial) + "\n"
	result += "COSTE REAL SIMULACION: " +
		fmt.Sprintf("%v", elapsedTime.String()) + "\n"
	fmt.Println(result)
}
