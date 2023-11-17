package main

import (
	"math/big"
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
	readWriterSize = keySize/8;
)

type remoteConn struct {
	c *net.TCPConn;
	pubK *rsa.PublicKey;
}

func checkErr(err error){ 
	if err != nil {
		fmt.Println(err);
		os.Exit(1);
	}
}

var listenAddr = &net.TCPAddr{IP: net.IPv4(127,0,0,1), Port: 0};

func getRemoteConn(c *net.TCPConn) *remoteConn{
	return &remoteConn{c, waitPubKey(bufio.NewReader(c))};
}

func waitPubKey(buf *bufio.Reader) (*rsa.PublicKey) {
	
	// Читаем строку из буфера
	line, _, err := buf.ReadLine();
	checkErr(err);
	
	// Так как тип line - []byte (срез байт)
	// то для удобства сравнения переконвертируем <code><b>line</b></code> в строку
	if string(line) == "CONNECT" {
		
		// Далее мы будем читать буфер в том же порядке, в котором отправляем данные с клиента
		line, _, err := buf.ReadLine(); checkErr(err); // Читаем PublicKey.N

		// Создаём пустой rsa.PublicKey
		pubKey := rsa.PublicKey{N: big.NewInt(0)};
		// pubKey.N == 0 
		// тип pubKey.N big.Int http://golang.org/pkg/big/#Int
		
		// Конвертируем полученную строку и запихиваем в pubKey.N big.Int
		pubKey.N.SetString(string(line), 10);
		// Метод SetString() получает 2 параметра:
		// string(line) - конвертирует полученные байты в строку
		// 10 - система исчисления используемая в данной строке 
		// (2 двоичная, 8 восьмеричная, 10 десятичная, 16 шестнадцатеричная ...)
		
		// Читаем из буфера второе число для pubKey.E
		line, _, err = buf.ReadLine();
		checkErr(err);

		// Используемый пакет strconv для конвертации тип string в тип int
		pubKey.E, err = strconv.Atoi(string(line)); checkErr(err);
		
		// возвращаем ссылку на rsa.PublicKey
		return &pubKey;
		
	} else {
		
		// В этом случае дальнейшее действия программы не предусмотренною. По этому:
		// Выводим что получили
		fmt.Println("Error: unkown command ", string(line));
		os.Exit(1);
	}
	return nil;
}

func (rConn *remoteConn) sendCommand(comm string){
	eComm, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, rConn.pubK, []byte(comm), nil);
	checkErr(err);
	rConn.c.Write(eComm);
}

func listen() {
	i, err := net.ListenTCP(tcpProtocol, listenAddr);
	checkErr(err);
	fmt.Println("Listen port: ", i.Addr().(*net.TCPAddr).Port);
	c, err := i.AcceptTCP();
	checkErr(err);
	fmt.Println("remote from: ", c.RemoteAddr());
	rConn := getRemoteConn(c);
	rConn.sendCommand("Go Language Server v0.1 for learning");
	rConn.sendCommand("Привет!");
	rConn.sendCommand("Привіт!");
	rConn.sendCommand("Прывітанне!");
	rConn.sendCommand("Hello!");
	rConn.sendCommand("Salut!");
	rConn.sendCommand("ハイ!");
	rConn.sendCommand("您好!");
	rConn.sendCommand("안녕!");
	rConn.sendCommand("Hej!");
}

func main(){
	listen();
}