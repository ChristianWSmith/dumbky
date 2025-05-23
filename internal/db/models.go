package db

type Collection struct {
	ID        int
	Name      string
	CreatedAt string
}

type Request struct {
	ID             int
	CollectionName string
	Name           string
	Payload        string
	CreatedAt      string
}

type Environment struct {
	ID        int
	Name      string
	Payload   string
	CreatedAt string
}
