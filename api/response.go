package main

import "github.com/ScMofeoluwa/minalytics/types"

type APIResponse struct {
	Data       interface{} `json:"data,omitempty"`
	statusCode int         `json:"-"`
	Message    string      `json:"message"`
}

type APIStatus struct {
	statusCode int    `json:"-"`
	Message    string `json:"message"`
}

type ReferralResponse struct {
	Data types.ReferralStats
	APIStatus
}

type PageResponse struct {
	Data types.PageStats
	APIStatus
}

type BrowserResponse struct {
	Data types.BrowserStats
	APIStatus
}

type CountryResponse struct {
	Data types.CountryStats
	APIStatus
}

type DeviceResponse struct {
	Data types.DeviceStats
	APIStatus
}

type OSResponse struct {
	Data types.OSStats
	APIStatus
}

func NewSuccessResponse(data interface{}, code int, message string) APIResponse {
	return APIResponse{
		Data:       data,
		statusCode: code,
		Message:    message,
	}
}

func NewErrorResponse(code int, message string) APIResponse {
	return APIResponse{
		statusCode: code,
		Message:    message,
	}
}
