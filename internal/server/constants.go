package server

const ServerPrefix = "/api"

const (
	Unauthorised        = "Unauthorised"
	FailedReadingBody   = "Failed to read request body"
	FailedUnmarshalling = "Failed to unmarhsall json body"
	FailedMarshalling   = "Failed to marshall json body"
	FailedParsingUserId = "Failed to parse user id"
	FailedParsingPathId = "Failed to parse path id"
	FailedParsingLimit  = "Failed to parse limit"
	FailedParsingOffset = "Failed to parse offset"
)

const (
	IDKey  = "id"
	IDKey2 = "id2"
)
