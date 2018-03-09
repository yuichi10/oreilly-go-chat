package main

type room struct {
	// forward は他のクライアントに転送するためのメッセージを保持するチャネルです
	forward chan []byte
}
