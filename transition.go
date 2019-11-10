package centralsim

import "fmt"

//--------------------------------------------------------------------------
// TransitionList is a list of transitions themselves
type TransitionList []Transition //Slice de transiciones como Lista

// length return length of TransitionList with type adapted to IndLocalTrans
func (self TransitionList) length() IndLocalTrans {
	return IndLocalTrans(len(self))
}

//--------------------------------------------------------------------------

// IndGlobalTrans is a index of a transition in the global list
//type IndGlobalTrans int32

// IndLocalTrans is a index of a transition in the local lefs list
type IndLocalTrans int32

//TypeConst is the constant to propagate in lefs
type TypeConst int32

type TransitionConstant struct {
	INextTrans IndLocalTrans
	Cnstnt     TypeConst
}

//------------------------------------------------------------------------

// -----------------------------------------------------------------------
// Tipo abstracto Transition para guardar la informacion de una transicion
// -----------------------------------------------------------------------
type Transition struct {
	// indice en la tabla global de transiciones
	IdLocal IndLocalTrans

	// iiValorLef es el valor que tiene la funcion de
	// sensibilizacion en el instante de tiempo que nos da
	// la variable ii_tiempo
	IiValorLef TypeConst
	Ii_tiempo  TypeClock

	// tiempo que dura el disparo de la transicion
	Ii_duracion_disparo TypeClock

	// vector de transiciones a las que tengo que propagar cte
	// cuando se dispare esta transicion, junto con la cte que
	// tengo que propagar
	Ii_listactes []TransitionConstant
}

/*
	-----------------------------------------------------------------
	   METODO: Imprime
	   RECIBE: Nada
	   DEVUELVE: Nada
	   PROPOSITO: Imprimir los atributos de la clase para depurar errores
		-----------------------------------------------------------------
*/
func (self *Transition) Imprime() {
	fmt.Println("Dato Transicion:")
	fmt.Println("IDGLOBAL: ", self.IdLocal)
	fmt.Println(" VALOR LEF: ", self.IiValorLef)
	fmt.Println(" TIEMPO: ", self.Ii_tiempo)
	fmt.Println(" DURACION DISPARO: ", self.Ii_duracion_disparo)
	fmt.Println(" LISTA DE CTES: ")
	for _, v := range self.Ii_listactes {
		fmt.Println("\tTRANSICION: ", v.INextTrans, "\t\tCTE: ", v.Cnstnt)
	}
}

/*
	-----------------------------------------------------------------
   METODO: Imprime_valores
   RECIBE: Nada
   DEVUELVE: Nada
   PROPOSITO: Imprimir simplemente el valor de la transicion
	COMENTARIO : es solo de lectura
	-----------------------------------------------------------------
*/
func (self Transition) Imprime_valores() {
	fmt.Println("Transicion -> ")
	//	fmt.Println("\tIDGLOBAL: ", self.Ii_idglobal)
	fmt.Println("\t\tVALOR LEF: ", self.IiValorLef)
	fmt.Println("\t\tTIEMPO: ", self.Ii_tiempo)
}

//----------------------------------------------------------------------

// Stack Transition is a Stack of transitions indices
type StackTransitions []IndLocalTrans

//Push transition id to stack
func (self *StackTransitions) push(i_tr IndLocalTrans) {
	*self = append(*self, i_tr)
}

//Pop transition id from stack
func (self *StackTransitions) pop() IndLocalTrans {
	if (*self).isEmpty() {
		return -1
	} else {
		i_tr := (*self)[len(*self)-1]  // obtener dato de lo alto de la pila
		*self = (*self)[:len(*self)-1] //desempilar
		return i_tr
	}
}

func (self StackTransitions) isEmpty() bool {
	return len(self) == 0
}
