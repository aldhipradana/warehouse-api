package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ActionLogger logs details of POST, PUT, and DELETE requests to a file based on the date
func ActionLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		// Only log for data-modifying methods
		if method == "POST" || method == "PUT" || method == "DELETE" {
			start := time.Now()
			path := c.Request.URL.Path
			rawQuery := c.Request.URL.RawQuery

			// Capture the payload
			var body []byte
			if c.Request.Body != nil {
				body, _ = io.ReadAll(c.Request.Body)
				// Restore the request body for subsequent handlers
				c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			}

			// Process request
			c.Next()

			// Log after request is processed
			end := time.Now()
			latency := end.Sub(start)
			status := c.Writer.Status()
			clientIP := c.ClientIP()

			// Prepare payload for logging (flatten to single line)
			payload := string(body)
			if payload == "" {
				payload = "none"
			} else {
				payload = strings.ReplaceAll(payload, "\n", " ")
				payload = strings.ReplaceAll(payload, "\r", "")
				payload = strings.Join(strings.Fields(payload), " ") // Remove extra spaces
			}

			// Prepare query for logging
			query := rawQuery
			if query == "" {
				query = "none"
			}

			// Prepare log message
			logMsg := fmt.Sprintf("[ACTION] %s %s | Status: %d | Latency: %v | IP: %s | Query: %s | Payload: %s\n",
				method,
				path,
				status,
				latency,
				clientIP,
				query,
				payload,
			)

			// Write to file based on date
			logDir := "log"
			fileName := time.Now().Format("2006-01-02") + ".log"
			filePath := filepath.Join(logDir, fileName)

			// Ensure directory exists
			if err := os.MkdirAll(logDir, 0755); err != nil {
				log.Printf("Error creating log directory: %v", err)
				return
			}

			// Open file in append mode, create if not exists
			f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Printf("Error opening log file: %v", err)
				return
			}
			defer f.Close()

			// Add timestamp to the line and write
			timestamp := time.Now().Format("15:04:05")
			if _, err := f.WriteString(fmt.Sprintf("%s %s", timestamp, logMsg)); err != nil {
				log.Printf("Error writing to log file: %v", err)
			}
		} else {
			c.Next()
		}
	}
}
