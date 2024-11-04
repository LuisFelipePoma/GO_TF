package errors

import (
	"encoding/json"
	"fmt"
	"net"
	"github.com/LuisFelipePoma/Movies_Recomender_With_Golang/src/api/types"
)

func ReturnError(conn net.Conn, errMsg string) {
	response := types.Response{
		Error: errMsg,
	}
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(response); err != nil {
		fmt.Println("Error al codificar JSON:", err)
	}
}

func SendJSONResponse(conn net.Conn, data types.Response) error {
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(data); err != nil {
		fmt.Println("Error al codificar JSON:", err)
		return fmt.Errorf("error al codificar JSON: %w", err)
	}
	return nil
}
