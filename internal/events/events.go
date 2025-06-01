package events

type RequestSelected struct {
	RequestName string
	IsSelected  bool
}

type CollectionSelected struct {
	CollectionName string
	IsSelected     bool
}

type RequestSaved struct{}
