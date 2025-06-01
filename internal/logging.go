package internal

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// API Error Response
type ErrorResponse struct {
	Message string `json:"message" example:"Server Error"` // Error Message
}

// Log API Errors and return error response
func APIError(w http.ResponseWriter, module string, message string, errorCode int, err error) {
	// Log Error Response
	Error(module, message, err)

	// Create JSON Error Response
	errorJSON, err := json.Marshal(ErrorResponse{
		Message: message,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send Error Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errorCode)
	w.Write(errorJSON)
}

// Log Error with Module Name (Module::Method), Error Message and Error Object
func Error(module string, message string, err error) {
	log.Println("\033[31mERROR\033[0m "+module+" - "+message, err)
}

// Debug Log with Module Name (Module::Method), Message and Objects
func Debug(module string, message string, objects ...any) {
	if len(objects) > 0 {
		log.Println("DEBUG "+module+" - "+message, objects)
	} else {
		log.Println("DEBUG " + module + " - " + message)
	}
}

// Warn Log with Module Name (Module::Method), Message and Objects
func Warn(module string, message string, objects ...any) {
	if len(objects) > 0 {
		log.Println("\033[33mWARN\033[0m  "+module+" - "+message, objects)
	} else {
		log.Println("\033[33mWARN\033[0m  " + module + " - " + message)
	}
}

// Warn Log with Module Name (Module::Method), Message and Objects
func API(statusCode int, method string, path string, responseTime string) {
	log.Println("\033[34mAPI\033[0m  ", strconv.Itoa(statusCode)+" "+method+" "+path+" "+responseTime)
}
