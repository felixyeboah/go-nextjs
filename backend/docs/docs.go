// Package docs provides Swagger documentation for the API.
package docs

import (
	"github.com/swaggo/swag"
)

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Fullstack API",
	Description:      "A modern fullstack application with Go and Next.js",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

const docTemplate = `{
    "swagger": "2.0",
    "info": {
        "description": "A modern fullstack application with Go and Next.js",
        "title": "Fullstack API",
        "version": "1.0"
    },
    "host": "localhost:8080",
    "basePath": "/api/v1",
    "paths": {}
}`
