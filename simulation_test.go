package centralsim

import (
	//"log"
	"testing"
	"golang.org/x/crypto/ssh"
	"fmt"
	"strings"
)


// var defaultAddresses = []string{"155.210.154.200:17431", "155.210.154.199:17432", "155.210.154.197:17433"}
var defaultAddresses = []string{"192.168.1.70:17431", "192.168.1.70:17432", "192.168.1.70:17433"}
var defaultPorts = []string{"17431", "17432", "17433"}
var simulationTime = TypeClock(6)


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
					TransitionConstant{-3, -1},
				},
			},
			Transition{
				IdLocal:             1,
				IdGlobal:						 5,
				IiValorLef:          2,
				Ii_duracion_disparo: 1,
				Ii_listactes: []TransitionConstant{
					TransitionConstant{1, 2},
					TransitionConstant{0, -1},
				},
			},
		},
		Il_pos: map[IndLocalTrans]string{
			1: defaultAddresses[1],
			3: defaultAddresses[2],
		},
		Il_pre: map[IndLocalTrans]string{
			2: defaultAddresses[1],
			4: defaultAddresses[2],
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
					TransitionConstant{-5, -1},
				},
			},
		},
		Il_pos: map[IndLocalTrans]string{
			5: defaultAddresses[0],
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
		Subnet: TransitionList{
			Transition{
				IdLocal:             0,
				IdGlobal:						 3,
				IiValorLef:          1,
				Ii_duracion_disparo: 1,
				Ii_listactes: []TransitionConstant{
					TransitionConstant{0, 1},
					TransitionConstant{1, -1},
				},
			},
			Transition{
				IdLocal:             1,
				IdGlobal:						 4,
				IiValorLef:          1,
				Ii_duracion_disparo: 1,
				Ii_listactes: []TransitionConstant{
					TransitionConstant{1, 1},
					TransitionConstant{-5, -1},
				},
			},
		},
		Il_pos: map[IndLocalTrans]string{
			5: defaultAddresses[0],
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

func TestSimulationDistrEngine(t *testing.T) {
	//t.Skip("skipping test simulation.")
	subnets := []string{"TestSimulationEngineBasic"}
	// subnets := []string{"TestSimulationEnginePartition1"}//, "TestSimulationEnginePartition2", "TestSimulationEnginePartition3"}
	defaultAddresses := []string{"127.0.0.1:17431"}
	// defaultAddresses := []string{"155.210.154.197:17433", "155.210.154.200:17431", "155.210.154.199:17432"}
	dir := "/home/francisco/go/src/MiniProyectoSD"
	// dir := "/home/a794893/go/src/MiniProyectoSD"
	rsa := "/home/francisco/.ssh/id_rsa"
	// rsa := "/home/a794893/.ssh/id_rsa"

	i := 0
	for i < len(defaultAddresses) {
	 config := &ssh.ClientConfig {
	  // User: "a794893",
	  User: "francisco",
	  Auth: []ssh.AuthMethod{PublicKey(rsa)}, HostKeyCallback: ssh.InsecureIgnoreHostKey()}

	 conn, err := ssh.Dial("tcp", strings.Split(defaultAddresses[i], ":")[0] + ":22", config)
	 if err != nil {
	 	panic(err)
	 }

	 // Start Snode
	 fmt.Println("ssh to:", defaultAddresses[i], len(defaultAddresses), i)
	 go RunCommand("cd " + dir + " && go test -run " + subnets[i], conn)
	 fmt.Println("cd " + dir + " && go test -run " + subnets[i])
	 // go RunCommand("cd " + dir + " && /usr/local/go/bin/go test -run " + subnets[i], conn)
	 // fmt.Println("cd " + dir + " && /usr/local/go/bin/go test -run " + subnets[i])
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
