package settings

import (
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/defaults"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
)

func server(s values.Store) *Server {
	return &Server{
		Port: values.NewDefaultInt(defaults.Port, s, "server.port"),
		CORS: cors(s),
	}
}

type Server struct {
	Port values.Int
	CORS CORS
}

func cors(s values.Store) CORS {
	return CORS{
		AllowedOrigins: values.NewDefaultStringSlice([]string{"*"}, s, "server.cors.allowedOrigins"),
		Disabled:       values.NewDefaultBool(defaults.CORSDisabled, s, "server.cors.disabled"),
	}
}

type CORS struct {
	AllowedOrigins values.StringSlice
	Disabled       values.Bool
}
