package handler

import (
	"github.com/gofiber/fiber/v3"
)

// SwaggerHandler handles swagger documentation
type SwaggerHandler struct {
	swaggerURL string
}

// NewSwaggerHandler creates new swagger handler
func NewSwaggerHandler() *SwaggerHandler {
	return &SwaggerHandler{
		swaggerURL: "/swagger/*",
	}
}

// RegisterRoutes registers swagger routes
func (h *SwaggerHandler) RegisterRoutes(app *fiber.App) {
	// Swagger JSON endpoint
	app.Get("/docs/swagger.json", h.ServeSwaggerJSON)

	// Swagger UI endpoint (simple HTML page)
	app.Get("/swagger/*", h.ServeSwaggerUI)
}

// ServeSwaggerJSON serves swagger JSON
func (h *SwaggerHandler) ServeSwaggerJSON(c fiber.Ctx) error {
	swaggerDoc := fiber.Map{
		"openapi": "3.0.0",
		"info": fiber.Map{
			"title":       "Public API Documentation",
			"version":     "1.0.0",
			"description": "Public API for courses and lessons",
		},
		"paths": fiber.Map{
			"/api/v1/public/courses": fiber.Map{
				"get": fiber.Map{
					"summary":     "Get all published courses",
					"description": "Get list of all published courses",
					"tags":        []string{"Public Courses"},
					"responses": fiber.Map{
						"200": fiber.Map{
							"description": "Successful response",
							"content": fiber.Map{
								"application/json": fiber.Map{
									"schema": fiber.Map{
										"type": "array",
										"items": fiber.Map{
											"$ref": "#/components/schemas/CourseResponse",
										},
									},
								},
							},
						},
					},
				},
			},
			"/api/v1/public/courses/{id}": fiber.Map{
				"get": fiber.Map{
					"summary":     "Get course by ID",
					"description": "Get published course by ID with lessons",
					"tags":        []string{"Public Courses"},
					"parameters": []fiber.Map{
						{
							"name":     "id",
							"in":       "path",
							"required": true,
							"schema": fiber.Map{
								"type": "integer",
							},
						},
					},
					"responses": fiber.Map{
						"200": fiber.Map{
							"description": "Successful response",
							"content": fiber.Map{
								"application/json": fiber.Map{
									"schema": fiber.Map{
										"$ref": "#/components/schemas/Course",
									},
								},
							},
						},
					},
				},
			},
			"/api/v1/public/lessons": fiber.Map{
				"get": fiber.Map{
					"summary":     "Get all published lessons",
					"description": "Get list of all published lessons",
					"tags":        []string{"Public Lessons"},
					"parameters": []fiber.Map{
						{
							"name":     "course_id",
							"in":       "query",
							"required": false,
							"schema": fiber.Map{
								"type": "integer",
							},
							"description": "Filter by course ID",
						},
					},
					"responses": fiber.Map{
						"200": fiber.Map{
							"description": "Successful response",
							"content": fiber.Map{
								"application/json": fiber.Map{
									"schema": fiber.Map{
										"type": "array",
										"items": fiber.Map{
											"$ref": "#/components/schemas/Lesson",
										},
									},
								},
							},
						},
					},
				},
			},
			"/api/v1/public/lessons/{id}": fiber.Map{
				"get": fiber.Map{
					"summary":     "Get lesson by ID",
					"description": "Get published lesson by ID with course info",
					"tags":        []string{"Public Lessons"},
					"parameters": []fiber.Map{
						{
							"name":     "id",
							"in":       "path",
							"required": true,
							"schema": fiber.Map{
								"type": "integer",
							},
						},
					},
					"responses": fiber.Map{
						"200": fiber.Map{
							"description": "Successful response",
							"content": fiber.Map{
								"application/json": fiber.Map{
									"schema": fiber.Map{
										"$ref": "#/components/schemas/Lesson",
									},
								},
							},
						},
					},
				},
			},
		},
		"components": fiber.Map{
			"schemas": fiber.Map{
				"Course": fiber.Map{
					"type": "object",
					"properties": fiber.Map{
						"id":           fiber.Map{"type": "integer"},
						"title":        fiber.Map{"type": "string"},
						"description":  fiber.Map{"type": "string"},
						"slug":         fiber.Map{"type": "string"},
						"is_published": fiber.Map{"type": "boolean"},
						"created_at":   fiber.Map{"type": "string", "format": "date-time"},
						"updated_at":   fiber.Map{"type": "string", "format": "date-time"},
						"lessons": fiber.Map{
							"type":  "array",
							"items": fiber.Map{"$ref": "#/components/schemas/Lesson"},
						},
					},
				},
				"CourseResponse": fiber.Map{
					"type": "object",
					"properties": fiber.Map{
						"id":            fiber.Map{"type": "integer"},
						"title":         fiber.Map{"type": "string"},
						"description":   fiber.Map{"type": "string"},
						"slug":          fiber.Map{"type": "string"},
						"lessons_count": fiber.Map{"type": "integer"},
						"created_at":    fiber.Map{"type": "string", "format": "date-time"},
					},
				},
				"Lesson": fiber.Map{
					"type": "object",
					"properties": fiber.Map{
						"id":           fiber.Map{"type": "integer"},
						"course_id":    fiber.Map{"type": "integer"},
						"title":        fiber.Map{"type": "string"},
						"content":      fiber.Map{"type": "string"},
						"slug":         fiber.Map{"type": "string"},
						"order_index":  fiber.Map{"type": "integer"},
						"is_published": fiber.Map{"type": "boolean"},
						"created_at":   fiber.Map{"type": "string", "format": "date-time"},
						"updated_at":   fiber.Map{"type": "string", "format": "date-time"},
						"course":       fiber.Map{"$ref": "#/components/schemas/Course"},
					},
				},
			},
		},
	}

	return c.JSON(swaggerDoc)
}

// ServeSwaggerUI serves Swagger UI HTML page
func (h *SwaggerHandler) ServeSwaggerUI(c fiber.Ctx) error {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Swagger UI - Public API</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.10.0/swagger-ui.css" />
    <style>
        html {
            box-sizing: border-box;
            overflow: -moz-scrollbars-vertical;
            overflow-y: scroll;
        }
        *, *:before, *:after {
            box-sizing: inherit;
        }
        body {
            margin:0;
            background: #fafafa;
        }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.0/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.0/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: "/docs/swagger.json",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>`

	c.Set("Content-Type", "text/html")
	return c.SendString(html)
}
