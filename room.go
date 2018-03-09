package main

type room struct {
	// forward は他のクライアントに転送するためのメッセージを保持するチャネルです
	forward chan []byte
	// joinはチャットルームに参加しようとしているクライアントのためのチャネルです。
	join chan *client
	// leaveはチャットルームから退出しようとしているクライアントのためのチャネル
	leave chan *client
	// clients には在室している全てのクライアントが保持される
	clients map[*client]bool
}
