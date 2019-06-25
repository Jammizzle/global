package mongo

type MongoCollection interface {
	// Events
	AssureeID() error
	BeforeSave() error
}
