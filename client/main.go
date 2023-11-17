package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"fmt"
	"net"
	"os"
	"strconv"
)

const(
	tcpProtocol = "tcp4";
	keySize = 1024;
	readWriterSize = keySize / 8;
)

func checkErr(err error){
	if err != nil{
		fmt.Println("error: ", err);
		os.Exit(1);
	}
}

var connectAddr = &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0};

func connectTo() *net.TCPConn{
	fmt.Print("enter port: ");
	fmt.Scanf("%d", &connectAddr.Port);
	fmt.Println("connect to: ", connectAddr);

	c, err := net.DialTCP(tcpProtocol, nil, connectAddr);
	checkErr(err);
	return c;
}

func sendKey(c *net.TCPConn, K *rsa.PrivateKey){
	c.Write([]byte("CONNECT\n"));
	c.Write([]byte(K.PublicKey.N.String() + "\n"));
	c.Write([]byte(strconv.Itoa(K.PublicKey.E) + "\n"));
}

func getBytes(buf *bufio.Reader, n int) []byte{
	bytes, err := buf.Peek(n);
	checkErr(err);
	skipBytes(buf, n);
	return bytes;
}

func skipBytes(buf *bufio.Reader, skipCount int){
	for i := 0; i < skipCount; i++{
		buf.ReadByte();
	}
}

func main(){
	c := connectTo();
	buf := bufio.NewReader(c);
	k, err := rsa.GenerateKey(rand.Reader, keySize);
	checkErr(err);
	sendKey(c, k);

	for {
		cryptoMsg := getBytes(buf, readWriterSize);
		msg, err := rsa.DecryptOAEP(sha1.New(), rand.Reader, k, cryptoMsg, nil);
		checkErr(err);
		fmt.Println(string(msg));
	}
}