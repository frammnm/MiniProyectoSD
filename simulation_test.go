package centralsim

import (
	//"log"
	"testing"
)

func TestSimulationEngine(t *testing.T) {
	//t.Skip("skipping test simulation.")
	lfs := Lefs{ //Ejemplo PN documento adjunto
		Subnet: TransitionList{
			Transition{
				IdLocal:             0,
				IiValorLef:          0,
				Ii_duracion_disparo: 1,
				Ii_listactes: []TransitionConstant{
					TransitionConstant{0, 1},
					TransitionConstant{1, -1},
					TransitionConstant{2, -1},
				},
			},
			Transition{
				IdLocal:             1,
				IiValorLef:          1,
				Ii_duracion_disparo: 1,
				Ii_listactes: []TransitionConstant{
					TransitionConstant{1, 1},
					TransitionConstant{2, -1},
				},
			},
			Transition{
				IdLocal:             2,
				IiValorLef:          2,
				Ii_duracion_disparo: 1,
				Ii_listactes: []TransitionConstant{
					TransitionConstant{2, 2},
					TransitionConstant{0, -1},
				},
			},
		},
	}
	ms := MakeMotorSimulation(lfs)
	ms.Simular(0, 3) // ciclo 0 hasta ciclo 3
}

/*
func TestTransition(t *testing.T) {

}

func TestLefs(t *testing.T) {

}
*/
