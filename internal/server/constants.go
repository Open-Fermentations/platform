package server

const ContextUserIdKey = "user_id"

const ServerPrefix = "/api"

const (
	BadBodyRead         = "Failed to read request body"
	FailedToUnmarshall  = "Failed to unmarhsall json body"
	FailedToMarshall    = "Failed to marshall json body"
	FailedToParseUserId = "Failed to parse user id"
	FailedToParsePathId = "Failed to parse path id"
	FailedParsingLimit  = "Failed to parse limit"
	FailedParsingOffset = "Failed to parse offset"
)

const (
	IDKey = "id"
)
