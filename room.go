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

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			// すべてのクライアントにメッセージを転送
			for client := range r.clients {
				select {
				case client.send <- msg:
					// メッセ＝時を送信
				default:
					// 送信に失敗
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}
