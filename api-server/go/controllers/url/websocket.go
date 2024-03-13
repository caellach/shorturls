package url

import (
	"encoding/json"
	"log"

	"github.com/caellach/shorturl/api-server/go/pkg/middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/websocket/v2"
	"github.com/valyala/fasthttp"
)

func disconnectWebsocket(wsContext *websocket.Conn, userId string) {
	if userId == "" {
		return
	}

	log.Printf("WebSocket disconnected: %s@%s", userId, getRemoteAddr(wsContext))

	for i, conn := range websocketConnections[userId] {
		if conn == wsContext {
			websocketConnections[userId] = append(websocketConnections[userId][:i], websocketConnections[userId][i+1:]...)
		}
	}
}

func forceDisconnectWebsocket(wsContext *websocket.Conn, userId string) {
	log.Printf("Forcefully disconnecting websocket: %s", getRemoteAddr(wsContext))
	// kill the socket
	wsContext.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	disconnectWebsocket(wsContext, userId)
}

func getRemoteAddr(wsContext *websocket.Conn) string {
	ctx := wsContext.Locals("ctx").(*fasthttp.RequestCtx)
	xForwardedFor := ctx.Request.Header.Peek("X-Forwarded-For")
	ip := string(xForwardedFor)
	if ip == "" {
		ip = wsContext.RemoteAddr().String()
	}
	return ip
}

var MAX_WS_MSG_SIZE = 1024

func urlWs(wsContext *websocket.Conn) {
	// WebSocket connected
	authorized := false
	userId := ""
	log.Printf("WebSocket connected: %s", getRemoteAddr(wsContext))

	// Listen for messages infinitely
	for {
		msgType, msg, err := wsContext.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			disconnectWebsocket(wsContext, userId)
			break
		}

		if !authorized {
			if msgType == websocket.TextMessage {
				if len(msg) == 0 || len(msg) > MAX_WS_MSG_SIZE {
					forceDisconnectWebsocket(wsContext, userId)
					break
				}

				// decode the message
				var message map[string]interface{}
				err := json.Unmarshal(msg, &message)
				if err != nil {
					log.Println("failed to unmarshal json message:", err)
					forceDisconnectWebsocket(wsContext, userId)
					break
				}

				// check the action
				action, ok := message["action"].(string)
				if !ok {
					log.Println("invalid message: action not found")
					forceDisconnectWebsocket(wsContext, userId)
					break
				}

				if action == "auth" {
					// check the token
					token, ok := message["token"].(string)
					if !ok {
						log.Println("invalid message: token not found")
						forceDisconnectWebsocket(wsContext, userId)
						break
					}

					// check the token, it should be signed by the server
					decodedToken, err := middleware.ValidateToken(token)
					if err != nil {
						log.Println("failed to verify token:", err)
						forceDisconnectWebsocket(wsContext, userId)
						break
					}

					// check the user id
					userId, ok = decodedToken.Claims.(jwt.MapClaims)["sub"].(string)
					if !ok {
						log.Println("invalid token: user id not found")
						forceDisconnectWebsocket(wsContext, userId)
						break
					}

					// register the websocket
					websocketConnections[userId] = append(websocketConnections[userId], wsContext)
					authorized = true
					// websocket id
					_logger.Printf("WebSocket authorized: %s@%s", userId, getRemoteAddr(wsContext))
					// send the auth response
					wsContext.WriteMessage(websocket.TextMessage, []byte(`{"action":"auth"}`))
					continue
				}
			}
			forceDisconnectWebsocket(wsContext, userId)
			break
		} else {
			if msgType == websocket.PingMessage {
				if string(msg) == "ping" {
					wsContext.WriteMessage(websocket.PongMessage, []byte("pong"))
				}
			} else if msgType == websocket.TextMessage {
				if len(msg) == 0 || len(msg) > MAX_WS_MSG_SIZE {
					forceDisconnectWebsocket(wsContext, userId)
					break
				}

				var message map[string]interface{}
				err := json.Unmarshal(msg, &message)
				if err != nil {
					log.Println("failed to unmarshal json message:", err)
					forceDisconnectWebsocket(wsContext, userId)
					break
				}

				// check the action
				action, ok := message["action"].(string)
				if !ok {
					log.Println("invalid message: action not found")
					forceDisconnectWebsocket(wsContext, userId)
					break
				}

				if action == "ping" {
					wsContext.WriteMessage(websocket.TextMessage, []byte(`{"action":"pong"}`))
					continue
				} else {
					log.Println("invalid action:", action)
					forceDisconnectWebsocket(wsContext, userId)
					break
				}
			}
		}
	}
}
