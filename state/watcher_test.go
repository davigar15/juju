package state_test

import (
	. "launchpad.net/gocheck"
	"launchpad.net/juju/go/state"
	"launchpad.net/juju/go/state/watcher"
	"time"
)

var serviceWatchConfigData = []map[string]interface{}{
	{},
	{"foo": "bar", "baz": "yadda"},
	{"baz": "yadda"},
}

func (s *StateSuite) TestServiceWatchConfig(c *C) {
	dummy := s.addDummyCharm(c)
	wordpress, err := s.st.AddService("wordpress", dummy)
	c.Assert(err, IsNil)
	c.Assert(wordpress.Name(), Equals, "wordpress")

	config, err := wordpress.Config()
	c.Assert(err, IsNil)
	c.Assert(config.Keys(), HasLen, 0)
	configWatcher := wordpress.WatchConfig()

	time.Sleep(100 * time.Millisecond)

	// Two change events.
	config.Set("foo", "bar")
	config.Set("baz", "yadda")
	changes, err := config.Write()
	c.Assert(err, IsNil)
	c.Assert(changes, DeepEquals, []state.ItemChange{{
		Key: "baz",
		Type: state.ItemAdded,
		NewValue: "yadda",
	}, {
		Key: "foo",
		Type: state.ItemAdded,
		NewValue: "bar",
	}})
	time.Sleep(100 * time.Millisecond)
	config.Delete("foo")
	changes, err = config.Write()
	c.Assert(err, IsNil)
	c.Assert(changes, DeepEquals, []state.ItemChange{{
		Key: "foo",
		Type: state.ItemDeleted,
		OldValue: "bar",
	}})

	for _, want := range serviceWatchConfigData {
		select {
		case got, ok := <-configWatcher.Changes():
			c.Logf("got configNode %p", got)
			c.Assert(ok, Equals, true)
			c.Assert(got.Map(), DeepEquals, want)
		case <-time.After(200 * time.Millisecond):
			c.Fatalf("didn't get change: %#v", want)
		}
	}

	select {
	case got, _ := <-configWatcher.Changes():
		c.Fatalf("got unexpected change: %#v", got)
	case <-time.After(100 * time.Millisecond):
	}

	err = configWatcher.Stop()
	c.Assert(err, IsNil)
}

func (s *StateSuite) TestServiceWatchConfigIllegalData(c *C) {
	dummy := s.addDummyCharm(c)
	wordpress, err := s.st.AddService("wordpress", dummy)
	c.Assert(err, IsNil)
	c.Assert(wordpress.Name(), Equals, "wordpress")
	configWatcher := wordpress.WatchConfig()

	// Receive empty change after service adding.
	select {
	case got, ok := <-configWatcher.Changes():
		c.Assert(ok, Equals, true)
		c.Assert(got.Map(), DeepEquals, map[string]interface{}{})
	case <-time.After(100 * time.Millisecond):
		c.Fatalf("unexpected timeout")
	}

	// Set config to illegal data.
	_, err = s.zkConn.Set("/services/service-0000000000/config", "---", -1)
	c.Assert(err, IsNil)

	select {
	case _, ok := <-configWatcher.Changes():
		c.Assert(ok, Equals, false)
	case <-time.After(100 * time.Millisecond):
	}

	err = configWatcher.Stop()
	c.Assert(err, ErrorMatches, "YAML error: .*")
}

type unitWatchNeedsUpgradeTest struct {
	test func(*state.Unit) error
	want state.NeedsUpgrade
}

var unitWatchNeedsUpgradeTests = []unitWatchNeedsUpgradeTest{
	{func(u *state.Unit) error { return u.SetNeedsUpgrade(false) }, state.NeedsUpgrade{true, false}},
	{func(u *state.Unit) error { return u.ClearNeedsUpgrade() }, state.NeedsUpgrade{false, false}},
	{func(u *state.Unit) error { return u.SetNeedsUpgrade(true) }, state.NeedsUpgrade{true, true}},
}

func (s *StateSuite) TestUnitWatchNeedsUpgrade(c *C) {
	dummy := s.addDummyCharm(c)
	wordpress, err := s.st.AddService("wordpress", dummy)
	c.Assert(err, IsNil)
	c.Assert(wordpress.Name(), Equals, "wordpress")
	unit, err := wordpress.AddUnit()
	c.Assert(err, IsNil)
	needsUpgradeWatcher := unit.WatchNeedsUpgrade()

	for _, test := range unitWatchNeedsUpgradeTests {
		err := test.test(unit)
		c.Assert(err, IsNil)
		select {
		case got, ok := <-needsUpgradeWatcher.Changes():
			c.Assert(ok, Equals, true)
			c.Assert(got, DeepEquals, test.want)
		case <-time.After(200 * time.Millisecond):
			c.Fatalf("didn't get change: %#v", test.want)
		}
	}

	select {
	case got, _ := <-needsUpgradeWatcher.Changes():
		c.Fatalf("got unexpected change: %#v", got)
	case <-time.After(100 * time.Millisecond):
	}

	err = needsUpgradeWatcher.Stop()
	c.Assert(err, IsNil)
}

type unitWatchResolvedTest struct {
	test func(*state.Unit) error
	want state.ResolvedMode
}

var unitWatchResolvedTests = []unitWatchResolvedTest{
	{func(u *state.Unit) error { return u.SetResolved(state.ResolvedRetryHooks) }, state.ResolvedRetryHooks},
	{func(u *state.Unit) error { return u.ClearResolved() }, state.ResolvedNone},
	{func(u *state.Unit) error { return u.SetResolved(state.ResolvedNoHooks) }, state.ResolvedNoHooks},
}

func (s *StateSuite) TestUnitWatchResolved(c *C) {
	dummy := s.addDummyCharm(c)
	wordpress, err := s.st.AddService("wordpress", dummy)
	c.Assert(err, IsNil)
	c.Assert(wordpress.Name(), Equals, "wordpress")
	unit, err := wordpress.AddUnit()
	c.Assert(err, IsNil)
	resolvedWatcher := unit.WatchResolved()

	for _, test := range unitWatchResolvedTests {
		err := test.test(unit)
		c.Assert(err, IsNil)
		select {
		case got, ok := <-resolvedWatcher.Changes():
			c.Assert(ok, Equals, true)
			c.Assert(got, Equals, test.want)
		case <-time.After(200 * time.Millisecond):
			c.Fatalf("didn't get change: %#v", test.want)
		}
	}

	select {
	case got, _ := <-resolvedWatcher.Changes():
		c.Fatalf("got unexpected change: %#v", got)
	case <-time.After(100 * time.Millisecond):
	}

	err = resolvedWatcher.Stop()
	c.Assert(err, IsNil)
}

type unitWatchPortsTest struct {
	test func(*state.Unit) error
	want []state.Port
}

var unitWatchPortsTests = []unitWatchPortsTest{
	{func(u *state.Unit) error { return u.OpenPort("tcp", 80) }, []state.Port{{"tcp", 80}}},
	{func(u *state.Unit) error { return u.OpenPort("udp", 53) }, []state.Port{{"tcp", 80}, {"udp", 53}}},
	{func(u *state.Unit) error { return u.ClosePort("tcp", 80) }, []state.Port{{"udp", 53}}},
}

func (s *StateSuite) TestUnitWatchPorts(c *C) {
	dummy := s.addDummyCharm(c)
	wordpress, err := s.st.AddService("wordpress", dummy)
	c.Assert(err, IsNil)
	c.Assert(wordpress.Name(), Equals, "wordpress")
	unit, err := wordpress.AddUnit()
	c.Assert(err, IsNil)
	portsWatcher := unit.WatchPorts()

	for _, test := range unitWatchPortsTests {
		err := test.test(unit)
		c.Assert(err, IsNil)
		select {
		case got, ok := <-portsWatcher.Changes():
			c.Assert(ok, Equals, true)
			c.Assert(got, DeepEquals, test.want)
		case <-time.After(200 * time.Millisecond):
			c.Fatalf("didn't get change: %#v", test.want)
		}
	}

	select {
	case got, _ := <-portsWatcher.Changes():
		c.Fatalf("got unexpected change: %#v", got)
	case <-time.After(100 * time.Millisecond):
	}

	err = portsWatcher.Stop()
	c.Assert(err, IsNil)
}

type machinesWatchTest struct {
	test func(*state.State) error
	want watcher.ChildrenChange
}

var machinesWatchTests = []machinesWatchTest{
	{
		func(s *state.State) error {
			_, err := s.AddMachine()
			return err
		},
		watcher.ChildrenChange{Added: []string{"machine-0000000000"}},
	},
	{
		func(s *state.State) error {
			_, err := s.AddMachine()
			return err
		},
		watcher.ChildrenChange{Added: []string{"machine-0000000001"}},
	},
	{
		func(s *state.State) error {
			return s.RemoveMachine(1)
		},
		watcher.ChildrenChange{Deleted: []string{"machine-0000000001"}},
	},
}

func (s *StateSuite) TestWatchMachines(c *C) {
	w := s.st.WatchMachines()

	for _, test := range machinesWatchTests {
		err := test.test(s.st)
		c.Assert(err, IsNil)
		select {
		case got, ok := <-w.Changes():
			c.Assert(ok, Equals, true)
			c.Assert(got, DeepEquals, test.want)
		case <-time.After(200 * time.Millisecond):
			c.Fatalf("didn't get change: %#v", test.want)
		}
	}

	select {
	case got, _ := <-w.Changes():
		c.Fatalf("got unexpected change: %#v", got)
	case <-time.After(100 * time.Millisecond):
	}

	c.Assert(w.Stop(), IsNil)
}
