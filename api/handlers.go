package api

import (
	"encoding/json"
	"encoding/xml"

	"net/http"

	"fmt"
	"time"

	"github.com/BlooperDB/API/db"
	"github.com/BlooperDB/API/utils"
	"github.com/gorilla/mux"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

type GenericResponse struct {
	Success bool           `json:"success"`
	Error   *ErrorResponse `json:"error,omitempty"`
	Data    interface{}    `json:"data,omitempty"`
}

type GenericHandle func(http.ResponseWriter, *http.Request) GenericResponse

func ProcessResponse(handle GenericHandle) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		response := handle(w, r)

		format := r.URL.Query().Get("format")
		pretty := len(r.URL.Query()["pretty"]) > 0

		switch format {
		default:
			w.Header().Set("Content-Type", "application/json")
		case "xml":
			w.Header().Set("Content-Type", "application/xml")
		}

		if response.Error != nil {
			w.WriteHeader(response.Error.Status)
		}

		switch format {
		default:
			encoder := json.NewEncoder(w)

			if pretty {
				encoder.SetIndent("", "    ")
			}

			encoder.Encode(response)
		case "xml":
			encoder := xml.NewEncoder(w)

			if pretty {
				encoder.Indent("", "    ")
			}

			encoder.Encode(response)
		}
	})
}

type DataHandle func(*http.Request) (interface{}, *ErrorResponse)

func DataHandler(handle DataHandle) GenericHandle {
	return func(w http.ResponseWriter, r *http.Request) GenericResponse {
		data, err := handle(r)

		return GenericResponse{
			Success: err == nil,
			Error:   err,
			Data:    data,
		}
	}
}

type RegisterRoute func(method string, path string, handle DataHandle)

func RouteHandler(router *mux.Router, prefix string) RegisterRoute {
	return func(method string, path string, handle DataHandle) {
		route := router.NewRoute()
		route.PathPrefix(prefix)
		route.Path(path)
		route.Methods(method)
		route.Handler(ProcessResponse(DataHandler(handle)))
	}
}

func LoggerHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := makeLogger(w)
		start := time.Now()

		h.ServeHTTP(logger, r)

		end := time.Now()
		latency := end.Sub(start)

		clientIP := r.RemoteAddr
		method := r.Method

		statusCode := logger.Status()

		statusColor := colorForStatus(statusCode)
		methodColor := colorForMethod(method)

		path := r.URL.Path

		fmt.Printf("[API] %v |%s %3d %s| %12v | %21s |%s %-7s %s %s\n",
			end.Format("2006/01/02 - 15:04:05"),
			statusColor, statusCode, reset,
			latency,
			clientIP,
			methodColor, method, reset,
			path,
		)
	})
}

var (
	green   = string([]byte{27, 91, 51, 48, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 51, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 51, 48, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return magenta
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return blue
	case "POST":
		return cyan
	case "PUT":
		return yellow
	case "DELETE":
		return red
	case "PATCH":
		return green
	case "HEAD":
		return magenta
	case "OPTIONS":
		return white
	default:
		return reset
	}
}

func NotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		fmt.Fprint(w, "404 page not found")
	})
}

type AuthDataHandle func(*db.User, *http.Request) (interface{}, *ErrorResponse)

func AuthHandler(handle AuthDataHandle) func(*http.Request) (interface{}, *ErrorResponse) {
	return func(r *http.Request) (interface{}, *ErrorResponse) {
		authUser := db.GetAuthUser(r)

		if authUser == nil {
			return nil, &utils.Error_blooper_token_invalid
		}

		return handle(authUser, r)
	}
}
