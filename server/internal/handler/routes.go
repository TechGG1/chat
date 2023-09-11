package handler

import (
	"encoding/json"
	"fmt"
	"github.com/TechGG1/chat/server/mywebsocket"
	"log"
	"net/http"
	"os"

	"github.com/TechGG1/chat/server/internal/handler/middleware"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

func (h *Handler) InitRoutes() *mux.Router {
	r := mux.NewRouter()

	auth := r.PathPrefix("/auth").Subrouter()
	auth.Use(middleware.HeaderMiddleware)

	auth.HandleFunc("/login", h.Login).Methods(http.MethodPost)
	auth.HandleFunc("/register", h.Register).Methods(http.MethodPost)

	chat := r.PathPrefix("/chat").Subrouter()
	chat.Use(middleware.HeaderMiddleware)
	chat.Use(middleware.Auth)

	chat.HandleFunc("/create", h.ChatCreate)
	chat.HandleFunc("/rooms", h.Rooms)
	chat.HandleFunc("/room-message", h.RoomMessages)

	return r
}

var RegisterWebsocketRoute = func(router *mux.Router) {
	pool := mywebsocket.NewPool()
	go pool.Start()
	sb := router.PathPrefix("/v1").Subrouter()

	sb.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		jwtToken := r.URL.Query().Get("jwt")
		jwtSecret := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			handleWebsocketAuthenticationErr(w, err)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			handleWebsocketAuthenticationErr(w, err)
			return
		}

		serveWS(pool, w, r, claims)
	})

}

func serveWS(pool *mywebsocket.Pool, w http.ResponseWriter, r *http.Request, claims jwt.MapClaims) {
	conn, err := mywebsocket.Upgrade(w, r)
	if err != nil {
		return
	}

	client := &mywebsocket.Client{
		Connection: conn,
		Pool:       pool,
		Email:      claims["Email"].(string),
		UserID:     uint(claims["UserID"].(float64)),
	}

	pool.Register <- client
	requestBody := make(chan []byte) // mywebsocket.Message byte array channel
	go client.Read(requestBody)
}

func handleWebsocketAuthenticationErr(w http.ResponseWriter, err error) {
	log.Println("mywebsocket error: ", err)
	w.WriteHeader(http.StatusUnauthorized)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	res := ErrorResponse{Message: err.Error(), Status: false, Code: http.StatusUnauthorized}
	data, err := json.Marshal(res)
	if err != nil {
		return
	}
	w.Write(data)
}
