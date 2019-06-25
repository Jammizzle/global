// Neo4J base
//
// WARNING: This needs some extra stuff https://github.com/neo4j/neo4j-go-driver/

package neo4j

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/sirupsen/logrus"
)

// Connection details from env vars or a .env file
var (
	User     string = os.Getenv("GRAPH_USERNAME")
	Password string = os.Getenv("GRAPH_PASSWORD")
	Host     string = os.Getenv("GRAPH_HOST")
	Bolt     string = "bolt"
	Port     int    = 7687
)

type Neo4JDriver struct {
	client  *neo4j.Driver
	session *neo4j.Session
}

func (n *Neo4JDriver) Connect() (err error) {

	useCustomConfig := func(level neo4j.LogLevel) func(config *neo4j.Config) {
		return func(config *neo4j.Config) {
			config.Log = neo4j.ConsoleLogger(level)
		}
	}

	connectionString := fmt.Sprintf("%s://%s:%d", Bolt, Host, Port)

	if n.client, err = neo4j.NewDriver(connectionString, neo4j.BasicAuth(User, Password, ""), useCustomConfig(neo4j.ERROR)); err != nil {
		return err // handle error
	}

	if n.session, err = n.client.Session(neo4j.AccessModeWrite); err != nil {
		return err
	}
}

func (n *Neo4JDriver) Disconnect() (err error) {
	err = n.client.Close()

	if err != nil {
		return err
	}

	err = n.session.Close()

	if err != nil {
		return err
	}

	return nil
}

func (n *Neo4JDriver) GetSession() neo4j.Session {
	return n.session
}

// #
// # Neo4J Logger
// #

type customLogger struct {
	level int
}

func (logger *customLogger) ErrorEnabled() bool {
	return ERROR <= logger.level
}

func (logger *customLogger) WarningEnabled() bool {
	return WARNING <= logger.level
}

func (logger *customLogger) InfoEnabled() bool {
	return INFO <= logger.level
}

func (logger *customLogger) DebugEnabled() bool {
	return DEBUG <= logger.level
}

func (cl *customLogger) Errorf(msg string, args ...interface{}) {
	logrus.Errorf(msg, args...)
}
func (cl *customLogger) Warningf(msg string, args ...interface{}) {
	logrus.Warnf(msg, args...)
}

func (cl *customLogger) Infof(msg string, args ...interface{}) {
	logrus.Infof(msg, args...)
}
func (cl *customLogger) Debugf(msg string, args ...interface{}) {
	logus.Debugf(msg, args...)
}
