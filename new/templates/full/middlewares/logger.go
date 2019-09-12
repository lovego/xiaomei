package middlewares

import (
	"github.com/lovego/config"
	"github.com/lovego/goa/middlewares"
)

var Logger = middlewares.NewLogger(config.HttpLogger())
