package types

//Error represents customized error object
//swagger:model
type Error struct {
	Path       string
	Message    string
	Error      error
	Type       string
	StatusCode int
	IsIgnore   bool
	Params     string
}
