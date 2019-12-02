package centralsim

import (
	//"log"
	"testing"
	"golang.org/x/crypto/ssh"
	"fmt"
	"strings"
)


// var defaultAddresses = []string{"155.210.154.199:17431", "155.210.154.200:17432", "155.210.154.200:17433", "155.210.154.204:17434", "155.210.154.208:17435"}
// var defaultAddresses = []string{"192.168.1.70:17431", "192.168.1.70:17432", "192.168.1.70:17433", "192.168.1.70:17434", "192.168.1.70:17435"}
// var defaultAddresses = []string{"10.1.24.55:17431", "10.1.24.55:17432", "10.1.24.55:17433", "10.1.24.55:17434", "10.1.24.55:17435"}
var defaultAddresses = []string{"127.0.0.1:17431", "127.0.0.1:17432", "127.0.0.1:17433", "127.0.0.1:17434", "127.0.0.1:17435"}
var defaultPorts = []string{"17431", "17432", "17433", "17434", "17435"}
var simulationTime = TypeClock(12)


func TestSimulationEngineBasic(t *testing.T) {
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

func TestSimulationEnginePartition1(t *testing.T) {
	//t.Skip("skipping test simulation.")
	lfs := Lefs{ //Ejemplo PN documento adjunto
		Subnet: TransitionList{
			Transition{
				IdLocal:             0,
				IdGlobal:						 0,
				IiValorLef:          0,
				Ii_duracion_disparo: 1,
				Ii_listactes: []TransitionConstant{
					TransitionConstant{0, 1},
					TransitionConstant{-1, -1},
					TransitionConstant{-5, -1},
					TransitionConstant{-9, -1},
					TransitionConstant{-13, -1},

				},
			},
			Transition{
				IdLocal:             1,
				IdGlobal:						 17,
				IiValorLef:          4,
				Ii_duracion_disparo: 1,
				Ii_listactes: []TransitionConstant{
					TransitionConstant{1, 4},
					TransitionConstant{0, -1},
				},
			},
		},
		Il_pos: map[IndLocalTrans]string{
			1: defaultAddresses[1],
			5: defaultAddresses[2],
			9: defaultAddresses[3],
			13: defaultAddresses[4],
		},
		Il_pre: map[IndLocalTrans]string{
			4: defaultAddresses[1],
			8: defaultAddresses[2],
			12: defaultAddresses[3],
			16: defaultAddresses[4],
		},
		Il_lookOuts: make(map[string]TypeClock), 
	}
	ms := MakeMotorSimulation(lfs)
	ms.se_addr = defaultAddresses[0]
	ms.se_port = defaultPorts[0]
	ms.Simular(0, simulationTime) // ciclo 0 hasta ciclo 3
	
}


func TestSimulationEnginePartition2(t *testing.T) {
	//t.Skip("skipping test simulation.")
	lfs := Lefs{ //Ejemplo PN documento adjunto
		Subnet: TransitionList{
			Transition{
				IdLocal:             0,
				IdGlobal:						 1,
				IiValorLef:          1,
				Ii_duracion_disparo: 1,
				Ii_listactes: []TransitionConstant{
					TransitionConstant{0, 1},
					TransitionConstant{1, -1},
				},
			},
			Transition{
				IdLocal:             1,
				IdGlobal:						 2,
				IiValorLef:          1,
				Ii_duracion_disparo: 1,
				Ii_listactes: []TransitionConstant{
					TransitionConstant{1, 1},
					TransitionConstant{2, -1},
				},
			},
			Transition{
				IdLocal:             2,
				IdGlobal:						 3,
				IiValorLef:          1,
				Ii_duracion_disparo: 1,
				Ii_listactes: []TransitionConstant{
					TransitionConstant{2, 1},
					TransitionConstant{3, -1},
				},
			},
			Transition{
				IdLocal:             3,
				IdGlobal:						 4,
				IiValorLef:          1,
				Ii_duracion_disparo: 1,
				Ii_listactes: []TransitionConstant{
					TransitionConstant{3, 1},
					TransitionConstant{-17, -1},
				},
			},
		},
		Il_pos: map[IndLocalTrans]string{
			17: defaultAddresses[0],
		},
		Il_pre: map[IndLocalTrans]string{
			0: defaultAddresses[0],
		},
		Il_lookOuts: make(map[string]TypeClock),
	}
	ms := MakeMotorSimulation(lfs)
	ms.se_addr = defaultAddresses[1]
	ms.se_port = defaultPorts[1]
	ms.Simular(0, simulationTime) // ciclo 0 hasta ciclo 3
	
}

func TestSimulationEnginePartition3(t *testing.T) {
	//t.Skip("skipping test simulation.")
	lfs := Lefs{ //Ejemplo PN documento adjunto
		Subnet: TransitionList{ //Ejemplo PN documento adjunto
			Transition{
				IdLocal:             0,
				IdGlobal:						 5,
				IiValorLef:          1,
				Ii_duracion_disparo: 1,
				Ii_listactes: []TransitionConstant{
					TransitionConstant{0, 1},
					TransitionConstant{1, -1},
				},
			},
			Transition{
				IdLocal:             1,
				IdGlobal:						 6,
				IiValorLef:          1,
				Ii_duracion_disparo: 1,
				Ii_listactes: []TransitionConstant{
					TransitionConstant{1, 1},
					TransitionConstant{2, -1},
				},
			},
			Transition{
				IdLocal:             2,
				IdGlobal:						 7,
				IiValorLef:          1,
				Ii_duracion_disparo: 1,
				Ii_listactes: []TransitionConstant{
					TransitionConstant{2, 1},
					TransitionConstant{3, -1},
				},
			},
			Transition{
				IdLocal:             3,
				IdGlobal:						 8,
				IiValorLef:          1,
				Ii_duracion_disparo: 1,
				Ii_listactes: []TransitionConstant{
					TransitionConstant{3, 1},
					TransitionConstant{-17, -1},
				},
			},
		},
		Il_pos: map[IndLocalTrans]string{
			17: defaultAddresses[0],
		},
		Il_pre: map[IndLocalTrans]string{
			0: defaultAddresses[0],
		},
		Il_lookOuts: make(map[string]TypeClock),
	}
	ms := MakeMotorSimulation(lfs)
	ms.se_addr = defaultAddresses[2]
	ms.se_port = defaultPorts[2]
	ms.Simular(0, simulationTime) // ciclo 0 hasta ciclo 3
	
}

func TestSimulationEnginePartition4(t *testing.T) {
	//t.Skip("skipping test simulation.")
	lfs := Lefs{ //Ejemplo PN documento adjunto
		Subnet: TransitionList{ //Ejemplo PN documento adjunto
			Transition{
				IdLocal:             0,
				IdGlobal:						 9,
				IiValorLef:          1,
				Ii_duracion_disparo: 1,
				Ii_listactes: []TransitionConstant{
					TransitionConstant{0, 1},
					TransitionConstant{1, -1},
				},
			},
			Transition{
				IdLocal:             1,
				IdGlobal:						 10,
				IiValorLef:          1,
				Ii_duracion_disparo: 1,
				Ii_listactes: []TransitionConstant{
					TransitionConstant{1, 1},
					TransitionConstant{2, -1},
				},
			},
			Transition{
				IdLocal:             2,
				IdGlobal:						 11,
				IiValorLef:          1,
				Ii_duracion_disparo: 1,
				Ii_listactes: []TransitionConstant{
					TransitionConstant{2, 1},
					TransitionConstant{3, -1},
				},
			},
			Transition{
				IdLocal:             3,
				IdGlobal:						 12,
				IiValorLef:          1,
				Ii_duracion_disparo: 1,
				Ii_listactes: []TransitionConstant{
					TransitionConstant{3, 1},
					TransitionConstant{-17, -1},
				},
			},
		},
		Il_pos: map[IndLocalTrans]string{
			17: defaultAddresses[0],
		},
		Il_pre: map[IndLocalTrans]string{
			0: defaultAddresses[0],
		},
		Il_lookOuts: make(map[string]TypeClock),
	}
	ms := MakeMotorSimulation(lfs)
	ms.se_addr = defaultAddresses[3]
	ms.se_port = defaultPorts[3]
	ms.Simular(0, simulationTime) // ciclo 0 hasta ciclo 3
	
}

func TestSimulationEnginePartition5(t *testing.T) {
	//t.Skip("skipping test simulation.")
	lfs := Lefs{ //Ejemplo PN documento adjunto
		Subnet: TransitionList{ //Ejemplo PN documento adjunto
			Transition{
				IdLocal:             0,
				IdGlobal:						 13,
				IiValorLef:          1,
				Ii_duracion_disparo: 1,
				Ii_listactes: []TransitionConstant{
					TransitionConstant{0, 1},
					TransitionConstant{1, -1},
				},
			},
			Transition{
				IdLocal:             1,
				IdGlobal:						 14,
				IiValorLef:          1,
				Ii_duracion_disparo: 1,
				Ii_listactes: []TransitionConstant{
					TransitionConstant{1, 1},
					TransitionConstant{2, -1},
				},
			},
			Transition{
				IdLocal:             2,
				IdGlobal:						 15,
				IiValorLef:          1,
				Ii_duracion_disparo: 1,
				Ii_listactes: []TransitionConstant{
					TransitionConstant{2, 1},
					TransitionConstant{3, -1},
				},
			},
			Transition{
				IdLocal:             3,
				IdGlobal:						 16,
				IiValorLef:          1,
				Ii_duracion_disparo: 1,
				Ii_listactes: []TransitionConstant{
					TransitionConstant{3, 1},
					TransitionConstant{-17, -1},
				},
			},
		},
		Il_pos: map[IndLocalTrans]string{
			17: defaultAddresses[0],
		},
		Il_pre: map[IndLocalTrans]string{
			0: defaultAddresses[0],
		},
		Il_lookOuts: make(map[string]TypeClock),
	}
	ms := MakeMotorSimulation(lfs)
	ms.se_addr = defaultAddresses[4]
	ms.se_port = defaultPorts[4]
	ms.Simular(0, simulationTime) // ciclo 0 hasta ciclo 3
	
}

func TestSimulationDistrEngineX(t *testing.T) {
	subnets := []string{"TestSimulationEnginePartition1", "TestSimulationEnginePartition2", "TestSimulationEnginePartition3", "TestSimulationEnginePartition4", "TestSimulationEnginePartition5"}
	// dir := "/home/francisco/go/src/MiniProyectoSD"
	dir := "/home/a794893/go/src/MiniProyectoSD"
	// rsa := "/home/francisco/.ssh/id_rsa"
	// rsa := "/home/a794893/.ssh/id_rsa"

	i := 0
	for i < len(defaultAddresses) {
	 config := &ssh.ClientConfig {
	  User: "a794893",
	  // User: "francisco",
	  Auth: []ssh.AuthMethod{ssh.Password("hetero64")}, HostKeyCallback: ssh.InsecureIgnoreHostKey()}
	  // Auth: []ssh.AuthMethod{PublicKey(rsa)}, HostKeyCallback: ssh.InsecureIgnoreHostKey()}

	 // fmt.Println("Ssh command: ", strings.Split(defaultAddresses[i], ":")[0] + ":22", config)
	 conn, err := ssh.Dial("tcp", strings.Split(defaultAddresses[i], ":")[0] + ":22", config)
	 if err != nil {
	 	panic(err)
	 }

	 // Start Snode
	 fmt.Println("ssh to:", defaultAddresses[i], len(defaultAddresses), i)
	 // go RunCommand("cd " + dir + " && go test -run " + subnets[i], conn)
	 // fmt.Println("cd " + dir + " && go test -run " + subnets[i])
	 go RunCommand("cd " + dir + " && /usr/local/go/bin/go test -run " + subnets[i], conn)
	 fmt.Println("cd " + dir + " && /usr/local/go/bin/go test -run " + subnets[i])
	 i++ 
	 fmt.Println("nmero", i)
	}

	for {}
}

/*
func TestTransition(t *testing.T) {

}

func TestLefs(t *testing.T) {

}
*/
