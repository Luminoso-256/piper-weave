package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"time"
)

var (
	config Config
	strg   Storage
)

func handleConnection(c net.Conn) {
	/* Get the request path */
	//stolen from petri. thanks former me!
	var requestBytes []byte
	var clen uint16
	path := ""
	clen = 0
	gotclen := false
	reader := bufio.NewReader(c)
	for {
		b, _ := reader.ReadByte()
		requestBytes = append(requestBytes, b)
		//read our "header" so we know when to stop
		if len(requestBytes) >= 2 && !gotclen {
			clb := requestBytes[:2]
			buf := bytes.NewBuffer(clb)
			binary.Read(buf, binary.LittleEndian, &clen)
			gotclen = true
		}
		//if the len of our byte array = the content len, it's time to stop recieving
		if (len(requestBytes) >= int(clen)+2) && gotclen {
			break
		}
	}
	path = string(requestBytes[2:])
	/* Determine Action */

	if strings.HasPrefix(path, "/add") {
		if strings.Contains(path, "?") {
			params := strings.Split(path, "?")[1]
			paramP := strings.Split(params, "&")
			if len(paramP) != 3 {
				c.Write(serializePiperResp([]byte("Wuh oh! You didn't supply enough parameters :(\n usage is /add?name=<name>&url=<url>&desc=<desc>\n"), 0x01))
				c.Close()
				return
			} else {
				currentTime := time.Now()
				strg.addEntry(strings.ReplaceAll(paramP[0], "name=", ""), "piper://"+strings.ReplaceAll(paramP[1], "url=", ""), strings.ReplaceAll(paramP[2], "desc=", "")+" [submitted "+currentTime.Format("2006-01-02")+"]")

				c.Write(serializePiperResp([]byte(fmt.Sprintf("Thanks for contributing to Weave! \n=> piper://%s/ go home", config.Hostname)), 0x01))
			}
		} else {
			c.Write(serializePiperResp([]byte("Wuh oh! You didn't supply parameters :(\n usage is /add?name=<name>&url=<url>&desc=<desc>\n"), 0x01))
		}

	} else {
		// default is listing the agregator
		entries := strg.getEntries()
		response := "# Weave: The Piper Agregator\n----------------------\n"
		for _, entry := range entries {
			response += "=> " + entry.url + " " + entry.name + "\n"
			response += "> " + entry.desc + "\n"
		}
		c.Write(serializePiperResp([]byte(response), 0x01))
	}
	c.Close()
}

func main() {
	fmt.Println("Weave: The Piper Content Agregator")
	strg.init()
	datab, _ := ioutil.ReadFile("data/config.json")
	json.Unmarshal(datab, &config)
	l, err := net.Listen("tcp4", fmt.Sprintf(":%v", (config.Port)))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}
}
