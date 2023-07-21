package utils

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os/exec"
	"runtime"
)

func LoadEnvVariables() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func ChunkSlice[T any](slice []T, chunkSize int) [][]T {
	var result [][]T

	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}

		result = append(result, slice[i:end])
	}

	return result
}

func PrettyPrintJson(j any) {
	marshaled, err := json.MarshalIndent(j, "", "   ")
	if err != nil {
		log.Fatalf("marshaling error: %s", err)
	}
	fmt.Println(string(marshaled))
}

func OpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
