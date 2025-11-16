package main

import (
	"chatapp/database"
	routes "chatapp/src"
	"chatapp/src/user"
	"chatapp/src/websocket"
	"flag"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func GinMiddleware(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Request.Header.Del("Origin")

		c.Next()
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Env can't be loaded")
	}
	db := database.InitDB()
	database.MigrateDB(db)
	sqlDb, _ := db.DB()
	defer sqlDb.Close()

	router := gin.New()

	flag.Parse()
	chatHub := websocket.NewHub(db)
	notifHub := websocket.NewHub(db)
	go chatHub.Run()
	go notifHub.Run()

	router.Use(GinMiddleware("http://localhost:5173"))

	router.POST("/user/register", user.RegisterUserHandler(db))
	router.POST("/user/login", user.LoginUserHandler(db))
	router.Any("/ws/chat", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(chatHub, w, r, "chat")
	}))
	router.Any("/ws/notif", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(notifHub, w, r, "notif")
	}))

	router.Use(user.AuthMiddlewareHandler)
	routes.RouteMap(router, db)

	if err := router.Run(":8080"); err != nil {
		log.Fatal("failed run app: ", err)
	}
}
