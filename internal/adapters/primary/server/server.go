package server

import "github.com/labstack/echo/v4"

type Server struct {
	echo *echo.Echo
	port int
}

func NewServer(e *echo.Echo, port int) *Server {
	return &Server{
		echo: e,
		port: port,
	}
}
