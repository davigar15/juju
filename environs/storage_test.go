// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package environs_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	gc "launchpad.net/gocheck"

	"launchpad.net/juju-core/environs"
	"launchpad.net/juju-core/provider/dummy"
	"launchpad.net/juju-core/testing"
	"launchpad.net/juju-core/utils"
)

var _ = gc.Suite(&datasourceSuite{})

type datasourceSuite struct {
	home    *testing.FakeHome
	storage environs.Storage
	baseURL string
}

func (s *datasourceSuite) SetUpTest(c *gc.C) {
	s.home = testing.MakeFakeHome(c, existingEnv, "existing")
	environ, err := environs.PrepareFromName("test")
	c.Assert(err, gc.IsNil)
	s.storage = environ.Storage()
	s.baseURL, err = s.storage.URL("")
	c.Assert(err, gc.IsNil)
}

func (s *datasourceSuite) TearDownTest(c *gc.C) {
	dummy.Reset()
	s.home.Restore()
}

func (s *datasourceSuite) TestFetch(c *gc.C) {
	sampleData := "hello world"
	s.storage.Put("foo/bar/data.txt", bytes.NewReader([]byte(sampleData)), int64(len(sampleData)))
	ds := environs.NewStorageSimpleStreamsDataSource(s.storage, "")
	rc, url, err := ds.Fetch("foo/bar/data.txt")
	c.Assert(err, gc.IsNil)
	defer rc.Close()
	c.Assert(url, gc.Equals, s.baseURL+"foo/bar/data.txt")
	data, err := ioutil.ReadAll(rc)
	c.Assert(data, gc.DeepEquals, []byte(sampleData))
}

func (s *datasourceSuite) TestFetchWithBasePath(c *gc.C) {
	sampleData := "hello world"
	s.storage.Put("base/foo/bar/data.txt", bytes.NewReader([]byte(sampleData)), int64(len(sampleData)))
	ds := environs.NewStorageSimpleStreamsDataSource(s.storage, "base")
	rc, url, err := ds.Fetch("foo/bar/data.txt")
	c.Assert(err, gc.IsNil)
	defer rc.Close()
	c.Assert(url, gc.Equals, s.baseURL+"base/foo/bar/data.txt")
	data, err := ioutil.ReadAll(rc)
	c.Assert(data, gc.DeepEquals, []byte(sampleData))
}

func (s *datasourceSuite) TestURL(c *gc.C) {
	ds := environs.NewStorageSimpleStreamsDataSource(s.storage, "")
	url, err := ds.URL("bar")
	c.Assert(err, gc.IsNil)
	expectedURL, _ := s.storage.URL("bar")
	c.Assert(url, gc.Equals, expectedURL)
}

func (s *datasourceSuite) TestURLWithBasePath(c *gc.C) {
	ds := environs.NewStorageSimpleStreamsDataSource(s.storage, "base")
	url, err := ds.URL("bar")
	c.Assert(err, gc.IsNil)
	expectedURL, _ := s.storage.URL("base/bar")
	c.Assert(url, gc.Equals, expectedURL)
}

var _ = gc.Suite(&storageSuite{})

type storageSuite struct{}

type fakeStorage struct {
	getName     string
	listPrefix  string
	invokeCount int
	shouldRetry bool
}

func (s *fakeStorage) Get(name string) (io.ReadCloser, error) {
	s.getName = name
	s.invokeCount++
	return nil, fmt.Errorf("an error")
}

func (s *fakeStorage) List(prefix string) ([]string, error) {
	s.listPrefix = prefix
	s.invokeCount++
	return nil, fmt.Errorf("an error")
}

func (s *fakeStorage) URL(name string) (string, error) {
	return "", nil
}

func (s *fakeStorage) DefaultConsistencyStrategy() utils.AttemptStrategy {
	return utils.AttemptStrategy{Min: 10}
}

func (s *fakeStorage) ShouldRetry(error) bool {
	return s.shouldRetry
}

func (s *storageSuite) TestGet(c *gc.C) {
	stor := &fakeStorage{shouldRetry: true}
	attempt := utils.AttemptStrategy{Min: 5}
	environs.Get(stor, "foo", attempt)
	c.Assert(stor.getName, gc.Equals, "foo")
	c.Assert(stor.invokeCount, gc.Equals, 5)
}

func (s *storageSuite) TestDefaultGet(c *gc.C) {
	stor := &fakeStorage{shouldRetry: true}
	environs.DefaultGet(stor, "foo")
	c.Assert(stor.getName, gc.Equals, "foo")
	c.Assert(stor.invokeCount, gc.Equals, 10)
}

func (s *storageSuite) TestGetRetry(c *gc.C) {
	stor := &fakeStorage{}
	environs.DefaultGet(stor, "foo")
	c.Assert(stor.getName, gc.Equals, "foo")
	c.Assert(stor.invokeCount, gc.Equals, 1)
}

func (s *storageSuite) TestList(c *gc.C) {
	stor := &fakeStorage{shouldRetry: true}
	attempt := utils.AttemptStrategy{Min: 5}
	environs.List(stor, "foo", attempt)
	c.Assert(stor.listPrefix, gc.Equals, "foo")
	c.Assert(stor.invokeCount, gc.Equals, 5)
}

func (s *storageSuite) TestDefaultList(c *gc.C) {
	stor := &fakeStorage{shouldRetry: true}
	environs.DefaultList(stor, "foo")
	c.Assert(stor.listPrefix, gc.Equals, "foo")
	c.Assert(stor.invokeCount, gc.Equals, 10)
}

func (s *storageSuite) TestListRetry(c *gc.C) {
	stor := &fakeStorage{}
	environs.DefaultList(stor, "foo")
	c.Assert(stor.listPrefix, gc.Equals, "foo")
	c.Assert(stor.invokeCount, gc.Equals, 1)
}
