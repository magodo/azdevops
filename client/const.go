package client

const (
	// header keys
	headerKeyAccept              = "Accept"
	headerKeyAuthorization       = "Authorization"
	headerKeyContentType         = "Content-Type"
	HeaderKeyContinuationToken   = "X-MS-ContinuationToken"
	headerKeyFedAuthRedirect     = "X-TFS-FedAuthRedirect"
	headerKeyForceMsaPassThrough = "X-VSS-ForceMsaPassThrough"
	headerKeySession             = "X-TFS-Session"
	headerUserAgent              = "User-Agent"

	// media types
	MediaTypeTextPlain       = "text/plain"
	MediaTypeApplicationJson = "application/json"
)
