package mongo

import (
	"github.com/hellomd/middlewares/config"
	"github.com/hellomd/middlewares/random"
	mgo "gopkg.in/mgo.v2"
)

var testSession *mgo.Session

// TestDB -
type TestDB struct {
	DB *mgo.Database
}

// NewTestDB -
func NewTestDB() *TestDB {
	if testSession == nil {
		var err error
		testSession, err = mgo.Dial(config.Get(MongoCfgKey))
		if err != nil {
			panic(err)
		}
	}

	return &TestDB{
		testSession.Copy().DB(random.String(10)),
	}
}

// Close -
func (tdb *TestDB) Close() {
	tdb.DB.DropDatabase()
	tdb.DB.Session.Close()
}
