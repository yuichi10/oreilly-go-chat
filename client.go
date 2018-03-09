package main

import "github.com/gorilla/websocket"

// client はチャットを行ってる一人のユーザーを表します。
type client struct {
	// docket はこのクライアントのためのwebsocket
	socket *websocket.Conn
	// sendはメッセージが送られるチャネル
	send chan []byte
	// roomはこのクライアントが参加してるチャットルーム
	room *room
}
