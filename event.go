/*
PROPOSITO: Guardar la informacion de los eventos necesarios para
   la simulacion
-----------------------------------------------------------------
*/
package centralsim

import "fmt"

type Event struct {
	// Tiempo para el que debemos considerar el evento
	Ii_tiempo TypeClock
	// A que transicion (indice transicion en subred)
	Ii_transicion IndLocalTrans
	// Constante que mandamos
	Ii_cte TypeConst
	// Si el evento debe ejecutarse remotamente o no
	// Ii_remote bool 
}

/*
-----------------------------------------------------------------
   METODO: NewEvento
   RECIBE: Tiempo, transicion y cte del evento a crear
   DEVUELVE: Event
   PROPOSITO: Crear evento con todos los datos del nuevo evento creados
-----------------------------------------------------------------

func NewEvento(ai_tiempo TypeClock, ai_transicion IndLocalTrans, ai_cte TypeConst) *Event {
	e := new(Event)
	set_tiempo(e.ai_tiempo)
	set_transicion(e.ai_transicion)
	set_cte(e.ai_cte)
}
*/

/*
-----------------------------------------------------------------
   METODO: set_tiempo
   RECIBE: Tiempo
   DEVUELVE: Nada
   PROPOSITO: Modificar el tiempo del evento
-----------------------------------------------------------------
*/
func (self *Event) Set_tiempo(ai_tiempo TypeClock) {
	self.Ii_tiempo = ai_tiempo
}

/*
-----------------------------------------------------------------
   METODO: set_transicion
   RECIBE: Identificador de la transicion (indice en array
	   de transiciones de esa subred)
   DEVUELVE: Nada
   PROPOSITO: Modificar la transicion del evento
-----------------------------------------------------------------
*/
func (self *Event) Set_transicion(ai_transicion IndLocalTrans) {
	self.Ii_transicion = ai_transicion
}

/*
-----------------------------------------------------------------
   METODO: set_cte
   RECIBE: Cte a transmitir
   DEVUELVE: Nada
   PROPOSITO: Modificar la cte del evento
-----------------------------------------------------------------
*/
func (self *Event) Set_cte(ai_cte TypeConst) {
	self.Ii_cte = ai_cte
}

/*
-----------------------------------------------------------------
   METODO: get_tiempo
   RECIBE: Nada
   DEVUELVE: El atributo Ii_tiempo
   PROPOSITO: Recoger el tiempo del evento
-----------------------------------------------------------------
*/
func (self Event) get_tiempo() TypeClock {
	return self.Ii_tiempo
}

/*
-----------------------------------------------------------------
   METODO: get_transicion
   RECIBE: Nada
   DEVUELVE: El atributo Ii_transicion
   PROPOSITO: Recoger la transicion del evento
-----------------------------------------------------------------
*/
func (self Event) get_transicion() IndLocalTrans {
	return self.Ii_transicion
}

/*
-----------------------------------------------------------------
   METODO: get_cte
   RECIBE: Nada
   DEVUELVE: El atributo Ii_cte
   PROPOSITO: Recoger la cte del evento
-----------------------------------------------------------------
*/
func (self Event) get_cte() TypeConst {
	return self.Ii_cte
}

/*
-----------------------------------------------------------------
   METODO: Imprime
   RECIBE: Nada
   DEVUELVE: Nada
   PROPOSITO: Visualizar los atributos de un evento para depurar errores
-----------------------------------------------------------------
*/
func (self Event) Imprime(i int) {
	fmt.Println("  EVENTO -> ", i)
	fmt.Println("    Tiempo: ", self.Ii_tiempo)
	fmt.Println("    Transicion: ", self.Ii_transicion)
	fmt.Println("    Constante: ", self.Ii_cte)
}

//----------------------------------------------------------------------------

// EventList es el tipo que almacena la lista de eventos necesaria
// para los motores de	simulacion.
type EventList []Event

/*
-----------------------------------------------------------------
   METODO: make
   RECIBE: Longitud de la lista de eventos
DEVUELVE: Nada
   PROPOSITO: crear la lista de tamaño aiLongitud
-----------------------------------------------------------------
*/
func MakeEventList(aiLongitud int) EventList {
	return make(EventList, aiLongitud)
}

/*
-----------------------------------------------------------------
   METODO: inserta
RECIBE: Evento a insertar
   DEVUELVE: Nada
PROPOSITO: Insertar el evento en la lista de eventos, de forma que
   	la insercion sea ordenada por tiempo.
-----------------------------------------------------------------
*/
func (self *EventList) inserta(ae_evento Event) {
	var i int // INITIALIZED to 0 !!!

	//fmt.Println("Insertar evento en lista : ", ae_evento, *self)

	// Obtengo la posicion ordenada del evento en slice con i
	for _, e := range *self {
		if e.get_tiempo() >= ae_evento.get_tiempo() {
			break
		}
		i++
	}

	//fmt.Println("POSICION a INSERTAR en lista de evnetos : ", i)
	*self = append((*self)[:i], append([]Event{ae_evento}, (*self)[i:]...)...)

	//fmt.Println("DESPUES de insertar : ", *self)
}

/*
 -----------------------------------------------------------------
    METODO: recoge_primer_evento
	RECIBE: Nada
    DEVUELVE: El primer evento de la lista
	PROPOSITO: Recoger el primer evento encolado
 -----------------------------------------------------------------
*/
func (self EventList) recoge_primer_evento() Event {
	if len(self) > 0 {
		return self[0]
	} else {
		return Event{} //sino devuelve el tipo Event, zeroed
	}
}

/*
-----------------------------------------------------------------
   METODO: elimina_primer_evento
RECIBE: Nada
   DEVUELVE: Nada
PROPOSITO: Eliminar el primer evento encolado
-----------------------------------------------------------------
*/
func (self *EventList) elimina_primer_evento() {
	if len(*self) > 0 {
		//suprimir con posibilidad de liberacion de memoria
		copy(*self, (*self)[1:])
		(*self)[len(*self)-1] = Event{} //pongo a zero el previo último Event
		(*self) = (*self)[:len(*self)-1]
	}
}

/*
-----------------------------------------------------------------
   METODO: longitud
RECIBE: Nada
   DEVUELVE: Numero de eventos encolados
PROPOSITO: Conocer el numero de elementos de la lista de eventos
-----------------------------------------------------------------
*/
func (self EventList) longitud() int {
	return len(self)
}

/*
-----------------------------------------------------------------
   METODO: Imprime
RECIBE: Nada
   DEVUELVE: Nada
PROPOSITO: Imprimir la lista de eventos
-----------------------------------------------------------------
*/
func (self EventList) Imprime() {
	fmt.Println("Estructura EventList")
	for i, e := range self {
		e.Imprime(i)
	}
}
