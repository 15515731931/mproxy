// Copyright 2015 The mp Authors. All rights reserved.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile)

	client := new(Client)
	err := client.Connet("tcp", "192.168.22.76:3306")
	if err != nil {
		log.Fatal(err)
	}
}
