package feed

// CheerLights TODO
type CheerLights struct {
	Channel Channel
	Feeds   []Feed
}

// Channel TODO
type Channel struct {
	ID          int
	Name        string
	Description string
	Latitude    string
	Longitude   string
	Field1      string
	Field2      string
	CreatedAt   string
	UpdatedAt   string
	LastEntryID int
}

// Feed TODO
type Feed struct {
	CreatedAt string
	EntryID   int
	Field1    string
	Field2    string
}
