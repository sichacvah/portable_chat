package api

import "github.com/sichacvah/portable_chat/model"

type Hub struct {
	register    chan *WebConn
	unregister  chan *WebConn
	broadcast   chan *model.Message
	stop        chan string
	connections map[*WebConn]bool
}

var hub = &Hub{
	register:    make(chan *WebConn),
	unregister:  make(chan *WebConn),
	broadcast:   make(chan *model.Message),
	stop:        make(chan string),
	connections: make(map[*WebConn]bool),
}

func PublishAndForget(message *model.Message) {
	go func() {
		hub.Broadcast(message)
	}()
}

func (h *Hub) Register(webConn *WebConn) {
	h.register <- webConn
}

func (h *Hub) Unregister(webConn *WebConn) {
	h.unregister <- webConn
}

func (h *Hub) Broadcast(message *model.Message) {
	if message != nil {
		h.broadcast <- message
	}
}

func (h *Hub) Stop() {
	h.stop <- "all"
}

func (h *Hub) Start() {
	go func() {
		for {
			select {

			case webCon := <-h.register:
				h.connections[webCon] = true
			case webCon := <-h.unregister:
				if _, ok := h.connections[webCon]; ok {
					delete(h.connections, webCon)
					close(webCon.Send)
				}
			case msg := <-h.broadcast:
				for webCon := range h.connections {
					if ShouldSendEvent(webCon, msg) {
						select {
						case webCon.Send <- msg:
						default:
							close(webCon.Send)
							delete(h.connections, webCon)
						}
					}
				}
			}
		}
	}()
}

func ShouldSendEvent(webCon *WebConn, msg *model.Message) bool {
	_, _ = webCon, msg
	return true
}
