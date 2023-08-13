// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "MIT",
            "url": "http://www.opensource.org/licenses/MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/repository": {
            "get": {
                "description": "Get a list of all backup repositories",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "backup-repositories"
                ],
                "summary": "Get all backup repositories",
                "responses": {
                    "200": {
                        "description": "List of backup repositories",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/backuprepo.BackupRepo"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new backup repository configuration with the given data",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "backup-repositories"
                ],
                "summary": "Create a new backup repository configuration",
                "parameters": [
                    {
                        "description": "Backup repository data",
                        "name": "backupRepo",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/backuprepo.BackupRepo"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success response\"“",
                        "schema": {
                            "$ref": "#/definitions/server.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid request body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/repository/{repo_name}": {
            "get": {
                "description": "Get the backup repository with the given name",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "backup-repositories"
                ],
                "summary": "Get a backup repository by name",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Name of the backup repository to retrieve",
                        "name": "repo_name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Backup repository data",
                        "schema": {
                            "$ref": "#/definitions/backuprepo.BackupRepo"
                        }
                    },
                    "404": {
                        "description": "Backup repository not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a backup repository by its name",
                "tags": [
                    "backup-repositories"
                ],
                "summary": "Delete a backup repository",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Name of the backup repository to delete",
                        "name": "repo_name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success response",
                        "schema": {
                            "$ref": "#/definitions/server.SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/repository/{repo_name}/storage/": {
            "get": {
                "description": "Get all storages associated with a backup repository by its name",
                "tags": [
                    "backup-repositories"
                ],
                "summary": "Get backup repository storages",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Name of the backup repository",
                        "name": "repo_name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "List of storages associated with the backup repository",
                        "schema": {
                            "type": "array",
                            "items": {}
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/repository/{repo_name}/storage/{storage_name}": {
            "post": {
                "description": "Associate a storage with a backup repository by their names",
                "tags": [
                    "backup-repositories"
                ],
                "summary": "Add storage to backup repository",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Name of the backup repository",
                        "name": "repo_name",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Name of the storage",
                        "name": "storage_name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success response",
                        "schema": {
                            "$ref": "#/definitions/server.SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "description": "Remove the association between a storage and a backup repository by their names",
                "tags": [
                    "backup-repositories"
                ],
                "summary": "Remove storage from backup repository",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Name of the backup repository",
                        "name": "repo_name",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Name of the storage",
                        "name": "storage_name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success response",
                        "schema": {
                            "$ref": "#/definitions/server.SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/storage/": {
            "get": {
                "description": "Get all storage configurations",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "storages"
                ],
                "responses": {
                    "200": {
                        "description": "List of all storage configurations",
                        "schema": {
                            "type": "array",
                            "items": {}
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/storage/{storage_conf}": {
            "get": {
                "description": "Get the storage configuration by its name",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "storages"
                ],
                "summary": "Get storage by name",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Name of the storage",
                        "name": "storage_conf",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Storage configuration",
                        "schema": {}
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete the storage configuration by its name",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "storages"
                ],
                "summary": "Delete storage",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Name of the storage",
                        "name": "storage_conf",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success response",
                        "schema": {
                            "$ref": "#/definitions/server.SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/storage/{storage_type}": {
            "post": {
                "description": "Create a new storage configuration",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "storages"
                ],
                "summary": "Create storage",
                "parameters": [
                    {
                        "enum": [
                            "s3"
                        ],
                        "type": "string",
                        "description": "Storage type (s3)",
                        "name": "storage_type",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Storage configuration to create",
                        "name": "storage",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/storage.S3Storage"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success response",
                        "schema": {
                            "$ref": "#/definitions/server.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "backuprepo.BackupRepo": {
            "type": "object",
            "properties": {
                "credentials": {
                    "$ref": "#/definitions/backuprepo.Credentials"
                },
                "name": {
                    "type": "string"
                },
                "remote_url": {
                    "type": "string"
                },
                "schedule": {
                    "type": "string"
                }
            }
        },
        "backuprepo.Credentials": {
            "type": "object",
            "properties": {
                "key_path": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "server.SuccessResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "storage.S3Storage": {
            "type": "object",
            "properties": {
                "access_key": {
                    "type": "string"
                },
                "bucket_name": {
                    "type": "string"
                },
                "endpoint": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "region": {
                    "type": "string"
                },
                "secret_key": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "GitEcho API",
	Description:      "REST API for GitEcho, a tool for backing up Git repositories",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}