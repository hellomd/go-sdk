package mongo

import (
	"github.com/hellomd/go-sdk/config"
	"github.com/hellomd/go-sdk/random"
	mgo "gopkg.in/mgo.v2"
)

var testSession *mgo.Session

// TestDB -
type TestDB struct {
	DB     *mgo.Database
	DBName string
}

// NewTestDB -
func NewTestDB() *TestDB {
	if testSession == nil {
		var err error
		testSession, err = mgo.Dial(config.Get(URLCfgKey))
		if err != nil {
			panic(err)
		}
	}

	dbName := random.String(10)

	return &TestDB{
		DB:     testSession.Copy().DB(dbName),
		DBName: dbName,
	}
}

// Close -
func (tdb *TestDB) Close() {
	tdb.DB.DropDatabase()
	tdb.DB.Session.Close()
}
