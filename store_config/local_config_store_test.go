package store_config

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/net/context"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct {
	tmpdir string
}

var _ = Suite(&MySuite{})

func (s *MySuite) SetUpTest(c *C) {
	var err error
	s.tmpdir, err = ioutil.TempDir(os.TempDir(), "openinstrument-test")
	c.Assert(err, IsNil)

	// Make a temporary config file for each test
	func() {
		log.Println("Copying file")
		in, err := os.Open("testdata/config.txt")
		c.Assert(err, IsNil)
		defer in.Close()
		out, err := os.Create(filepath.Join(s.tmpdir, "config.txt"))
		c.Assert(err, IsNil)
		defer out.Close()
		_, err = io.Copy(out, in)
		cerr := out.Close()
		c.Assert(cerr, IsNil)
		log.Println("Done copying file")
	}()

	cs := NewLocalConfigStore(filepath.Join(s.tmpdir, "config.txt"), "server1:8022")
	cs.Start(context.Background())
	Set(cs)
}

func (s *MySuite) TearDownTest(c *C) {
	Get().Stop()
	os.RemoveAll(s.tmpdir)
}

func (s *MySuite) TestGetThisServer(c *C) {
	store := Get()
	member := store.GetThisServer(context.Background())
	c.Check(member.Name, Equals, "server1:8022")
	c.Assert(store.Stop(), IsNil)
}

func (s *MySuite) TestGetServer(c *C) {
	store := Get()
	member, err := store.GetServer(context.Background(), "server1:8022")
	c.Assert(err, IsNil)
	c.Check(member.Name, Equals, "server1:8022")
	c.Assert(store.Stop(), IsNil)
}

func (s *MySuite) TestDeleteServer(c *C) {
	store := Get()

	member, err := store.GetServer(context.Background(), "server1:8022")
	c.Assert(err, IsNil)
	c.Check(member.Name, Equals, "server1:8022")

	c.Check(store.DeleteServer(context.Background(), member.Name), IsNil)
	c.Assert(err, IsNil)

	_, err = store.GetServer(context.Background(), "server1:8022")
	c.Assert(err, Not(IsNil))
}

func (s *MySuite) TestUpdateServer(c *C) {
	store := Get()

	member, err := store.GetServer(context.Background(), "server1:8022")
	c.Assert(err, IsNil)
	c.Check(member.Name, Equals, "server1:8022")

	member.Name = "server1:8322"
	c.Check(store.UpdateServer(context.Background(), member), IsNil)

	newMember, err := store.GetServer(context.Background(), "server1:8322")
	c.Assert(err, IsNil)
	c.Check(newMember.Name, Equals, "server1:8322")

	c.Assert(store.Stop(), IsNil)
}

func (s *MySuite) TestWatch(c *C) {
	store := Get()

	member, err := store.GetServer(context.Background(), "server1:8022")
	c.Assert(err, IsNil)

	watcher, err := store.SubscribeChanges()
	c.Assert(err, IsNil)

	member.Name = "server1:8322"
	c.Check(store.UpdateServer(context.Background(), member), IsNil)
	<-watcher.Chan()

	c.Check(store.UpdateServer(context.Background(), member), IsNil)
	<-watcher.Chan()

	store.UnsubscribeChanges(watcher)
	c.Check(store.UpdateServer(context.Background(), member), IsNil)
	c.Assert(store.Stop(), IsNil)
	select {
	case <-watcher.Chan():
		log.Println("Got an update from subscriber when none was expected")
		c.Fail()
	default:
	}
}
