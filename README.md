# GoDjangoSession

Valid for django 3.0.5

Usage:

```
package main

import (
	"encoding/base64"
	"fmt"
	"session/auth"
	"github.com/Kuzyashin/GoDjangoSession"
)

var Session = "NDgwZmI1ZjRlNTEzZjVjZDYyMmFlNDVlMWI0NThiOWJkZTJlZjNmZTp7Il9hdXRoX3VzZXJfaWQiOiIxIiwiX2F1dGhfdXNlcl9iYWNrZW5kIjoiZGphbmdvLmNvbnRyaWIuYXV0aC5iYWNrZW5kcy5Nb2RlbEJhY2tlbmQiLCJfYXV0aF91c2VyX2hhc2giOiJjMzdkNDQzODY5MWVmNWU2M2E3N2RiYmZiZDUzNmZhZDU1MTM3NWQ5In0="
var SecretKey = "LKasdnj1nJN81NDbf891nJANBgfkb>Ghvahv24"

func main() {
	cookie, _ := base64.StdEncoding.DecodeString(Session)
	session, err := DjangoSession.Decode(SecretKey, string(cookie))
	if err != nil {
		panic(err)
	}
	fmt.Printf("session Valid %+v", session)
}

```
