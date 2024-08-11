package types

type ContextKey string

const IsHtmxRequestKey ContextKey = "isHtmxRequest"

type contentkey string

const (
	overview contentkey = "overview"
	board    contentkey = "board"
	member   contentkey = "member"
	files    contentkey = "files"
	reports  contentkey = "reports"
	timeline contentkey = "timeline"
)
