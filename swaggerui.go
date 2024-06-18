package swaggerui

import (
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
)

//go:generate go run generate.go

//go:embed embed
var swagfs embed.FS

func byteHandler(b []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write(b)
		if err != nil {
			slog.Error(fmt.Sprintf("error writing swagger spec: %v", err))
			return
		}
	}
}

// Handler returns a handler that will serve a self-hosted Swagger UI with your spec embedded
func Handler(spec []byte) http.Handler {
	// render the index template with the proper spec name inserted
	static, _ := fs.Sub(swagfs, "embed")
	mux := http.NewServeMux()
	mux.HandleFunc("/swagger_spec", byteHandler(spec))
	mux.Handle("/", http.FileServer(http.FS(static)))
	return mux
}
