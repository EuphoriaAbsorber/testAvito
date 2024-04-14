// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/filldb": {
            "post": {
                "description": "Заполнение базы тестовыми данными",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "extra"
                ],
                "summary": "Заполнение базы тестовыми данными",
                "operationId": "fillDB",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Response"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/model.Error"
                        }
                    }
                }
            }
        },
        "/user_banner": {
            "get": {
                "description": "Получение баннера для пользователя",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Получение баннера для пользователя",
                "operationId": "getUserBanner",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "tag_id",
                        "name": "tag_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "feature_id",
                        "name": "feature_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "boolean",
                        "description": "use_last_revision",
                        "name": "use_last_revision",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "token",
                        "name": "token",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.UserBanner"
                        }
                    },
                    "400": {
                        "description": "Некорректные данные",
                        "schema": {
                            "$ref": "#/definitions/model.Error"
                        }
                    },
                    "401": {
                        "description": "Пользователь не авторизован",
                        "schema": {
                            "$ref": "#/definitions/model.Error"
                        }
                    },
                    "403": {
                        "description": "Пользователь не имеет доступа",
                        "schema": {
                            "$ref": "#/definitions/model.Error"
                        }
                    },
                    "404": {
                        "description": "Не найдено",
                        "schema": {
                            "$ref": "#/definitions/model.Error"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/model.Error"
                        }
                    }
                }
            }
        },
        "/users": {
            "get": {
                "description": "Получение пользователей и админов",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "extra"
                ],
                "summary": "Получение пользователей и админов",
                "operationId": "GetUsers",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Response"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/model.Error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.Error": {
            "type": "object",
            "properties": {
                "error": {}
            }
        },
        "model.Response": {
            "type": "object",
            "properties": {
                "body": {}
            }
        },
        "model.UserBanner": {
            "type": "object",
            "properties": {
                "text": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "127.0.0.1:8080",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Banner Service API",
	Description:      "Banner Service back server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
