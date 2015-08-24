// Copyright 2015 The mp Authors. All rights reserved.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

package main

type HandshakeResponsePacket struct {
	PacketHeader

	capabilityFlags []byte
	maxPacketSize   []byte
	characterSet    byte
	reserved        []byte
	username        []byte
	authResponseLen byte
	authResponse    []byte
	database        []byte
	authPluginName  []byte
	attrsLen        []byte
	attrs           []byte
}

func (p *HandshakeResponsePacket) Init(ihp *InitialHandshakePacket) []byte {
	capabilities := clientProtocol41 |
		clientSecureConnection |
		clientLongPassword |
		clientTransactions |
		clientLocalFiles |
		clientConnectWithDB
		// ihp.CapabilityFlag & clientLongFlag

	if ihp.CapabilityFlag & clientLongFlag {
		capabilities |= clientLongFlag
	}

	// 4 capability flags, CLIENT_PROTOCOL_41 always set
	// 4 max-packet size
	// 1 character set
	// string[23] reserved (all [0])
	l := 4 + 4 + 1 + 23

	// string[NUL] username
	l += len("test") + 1

	if capabilities & clientPluginAuthLenEncClientData {
		// lenenc-int length of auth-response
		// string[n] auth-response

	} else if capabilities & clientSecureConnection {
		// 1 length of auth-response
		// string[n] auth-response
		l += 1 + len("auth")

	} else {
		// string[NUL] auth-response

	}

	if capabilities & clientConnectWithDB {
		// string[NUL] database
		l += len("db") + 1
	}

	if capabilities & clientPluginAuth {
		// string[NUL] auth plugin name
		l += len(ihp.authPluginName) + 1
	}

	if capabilities & clientConnectAttrs {
		// lenenc-int length of all key-values
		// lenenc-str key
		// lenenc-str value
		// ...
		// if-more data in 'length of all key-values', more keys and value pairs

	}

	// header
	l += 4

	// merge data
	data := make([]byte, l)

	return data
}
