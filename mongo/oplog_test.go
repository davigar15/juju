// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package mongo_test

import (
	"reflect"
	"time"

	jujutesting "github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/juju/juju/mongo"
	coretesting "github.com/juju/juju/testing"
	"github.com/juju/juju/worker/peergrouper"
)

type oplogSuite struct {
	coretesting.BaseSuite
}

var _ = gc.Suite(&oplogSuite{})

func (s *oplogSuite) TestWithRealOplog(c *gc.C) {
	_, session := s.startMongoWithReplicaset(c)

	// Watch for oplog entries for the "bar" collection in the "foo"
	// DB.
	oplog := mongo.GetOplog(session)
	tailer := mongo.NewOplogTailer(
		oplog,
		bson.D{{"ns", "foo.bar"}},
		time.Now().Add(-time.Minute),
	)
	defer tailer.Stop()

	assertOplog := func(expectedOp string, expectedObj, expectedUpdate bson.D) {
		doc := s.getNextOplog(c, tailer)
		c.Assert(doc.Operation, gc.Equals, expectedOp)

		var actualObj bson.D
		err := doc.UnmarshalObject(&actualObj)
		c.Assert(err, jc.ErrorIsNil)
		c.Assert(actualObj, jc.DeepEquals, expectedObj)

		var actualUpdate bson.D
		err = doc.UnmarshalUpdate(&actualUpdate)
		c.Assert(err, jc.ErrorIsNil)
		c.Assert(actualUpdate, jc.DeepEquals, expectedUpdate)
	}

	// Insert into foo.bar and see that the oplog entry is reported.
	db := session.DB("foo")
	coll := db.C("bar")
	s.insertDoc(c, session, coll, bson.M{"_id": "thing"})
	assertOplog("i", bson.D{{"_id", "thing"}}, nil)

	// Update foo.bar and see the update reported.
	err := coll.UpdateId("thing", bson.M{"$set": bson.M{"blah": 42}})
	c.Assert(err, jc.ErrorIsNil)
	assertOplog("u", bson.D{{"$set", bson.D{{"blah", 42}}}}, bson.D{{"_id", "thing"}})

	// Insert into another collection (shouldn't be reported due to filter).
	s.insertDoc(c, session, db.C("elsewhere"), bson.M{"_id": "boo"})
	s.assertNoOplog(c, tailer)
}

func (s *oplogSuite) TestHonoursInitialTs(c *gc.C) {
	_, session := s.startMongo(c)

	t := time.Now()

	oplog := s.makeFakeOplog(c, session)
	for offset := -1; offset <= 1; offset++ {
		tDoc := t.Add(time.Duration(offset) * time.Second)
		s.insertDoc(c, session, oplog,
			&mongo.OplogDoc{Timestamp: mongo.NewMongoTimestamp(tDoc)},
		)
	}

	tailer := mongo.NewOplogTailer(oplog, nil, t)
	defer tailer.Stop()

	for offset := 0; offset <= 1; offset++ {
		doc := s.getNextOplog(c, tailer)
		tExpected := t.Add(time.Duration(offset) * time.Second)
		c.Assert(doc.Timestamp, gc.Equals, mongo.NewMongoTimestamp(tExpected))
	}
}

func (s *oplogSuite) TestStops(c *gc.C) {
	_, session := s.startMongo(c)

	oplog := s.makeFakeOplog(c, session)
	tailer := mongo.NewOplogTailer(oplog, nil, time.Time{})
	defer tailer.Stop()

	s.insertDoc(c, session, oplog, &mongo.OplogDoc{Timestamp: 1})
	s.getNextOplog(c, tailer)

	err := tailer.Stop()
	c.Assert(err, jc.ErrorIsNil)

	s.assertStopped(c, tailer)
	c.Assert(tailer.Err(), jc.ErrorIsNil)
}

type fakeIterator struct {
	docs    []*mongo.OplogDoc
	pos     int
	err     error
	timeout bool
}

func mkIterator(err error, docs ...*mongo.OplogDoc) *fakeIterator {
	return &fakeIterator{docs: docs, err: err}
}

func (i *fakeIterator) Next(result interface{}) bool {
	if i.pos >= len(i.docs) {
		return false
	}
	target := reflect.ValueOf(result).Elem()
	target.Set(reflect.ValueOf(*i.docs[i.pos]))
	i.pos++
	return true
}

func (i *fakeIterator) Err() error {
	if i.pos < len(i.docs) {
		return nil
	}
	return i.err
}

func (i *fakeIterator) Timeout() bool {
	if i.pos < len(i.docs) {
		return false
	}
	return i.timeout
}

type fakeSession struct {
	iterators  []*fakeIterator
	pos        int
	timestamps []bson.MongoTimestamp
	excludeIds [][]int64
}

var timeoutIterator = fakeIterator{timeout: true}

func (s *fakeSession) NewIter(ts bson.MongoTimestamp, ids []int64) mongo.OplogIterator {
	if s.pos >= len(s.iterators) {
		// We've run out of results - at this point the calls to get
		// more data would just keep timing out.
		return &timeoutIterator
		// TODO(babbageclunk): the tailer goroutine will just keep
		// calling this in a tight loop once it's gone through the
		// canned iterators (until the tailer.Stop kills it). Is this
		// a problem? The real call would block for a second waiting
		// for new data.
	}
	s.timestamps = append(s.timestamps, ts)
	s.excludeIds = append(s.excludeIds, ids)
	result := s.iterators[s.pos]
	s.pos++
	return result
}

func (s *fakeSession) Close() {}

func mkSession(iterators ...*fakeIterator) *fakeSession {
	return &fakeSession{iterators: iterators}
}

func (s *oplogSuite) TestRestartsOnError(c *gc.C) {
	session := mkSession(
		// First iterator terminates with an ErrCursor
		mkIterator(mgo.ErrCursor, &mongo.OplogDoc{Timestamp: 1, OperationId: 99}),
		mkIterator(nil, &mongo.OplogDoc{Timestamp: 2, OperationId: 42}),
	)
	tailer := mongo.NewOplogTailerWithSession(session, time.Time{})
	defer tailer.Stop()

	// First, ensure that the tailer is seeing oplog rows and handles
	// the ErrCursor that occurs at the end.
	doc := s.getNextOplog(c, tailer)
	c.Check(doc.Timestamp, gc.Equals, bson.MongoTimestamp(1))
	c.Check(session.timestamps[0], gc.Equals, mongo.NewMongoTimestamp(time.Time{}))
	c.Check(session.excludeIds[0], gc.IsNil)

	// Ensure that the tailer continues after getting a new iterator.
	doc = s.getNextOplog(c, tailer)
	c.Check(doc.Timestamp, gc.Equals, bson.MongoTimestamp(2))
	c.Check(session.timestamps[1], gc.Equals, bson.MongoTimestamp(1))
	c.Check(session.excludeIds[1], gc.DeepEquals, []int64{99})
}

func (s *oplogSuite) TestNoRepeatsAfterIterRestart(c *gc.C) {
	// A bunch of documents with the same timestamp but different ids.
	// These will be split across 2 iterators.
	docs := make([]*mongo.OplogDoc, 11)
	for i := 0; i < 10; i++ {
		id := int64(i + 10)
		docs[i] = &mongo.OplogDoc{
			Timestamp:   1,
			OperationId: id,
		}
	}
	// Add one more with a different timestamp.
	docs[10] = &mongo.OplogDoc{
		Timestamp:   2,
		OperationId: 42,
	}
	session := mkSession(
		// First block of documents, all time 1
		mkIterator(nil, docs[:5]...),
		// Second block, some time 1, one time 2
		mkIterator(nil, docs[5:]...),
	)
	tailer := mongo.NewOplogTailerWithSession(session, time.Time{})
	defer tailer.Stop()

	for id := int64(10); id < 15; id++ {
		doc := s.getNextOplog(c, tailer)
		c.Assert(doc.Timestamp, gc.Equals, bson.MongoTimestamp(1))
		c.Assert(doc.OperationId, gc.Equals, id)
	}

	// Check the query doesn't exclude any in the first request.
	c.Check(session.timestamps[0], gc.Equals, mongo.NewMongoTimestamp(time.Time{}))
	c.Check(session.excludeIds[0], gc.IsNil)

	// The OplogTailer will fall off the end of the iterator and get a new one.

	// Ensure that only previously unreported entries are now reported.
	for id := int64(15); id < 20; id++ {
		doc := s.getNextOplog(c, tailer)
		c.Assert(doc.Timestamp, gc.Equals, bson.MongoTimestamp(1))
		c.Assert(doc.OperationId, gc.Equals, id)
	}

	// Check we got the next block correctly
	c.Check(session.timestamps[1], gc.Equals, bson.MongoTimestamp(1))
	c.Check(session.excludeIds[1], gc.DeepEquals, []int64{10, 11, 12, 13, 14})

	doc := s.getNextOplog(c, tailer)
	c.Assert(doc.Timestamp, gc.Equals, bson.MongoTimestamp(2))
	c.Assert(doc.OperationId, gc.Equals, int64(42))
}

func (s *oplogSuite) TestDiesOnFatalError(c *gc.C) {
	_, session := s.startMongo(c)
	oplog := s.makeFakeOplog(c, session)
	s.insertDoc(c, session, oplog, &mongo.OplogDoc{Timestamp: 1})

	tailer := mongo.NewOplogTailer(oplog, nil, time.Time{})
	defer tailer.Stop()

	doc := s.getNextOplog(c, tailer)
	c.Assert(doc.Timestamp, gc.Equals, bson.MongoTimestamp(1))

	// Induce a fatal error by removing the oplog collection.
	err := oplog.DropCollection()
	c.Assert(err, jc.ErrorIsNil)

	s.assertStopped(c, tailer)
	// The actual error varies by MongoDB version so just check that
	// there is one.
	c.Assert(tailer.Err(), gc.Not(jc.ErrorIsNil))
}

func (s *oplogSuite) TestNewMongoTimestamp(c *gc.C) {
	t := time.Date(2015, 6, 24, 12, 47, 0, 0, time.FixedZone("somewhere", 5*3600))

	expected := bson.MongoTimestamp(6163845091342417920)
	c.Assert(mongo.NewMongoTimestamp(t), gc.Equals, expected)
	c.Assert(mongo.NewMongoTimestamp(t.In(time.UTC)), gc.Equals, expected)
}

func (s *oplogSuite) TestNewMongoTimestampBeforeUnixEpoch(c *gc.C) {
	c.Assert(mongo.NewMongoTimestamp(time.Time{}), gc.Equals, bson.MongoTimestamp(0))
}

func (s *oplogSuite) startMongoWithReplicaset(c *gc.C) (*jujutesting.MgoInstance, *mgo.Session) {
	inst := &jujutesting.MgoInstance{
		Params: []string{
			"--replSet", "juju",
		},
	}
	err := inst.Start(nil)
	c.Assert(err, jc.ErrorIsNil)
	s.AddCleanup(func(*gc.C) { inst.Destroy() })

	// Initiate replicaset.
	info := inst.DialInfo()
	args := peergrouper.InitiateMongoParams{
		DialInfo:       info,
		MemberHostPort: inst.Addr(),
	}
	err = peergrouper.InitiateMongoServer(args)
	c.Assert(err, jc.ErrorIsNil)

	return inst, s.dialMongo(c, inst)
}

func (s *oplogSuite) startMongo(c *gc.C) (*jujutesting.MgoInstance, *mgo.Session) {
	inst := &jujutesting.MgoInstance{
		Params: []string{
			"--setParameter", "enableTestCommands=1", // allows "emptycapped" command
		},
	}
	err := inst.Start(nil)
	c.Assert(err, jc.ErrorIsNil)
	s.AddCleanup(func(*gc.C) { inst.Destroy() })
	return inst, s.dialMongo(c, inst)
}

func (s *oplogSuite) emptyCapped(c *gc.C, coll *mgo.Collection) {
	// Call the emptycapped (test) command on a capped
	// collection. This invalidates any cursors on the collection.
	err := coll.Database.Run(bson.D{{"emptycapped", coll.Name}}, nil)
	c.Assert(err, jc.ErrorIsNil)
}

func (s *oplogSuite) dialMongo(c *gc.C, inst *jujutesting.MgoInstance) *mgo.Session {
	session, err := inst.Dial()
	c.Assert(err, jc.ErrorIsNil)
	s.AddCleanup(func(*gc.C) { session.Close() })
	return session
}

func (s *oplogSuite) makeFakeOplog(c *gc.C, session *mgo.Session) *mgo.Collection {
	db := session.DB("foo")
	oplog := db.C("oplog.fake")
	err := oplog.Create(&mgo.CollectionInfo{
		Capped:   true,
		MaxBytes: 1024 * 1024,
	})
	c.Assert(err, jc.ErrorIsNil)
	return oplog
}

func (s *oplogSuite) insertDoc(c *gc.C, srcSession *mgo.Session, coll *mgo.Collection, doc interface{}) {
	session := srcSession.Copy()
	defer session.Close()
	err := coll.With(session).Insert(doc)
	c.Assert(err, jc.ErrorIsNil)
}

func (s *oplogSuite) getNextOplog(c *gc.C, tailer *mongo.OplogTailer) *mongo.OplogDoc {
	select {
	case doc, ok := <-tailer.Out():
		if !ok {
			c.Fatalf("tailer unexpectedly died: %v", tailer.Err())
		}
		return doc
	case <-time.After(coretesting.LongWait):
		c.Fatal("timed out waiting for oplog doc")
	}
	return nil
}

func (s *oplogSuite) assertNoOplog(c *gc.C, tailer *mongo.OplogTailer) {
	select {
	case _, ok := <-tailer.Out():
		if !ok {
			c.Fatalf("tailer unexpectedly died: %v", tailer.Err())
		}
		c.Fatal("unexpected oplog activity reported")
	case <-time.After(coretesting.ShortWait):
		// Success
	}
}

func (s *oplogSuite) assertStopped(c *gc.C, tailer *mongo.OplogTailer) {
	// Output should close.
	select {
	case _, ok := <-tailer.Out():
		c.Assert(ok, jc.IsFalse)
	case <-time.After(coretesting.LongWait):
		c.Fatal("tailer output should have closed")
	}

	// OplogTailer should die.
	select {
	case <-tailer.Dying():
		// Success.
	case <-time.After(coretesting.LongWait):
		c.Fatal("tailer should have died")
	}
}
