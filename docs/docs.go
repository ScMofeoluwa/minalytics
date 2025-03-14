// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/analytics/browsers": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Retrieves browser stats",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Analytics"
                ],
                "summary": "Get Browsers",
                "parameters": [
                    {
                        "type": "string",
                        "description": "app tracking ID",
                        "name": "trackingID",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "start date",
                        "name": "startDate",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "end date",
                        "name": "endDate",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "stats fetched successfully",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.BrowserResponse"
                        }
                    },
                    "400": {
                        "description": "invalid request paramaters",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.APIStatus"
                        }
                    },
                    "500": {
                        "description": "failed to fetch browsers",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.APIStatus"
                        }
                    }
                }
            }
        },
        "/analytics/countries": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Retrieves country stats",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Analytics"
                ],
                "summary": "Get Countries",
                "parameters": [
                    {
                        "type": "string",
                        "description": "app tracking ID",
                        "name": "trackingID",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "start date",
                        "name": "startDate",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "end date",
                        "name": "endDate",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "stats fetched successfully",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.CountryResponse"
                        }
                    },
                    "400": {
                        "description": "invalid request paramaters",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.APIStatus"
                        }
                    },
                    "500": {
                        "description": "failed to fetch countries",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.APIStatus"
                        }
                    }
                }
            }
        },
        "/analytics/devices": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Retrieves device stats",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Analytics"
                ],
                "summary": "Get Devices",
                "parameters": [
                    {
                        "type": "string",
                        "description": "app tracking ID",
                        "name": "trackingID",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "start date",
                        "name": "startDate",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "end date",
                        "name": "endDate",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "stats fetched successfully",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.DeviceResponse"
                        }
                    },
                    "400": {
                        "description": "invalid request paramaters",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.APIStatus"
                        }
                    },
                    "500": {
                        "description": "failed to fetch devices",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.APIStatus"
                        }
                    }
                }
            }
        },
        "/analytics/os": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Retrieves operating system stats",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Analytics"
                ],
                "summary": "Get OS",
                "parameters": [
                    {
                        "type": "string",
                        "description": "app tracking ID",
                        "name": "trackingID",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "start date",
                        "name": "startDate",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "end date",
                        "name": "endDate",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "stats fetched successfully",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.OSResponse"
                        }
                    },
                    "400": {
                        "description": "invalid request paramaters",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.APIStatus"
                        }
                    },
                    "500": {
                        "description": "failed to fetch operating systems",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.APIStatus"
                        }
                    }
                }
            }
        },
        "/analytics/pages": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Retrieves page stats",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Analytics"
                ],
                "summary": "Get Pages",
                "parameters": [
                    {
                        "type": "string",
                        "description": "app tracking ID",
                        "name": "trackingID",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "start date",
                        "name": "startDate",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "end date",
                        "name": "endDate",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "stats fetched successfully",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.PageResponse"
                        }
                    },
                    "400": {
                        "description": "invalid request paramaters",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.APIStatus"
                        }
                    },
                    "500": {
                        "description": "failed to fetch pages",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.APIStatus"
                        }
                    }
                }
            }
        },
        "/analytics/pageviews": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Retrieves page views",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Analytics"
                ],
                "summary": "Get PageViews",
                "parameters": [
                    {
                        "type": "string",
                        "description": "app tracking ID",
                        "name": "trackingID",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "start date",
                        "name": "startDate",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "end date",
                        "name": "endDate",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "stats fetched successfully",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.PageViewResponse"
                        }
                    },
                    "400": {
                        "description": "invalid request paramaters",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.APIStatus"
                        }
                    },
                    "500": {
                        "description": "failed to fetch visitors",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.APIStatus"
                        }
                    }
                }
            }
        },
        "/analytics/referrals": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Retrieves referral stats",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Analytics"
                ],
                "summary": "Get Referrals",
                "parameters": [
                    {
                        "type": "string",
                        "description": "app tracking ID",
                        "name": "trackingID",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "start date",
                        "name": "startDate",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "end date",
                        "name": "endDate",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "stats fetched successfully",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.ReferralResponse"
                        }
                    },
                    "400": {
                        "description": "invalid request paramaters",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.APIStatus"
                        }
                    },
                    "500": {
                        "description": "failed to fetch referrals",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.APIStatus"
                        }
                    }
                }
            }
        },
        "/analytics/track": {
            "get": {
                "description": "Tracks an event based on encoded data",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Analytics"
                ],
                "summary": "Track an event",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Base64 encoded event data",
                        "name": "data",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Event tracked successfully",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.APIStatus"
                        }
                    },
                    "400": {
                        "description": "Invalid base64 or JSON data",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.APIStatus"
                        }
                    },
                    "500": {
                        "description": "Failed to resolve geolocation or track event",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.APIStatus"
                        }
                    }
                }
            }
        },
        "/analytics/visitors": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Retrieves visitors",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Analytics"
                ],
                "summary": "Get Visitors",
                "parameters": [
                    {
                        "type": "string",
                        "description": "app tracking ID",
                        "name": "trackingID",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "start date",
                        "name": "startDate",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "end date",
                        "name": "endDate",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "stats fetched successfully",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.VisitorResponse"
                        }
                    },
                    "400": {
                        "description": "invalid request paramaters",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.APIStatus"
                        }
                    },
                    "500": {
                        "description": "failed to fetch visitors",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.APIStatus"
                        }
                    }
                }
            }
        },
        "/apps": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Retrieves user apps",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Apps"
                ],
                "summary": "Get Apps",
                "responses": {
                    "200": {
                        "description": "apps fetched successfully",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.AppResponse"
                        }
                    },
                    "401": {
                        "description": "userID not found in context",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.APIStatus"
                        }
                    },
                    "500": {
                        "description": "failed to fetch apps",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.APIStatus"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "creates an app",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Apps"
                ],
                "summary": "Create App",
                "parameters": [
                    {
                        "description": "app name",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.CreateAppRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "app created successfully",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.AppResponse"
                        }
                    },
                    "400": {
                        "description": "invalid request body",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.APIStatus"
                        }
                    },
                    "401": {
                        "description": "userID not found in context",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.APIStatus"
                        }
                    },
                    "500": {
                        "description": "failed to create app",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.APIStatus"
                        }
                    }
                }
            }
        },
        "/apps/{trackingID}": {
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Updates app by tracking ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Apps"
                ],
                "summary": "Delete App",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Tracking ID of the app to delete",
                        "name": "trackingID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "app successfully deleted",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "invalid request body",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.APIStatus"
                        }
                    },
                    "401": {
                        "description": "userID not found in context",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.APIStatus"
                        }
                    },
                    "500": {
                        "description": "failed to delete app",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.APIStatus"
                        }
                    }
                }
            },
            "patch": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Updates app by tracking ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Apps"
                ],
                "summary": "Update App",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Tracking ID of the app to delete",
                        "name": "trackingID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "app name",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.CreateAppRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "app name successfully changed",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.AppResponse"
                        }
                    },
                    "400": {
                        "description": "invalid request body",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.APIStatus"
                        }
                    },
                    "401": {
                        "description": "userID not found in context",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.APIStatus"
                        }
                    },
                    "500": {
                        "description": "failed to update app",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.APIStatus"
                        }
                    }
                }
            }
        },
        "/auth/{provider}": {
            "get": {
                "description": "Initiates OAuth authentication with the specified provider and returns a JWT token upon successful login.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "User Sign-In",
                "parameters": [
                    {
                        "type": "string",
                        "description": "OAuth provider (e.g., google, github)",
                        "name": "provider",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "JWT token",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid provider or missing provider",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.APIStatus"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.APIStatus"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "github_com_ScMofeoluwa_minalytics_shared.APIStatus": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "github_com_ScMofeoluwa_minalytics_shared.App": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "trackingID": {
                    "type": "string"
                }
            }
        },
        "github_com_ScMofeoluwa_minalytics_shared.AppResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.App"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "github_com_ScMofeoluwa_minalytics_shared.BrowserResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.BrowserStats"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "github_com_ScMofeoluwa_minalytics_shared.BrowserStats": {
            "type": "object",
            "properties": {
                "browser": {
                    "type": "string"
                },
                "percentage": {
                    "type": "integer"
                }
            }
        },
        "github_com_ScMofeoluwa_minalytics_shared.CountryResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.CountryStats"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "github_com_ScMofeoluwa_minalytics_shared.CountryStats": {
            "type": "object",
            "properties": {
                "country": {
                    "type": "string"
                },
                "percentage": {
                    "type": "integer"
                }
            }
        },
        "github_com_ScMofeoluwa_minalytics_shared.CreateAppRequest": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        },
        "github_com_ScMofeoluwa_minalytics_shared.DeviceResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.DeviceStats"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "github_com_ScMofeoluwa_minalytics_shared.DeviceStats": {
            "type": "object",
            "properties": {
                "device": {
                    "type": "string"
                },
                "percentage": {
                    "type": "integer"
                }
            }
        },
        "github_com_ScMofeoluwa_minalytics_shared.OSResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.OSStats"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "github_com_ScMofeoluwa_minalytics_shared.OSStats": {
            "type": "object",
            "properties": {
                "operating_system": {
                    "type": "string"
                },
                "percentage": {
                    "type": "integer"
                }
            }
        },
        "github_com_ScMofeoluwa_minalytics_shared.PageResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.PageStats"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "github_com_ScMofeoluwa_minalytics_shared.PageStats": {
            "type": "object",
            "properties": {
                "path": {
                    "type": "string"
                },
                "visitor_count": {
                    "type": "integer"
                }
            }
        },
        "github_com_ScMofeoluwa_minalytics_shared.PageViewResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.PageViewStats"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "github_com_ScMofeoluwa_minalytics_shared.PageViewStats": {
            "type": "object",
            "properties": {
                "time": {
                    "type": "string"
                },
                "views": {
                    "type": "integer"
                }
            }
        },
        "github_com_ScMofeoluwa_minalytics_shared.ReferralResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.ReferralStats"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "github_com_ScMofeoluwa_minalytics_shared.ReferralStats": {
            "type": "object",
            "properties": {
                "referrer": {
                    "type": "string"
                },
                "visitor_count": {
                    "type": "integer"
                }
            }
        },
        "github_com_ScMofeoluwa_minalytics_shared.VisitorResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/github_com_ScMofeoluwa_minalytics_shared.VisitorStats"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "github_com_ScMofeoluwa_minalytics_shared.VisitorStats": {
            "type": "object",
            "properties": {
                "time": {
                    "type": "string"
                },
                "visitors": {
                    "type": "integer"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Minalytics API",
	Description:      "Analytics API for tracking and managing app data.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
