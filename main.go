package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"os"
)

func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}
	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func EstablishConnection(config *ssh.ClientConfig) (*ssh.Session, error) {
	connection, err := ssh.Dial("tcp", "51.79.30.99:22", config)

	if err != nil {
		return nil, err
	}

	session, err := connection.NewSession()

	if err != nil {
		return nil, err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		return nil, fmt.Errorf("request for pseudo terminal failed: %s", err)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("unable to setup stdin for session: %v", err)
	}
	go io.Copy(stdin, os.Stdin)

	stdout, err := session.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("unable to setup stdout for session: %v", err)
	}
	go io.Copy(os.Stdout, stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("unable to setup stderr for session: %v", err)
	}

	go io.Copy(os.Stderr, stderr)

	return session, nil
}

func main() {
	sshConfig := ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			PublicKeyFile("test"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	session, err := EstablishConnection(&sshConfig)
	defer session.Close()

	if err != nil {
		panic(err)
	}

	public, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	fmt.Println(base64.StdEncoding.EncodeToString(public[:]))
	fmt.Println(base64.StdEncoding.EncodeToString(priv[:]))

	session.Run("sh script.sh UwU")
}
