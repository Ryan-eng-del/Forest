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
        "/admin/change_pwd": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "修改密码",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "summary": "修改密码",
                "operationId": "/admin/change_pwd",
                "parameters": [
                    {
                        "description": "body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/adminDto.ChangePwdInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/public.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/admin/info": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "管理员信息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "summary": "管理员信息",
                "operationId": "/admin/info",
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/public.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/adminDto.AdminInfoOutput"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/admin_login/login": {
            "post": {
                "description": "管理员登陆",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "summary": "管理员登陆",
                "operationId": "/admin_login/login",
                "parameters": [
                    {
                        "description": "body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/adminDto.AdminLoginInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/public.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/commonDto.TokensOutput"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/service": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "服务列表",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Service"
                ],
                "summary": "服务列表",
                "operationId": "/service",
                "parameters": [
                    {
                        "type": "string",
                        "description": "关键词",
                        "name": "info",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "每页个数",
                        "name": "page_size",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "当前页数",
                        "name": "page_no",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/public.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/serviceDto.ServiceListOutput"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/service/{service_id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "服务列表",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Service"
                ],
                "summary": "服务详情",
                "operationId": "/service/{service_id}",
                "parameters": [
                    {
                        "type": "string",
                        "description": "服务id",
                        "name": "service_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/public.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/model.ServiceDetail"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "服务删除",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Service"
                ],
                "summary": "服务删除",
                "operationId": "/service/{service_id}",
                "parameters": [
                    {
                        "type": "string",
                        "description": "服务ID",
                        "name": "service_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/public.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "adminDto.AdminInfoOutput": {
            "type": "object",
            "properties": {
                "avatar": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "introduction": {
                    "type": "string"
                },
                "login_time": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "roles": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "adminDto.AdminLoginInput": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string",
                    "example": "123456"
                },
                "username": {
                    "type": "string",
                    "example": "admin"
                }
            }
        },
        "adminDto.ChangePwdInput": {
            "type": "object",
            "required": [
                "password"
            ],
            "properties": {
                "password": {
                    "description": "密码",
                    "type": "string",
                    "example": "123456"
                }
            }
        },
        "commonDto.TokensOutput": {
            "type": "object",
            "properties": {
                "access_token": {
                    "description": "access_token",
                    "type": "string"
                },
                "expires_in": {
                    "description": "expires_in",
                    "type": "integer"
                },
                "scope": {
                    "description": "scope",
                    "type": "string"
                },
                "token_type": {
                    "description": "token_type",
                    "type": "string"
                }
            }
        },
        "model.AccessControl": {
            "type": "object",
            "properties": {
                "black_list": {
                    "type": "string"
                },
                "clientip_flow_limit": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "open_auth": {
                    "type": "boolean"
                },
                "service": {
                    "$ref": "#/definitions/model.Service"
                },
                "service_flow_limit": {
                    "type": "integer"
                },
                "service_id": {
                    "type": "integer"
                },
                "white_host_name": {
                    "type": "string"
                },
                "white_list": {
                    "type": "string"
                }
            }
        },
        "model.GrpcRule": {
            "type": "object",
            "properties": {
                "header_transfor": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "port": {
                    "type": "integer"
                },
                "service": {
                    "$ref": "#/definitions/model.Service"
                },
                "service_id": {
                    "type": "integer"
                }
            }
        },
        "model.HttpRule": {
            "type": "object",
            "properties": {
                "header_transfor": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "need_https": {
                    "type": "integer"
                },
                "need_strip_uri": {
                    "type": "boolean"
                },
                "need_websocket": {
                    "type": "boolean"
                },
                "rule": {
                    "type": "string"
                },
                "rule_type": {
                    "type": "integer"
                },
                "service": {
                    "$ref": "#/definitions/model.Service"
                },
                "service_id": {
                    "type": "integer"
                },
                "url_rewrite": {
                    "type": "string"
                }
            }
        },
        "model.LoadBalance": {
            "type": "object",
            "properties": {
                "check_interval": {
                    "type": "integer"
                },
                "check_method": {
                    "type": "integer"
                },
                "check_timeout": {
                    "type": "integer"
                },
                "forbid_list": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "ip_list": {
                    "type": "string"
                },
                "round_type": {
                    "type": "integer"
                },
                "service": {
                    "$ref": "#/definitions/model.Service"
                },
                "service_id": {
                    "type": "integer"
                },
                "upstream_connect_timeout": {
                    "type": "integer"
                },
                "upstream_header_timeout": {
                    "type": "integer"
                },
                "upstream_idle_timeout": {
                    "type": "integer"
                },
                "upstream_max_idle": {
                    "type": "integer"
                },
                "weight_list": {
                    "type": "string"
                }
            }
        },
        "model.Service": {
            "type": "object",
            "properties": {
                "create_at": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "isDelete": {
                    "type": "integer"
                },
                "load_type": {
                    "type": "integer"
                },
                "service_desc": {
                    "type": "string"
                },
                "service_name": {
                    "type": "string"
                },
                "update_at": {
                    "type": "string"
                }
            }
        },
        "model.ServiceDetail": {
            "type": "object",
            "properties": {
                "access_control": {
                    "$ref": "#/definitions/model.AccessControl"
                },
                "grpc_rule": {
                    "$ref": "#/definitions/model.GrpcRule"
                },
                "http_rule": {
                    "$ref": "#/definitions/model.HttpRule"
                },
                "info": {
                    "$ref": "#/definitions/model.Service"
                },
                "load_balance": {
                    "$ref": "#/definitions/model.LoadBalance"
                },
                "tcp_rule": {
                    "$ref": "#/definitions/model.TcpRule"
                }
            }
        },
        "model.TcpRule": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "port": {
                    "type": "integer"
                },
                "service": {
                    "$ref": "#/definitions/model.Service"
                },
                "service_id": {
                    "type": "integer"
                }
            }
        },
        "public.Response": {
            "type": "object",
            "properties": {
                "data": {},
                "errmsg": {
                    "type": "string"
                },
                "errno": {
                    "$ref": "#/definitions/public.ResponseCode"
                },
                "stack": {},
                "trace_id": {}
            }
        },
        "public.ResponseCode": {
            "type": "integer",
            "enum": [
                0,
                1,
                2,
                3,
                401,
                1000,
                2001
            ],
            "x-enum-varnames": [
                "SuccessCode",
                "UndefErrorCode",
                "ValidErrorCode",
                "InternalErrorCode",
                "InvalidRequestErrorCode",
                "CustomizeCode",
                "GROUPALL_SAVE_FLOWERROR"
            ]
        },
        "serviceDto.ServiceListItemOutput": {
            "type": "object",
            "properties": {
                "createAt": {
                    "type": "string"
                },
                "id": {
                    "description": "id",
                    "type": "integer"
                },
                "load_type": {
                    "description": "类型",
                    "type": "integer"
                },
                "qpd": {
                    "description": "qpd",
                    "type": "integer"
                },
                "qps": {
                    "description": "qps",
                    "type": "integer"
                },
                "service_addr": {
                    "description": "服务地址",
                    "type": "string"
                },
                "service_desc": {
                    "description": "服务描述",
                    "type": "string"
                },
                "service_name": {
                    "description": "服务名称",
                    "type": "string"
                },
                "total_node": {
                    "description": "节点数",
                    "type": "integer"
                },
                "updateAt": {
                    "type": "string"
                }
            }
        },
        "serviceDto.ServiceListOutput": {
            "type": "object",
            "properties": {
                "list": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/serviceDto.ServiceListItemOutput"
                    }
                },
                "total": {
                    "type": "integer"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        },
        "BasicAuth": {
            "type": "basic"
        },
        "OAuth2AccessCode": {
            "type": "oauth2",
            "flow": "accessCode",
            "authorizationUrl": "https://example.com/oauth/authorize",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": " Grants read and write access to administrative information"
            }
        },
        "OAuth2Application": {
            "type": "oauth2",
            "flow": "application",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": " Grants read and write access to administrative information",
                "write": " Grants write access"
            }
        },
        "OAuth2Implicit": {
            "type": "oauth2",
            "flow": "implicit",
            "authorizationUrl": "https://example.com/oauth/authorize",
            "scopes": {
                "admin": " Grants read and write access to administrative information",
                "write": " Grants write access"
            }
        },
        "OAuth2Password": {
            "type": "oauth2",
            "flow": "password",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": " Grants read and write access to administrative information",
                "read": " Grants read access",
                "write": " Grants write access"
            }
        }
    },
    "x-extension-openapi": {
        "example": "value on a json format"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8880",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Go-Gateway  API",
	Description:      "This is a sample server celler server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
