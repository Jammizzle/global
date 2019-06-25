package mongoconnector

type MongoCollection interface {
	// Events
	AssureeID() error
	BeforeSave() error
}
