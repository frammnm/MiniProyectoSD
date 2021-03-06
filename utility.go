package centralsim

import (
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"os"
	"encoding/gob"
)




type MsgI interface {
	// SetFrom(string)
	// GetClock() v.VClock
	// GetFrom() string

}

type MsgEvent struct {
	Value Event
}

type MsgLookAhead struct {
	Value TypeClock
	From string
}

type MsgNull struct {
	Value IndLocalTrans
	From string
}

func init() {
	gob.Register(&MsgEvent{})
	gob.Register(&MsgLookAhead{})
	gob.Register(&MsgNull{})
}



func stringInArray(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func RunCommand(cmd string, conn *ssh.Client) {
	sess, err := conn.NewSession()
	defer conn.Close()
	if err != nil {
		panic(err)
	}
	defer sess.Close()
	sessStdOut, err := sess.StdoutPipe()
	if err != nil {
		panic(err)
	}
	go io.Copy(os.Stdout, sessStdOut)
	sessStderr, err := sess.StderrPipe()
	if err != nil {
		panic(err)
	}
	go io.Copy(os.Stderr, sessStderr)
	err = sess.Run(cmd)
	if err != nil {
		panic(err)
	}
}

func PublicKey(path string) ssh.AuthMethod {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		panic(err)
	}
	return ssh.PublicKeys(signer)
}

func RemoveFromSlice(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func TrueMap(m map[string]bool) bool {
	for _, v := range(m) {
	  if !v {
	    return false
	  }
	}
	return true
}

func min(a, b TypeClock) TypeClock {
    if a < b {
        return a
    }
    return b
}