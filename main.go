package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"hintword.com/api/migrations"

	"hintword.com/api/app/common/constants"
	cfg "hintword.com/api/app/configs"
	db "hintword.com/api/app/database"
	"hintword.com/api/app/routes"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	_ = godotenv.Load()

	/*appNew, err := newrelic.NewApplication(
		newrelic.ConfigAppName("Offybox"),
		newrelic.ConfigLicense("58ec1abc7a9f58d4126385f2bfd6cc73FFFFNRAL"),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)

	if err != nil {
		log.Fatal(err)
	}*/

	app := fiber.New()

	config := cfg.Config{}

	envMethod := os.Getenv("ENV_LOAD_METHOD")

	log.Infof("Fetch ENV_LOAD_METHOD: %s\n", envMethod)

	if envMethod == "" {
		fmt.Println("Setting ENV_LOAD_METHOD as LOCAL 😏")
		envMethod = constants.EnvLoadMethodLocal
	}

	log.Infof("Processing with ENV_LOAD_METHOD: %s\n", envMethod)

	switch envMethod {
	case constants.EnvLoadMethodLocal:
		cfg.LoadLocalConfig()
		config = cfg.GetConfig()
		break
	case constants.EnvLoadMethodSSM:
		cfg.LoadSSMConfig()
		config = cfg.GetConfig()
		break
	default:
		log.Fatalln("There is issue you shouldn't be here 🧘🏻")
	}

	/*
		====== Setup DB ============
	*/

	// Connect to Mysql
	db.ConnectMysql()

	// Connect to Google Oauth

	// Migration
	err := migrations.RunMigrations()
	if err != nil {
		log.Error(err)
		return
	}

	// Auto-migrate the GmailMessage model to create the table if it doesn't exist
	app.Use(logger.New())
	// cors
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, Access-Control-Allow-Headers, X-Platform",
		AllowMethods: "*",
	}))

	// Setup routes
	routes.SetupRoutes(app)

	// Load tenant level configuration

	// Run the app and listen on given port
	port := fmt.Sprintf(":%s", config.Port)
	err = app.Listen(port)
	if err != nil {
		log.Error(err)
		return
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}
}
