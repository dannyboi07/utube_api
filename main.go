package main

import (
	"fmt"
	"net/http"
	"os"
	"utube/common"
	"utube/controller"
	"utube/db"
	"utube/redis"
	"utube/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	utils.ReadEnv()
	if err := common.InitKeys(); err != nil {
		utils.Log.Fatalln("Failed to load keys, err:", err)
	}

	if err := db.Init(); err != nil {
		utils.Log.Fatalln("Failed to connect to database, exiting..., err: ", err)
	}
	defer db.CloseDb()

	if err := redis.Init(); err != nil {
		utils.Log.Fatalln("Failed to connect to Redis, exiting..., err:", err)
	}
	defer redis.Close()

	var r *chi.Mux = chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowOriginFunc: func(r *http.Request, origin string) bool {
			return origin == "http://localhost:3000"
		},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"*"},
	}))

	r.Route("/utubeapi", func(r chi.Router) {

		r.Get("/ping", controller.Ping)

		r.Route("/v1", func(r chi.Router) {
			r.Group(func(r chi.Router) {

				r.Route("/auth", func(r chi.Router) {
					r.Get("/refresh", controller.RefreshToken)

					r.Group(func(r chi.Router) {
						r.Use(utils.JsonRoute)
						r.Post("/login", controller.Login)
						r.Post("/register", controller.Register)
					})
				})

				r.Route("/video", func(r chi.Router) {
					// r.Use(utils.JsonRoute)

					r.Route("/upload", func(r chi.Router) {
						r.Use(utils.AuthMiddleware)
						// Get video upload status
						r.Get("/status/{video_name}", controller.UploadStatus)

						r.Group(func(r chi.Router) {
							r.Use(utils.JsonRoute)
							// Register video upload request
							r.Post("/create", controller.CreateVideo)
							// Upload video chunk
							r.Post("/chunk", controller.UploadChunk)
						})
					})
				})
			})

			// multipart/form-data routes
			// r.Group(func(r chi.Router) {
			// 	r.Use(utils.MultipartFormDataRoute)

			// 	r.Post("/video", controller.UploadVideo)
			// })
		})
	})

	utils.Log.Println(os.Getenv("PORT"), os.Getenv("ENV"))

	utils.Log.Printf("Starting server on port: %s", os.Getenv("PORT"))
	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT")), r); err != nil {
		utils.Log.Fatalf("Failed to start server on port: %s, err: %s", os.Getenv("PORT"), err.Error())
	}
}
