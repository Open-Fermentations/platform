package server

const ContextUserIdKey = "user_id"

const (
	ContentTypeJSON = "application/json"
)

const ServerPrefix = "/api"

const (
	BadBodyRead         = "Failed to read request body"
	FailedToMarshall    = "Failed to marshall json body"
	FailedToParseUserId = "Failed to parse user id"
)
