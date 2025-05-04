package middlewares

import (
	"os"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"hintword.com/api/app/common/utility"
)

func CustomLogger(c *fiber.Ctx) error {
	//startTime := time.Now()
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)

	// Continue processing the request
	if err := c.Next(); err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"method":   c.Method(),
		"header":   utility.CleanString(string(c.Request().Header.Header())),
		"query":    c.Request().URI().QueryString(),
		"request":  utility.CleanString(string(c.Request().Body())),
		"status":   c.Response().StatusCode(),
		"response": utility.CleanString(string(c.Response().Body())),
	}).Info("API LOGGING PATH :" + c.Path())
	//
	// Log response status and body
	//fmt.Printf("Response Status: %d\n", c.Response().StatusCode())
	//fmt.Printf("Response Body: %s\n", c.Response().Body())
	//
	// Calculate and log the request processing time
	//elapsed := time.Since(startTime)
	//fmt.Printf("Elapsed Time: %v\n", elapsed)

	return nil
}
