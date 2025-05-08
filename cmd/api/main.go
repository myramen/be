package main

import (
	"log"

	"github.com/myramen/be/internal/app/coupon"
	"github.com/myramen/be/internal/app/order"
	"github.com/myramen/be/internal/pkg/config"
	"github.com/myramen/be/internal/pkg/db/mysql"
	"github.com/myramen/be/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()

	db, err := mysql.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	orderRepo := mysql.NewOrderRepository(db)
	couponRepo := mysql.NewCouponRepository(db)

	couponService := coupon.NewService(couponRepo)
	orderService := order.NewService(orderRepo, couponRepo)

	couponHandler := coupon.NewHandler(couponService)
	orderHandler := order.NewHandler(orderService)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(middleware.ErrorHandler())

	api := router.Group("/api/v1")
	{
		orderHandler.RegisterRoutes(api)
		couponHandler.RegisterRoutes(api)
	}

	log.Printf("Starting server on port %s", config.AppConfig.AppPort)
	if err := router.Run(":" + config.AppConfig.AppPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
