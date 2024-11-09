package errors

import (
	"encoding/json"
	"fmt"
	"net"
)

func ReturnError(conn net.Conn, errMsg string) {
	type ErrorResponse struct {
		Error string `json:"error"`
	}
	response := ErrorResponse{
		Error: errMsg,
	}
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(response); err != nil {
		fmt.Println("Error al codificar JSON:", err)
	}
}

func SendJSONResponse(conn net.Conn, data interface{}) error {
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(data); err != nil {
		fmt.Println("Error al codificar JSON:", err)
		return fmt.Errorf("error al codificar JSON: %w", err)
	}
	return nil
}
