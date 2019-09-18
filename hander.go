package echat

import (
	log "github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
)

type Handler struct {
	grouper     Grouper
	hub         *Hub
	transformer websocket.Upgrader
}

type HandleOption func(option *Option)

type Option struct {
	grouper Grouper
	parser  MessageParser
	logger  Logger
}

func DefaultOption() *Option {
	return &Option{grouper: NewSimpleGrouper(), parser: DefaultParser, logger: &DefaultLogger{}}
}

func WithGrouper(grouper Grouper) HandleOption {
	return func(option *Option) {
		option.grouper = grouper
	}
}

func WithParser(parser MessageParser) HandleOption {
	return func(option *Option) {
		option.parser = parser
	}
}

func WithLogger(logger Logger) HandleOption {
	return func(option *Option) {
		option.logger = logger
	}
}

func NewHandler(opts ...HandleOption) *Handler {
	opt := DefaultOption()
	for _, val := range opts {
		val(opt)
	}
	hub := NewHub(opt)
	go hub.Run()
	transformer := websocket.Upgrader{}
	hd := &Handler{hub: hub, grouper: opt.grouper, transformer: transformer}
	return hd
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	id := q.Get("id")
	if len(id) == 0 {
		id = uuid.New().String()
		log.Info(id)
	}
	conn, err := h.transformer.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		_ = conn.Close()
		return
	}
	client := &Client{hub: h.hub, conn: conn, send: make(chan []byte, 256), ID: id}
	client.hub.register <- client
	go client.writePump()
	go client.readPump()
	return
}
