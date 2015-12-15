package consul

import (
	"github.com/hashicorp/consul/consul/structs"
	"github.com/hashicorp/consul/testutil"
	"os"
	"testing"
)

func TestHealth_ChecksInState(t *testing.T) {
	dir1, s1 := testServer(t)
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()
	client := rpcClient(t, s1)
	defer client.Close()

	testutil.WaitForLeader(t, client.Call, "dc1")

	arg := structs.RegisterRequest{
		Datacenter: "dc1",
		Node:       "foo",
		Address:    "127.0.0.1",
		Check: &structs.HealthCheck{
			Name:   "memory utilization",
			Status: structs.HealthPassing,
		},
	}
	var out struct{}
	if err := client.Call("Catalog.Register", &arg, &out); err != nil {
		t.Fatalf("err: %v", err)
	}

	var out2 structs.IndexedHealthChecks
	inState := structs.ChecksInStateRequest{
		Datacenter: "dc1",
		State:      structs.HealthPassing,
	}
	if err := client.Call("Health.ChecksInState", &inState, &out2); err != nil {
		t.Fatalf("err: %v", err)
	}

	checks := out2.HealthChecks
	if len(checks) != 2 {
		t.Fatalf("Bad: %v", checks)
	}

	// First check is automatically added for the server node
	if checks[0].CheckID != SerfCheckID {
		t.Fatalf("Bad: %v", checks[0])
	}
	if checks[1].Name != "memory utilization" {
		t.Fatalf("Bad: %v", checks[1])
	}
}

func TestHealth_NodeChecks(t *testing.T) {
	dir1, s1 := testServer(t)
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()
	client := rpcClient(t, s1)
	defer client.Close()

	testutil.WaitForLeader(t, client.Call, "dc1")

	arg := structs.RegisterRequest{
		Datacenter: "dc1",
		Node:       "foo",
		Address:    "127.0.0.1",
		Check: &structs.HealthCheck{
			Name:   "memory utilization",
			Status: structs.HealthPassing,
		},
	}
	var out struct{}
	if err := client.Call("Catalog.Register", &arg, &out); err != nil {
		t.Fatalf("err: %v", err)
	}

	var out2 structs.IndexedHealthChecks
	node := structs.NodeSpecificRequest{
		Datacenter: "dc1",
		Node:       "foo",
	}
	if err := client.Call("Health.NodeChecks", &node, &out2); err != nil {
		t.Fatalf("err: %v", err)
	}

	checks := out2.HealthChecks
	if len(checks) != 1 {
		t.Fatalf("Bad: %v", checks)
	}
	if checks[0].Name != "memory utilization" {
		t.Fatalf("Bad: %v", checks)
	}
}

func TestHealth_ServiceChecks(t *testing.T) {
	dir1, s1 := testServer(t)
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()
	client := rpcClient(t, s1)
	defer client.Close()

	testutil.WaitForLeader(t, client.Call, "dc1")

	arg := structs.RegisterRequest{
		Datacenter: "dc1",
		Node:       "foo",
		Address:    "127.0.0.1",
		Service: &structs.NodeService{
			ID:      "db",
			Service: "db",
		},
		Check: &structs.HealthCheck{
			Name:      "db connect",
			Status:    structs.HealthPassing,
			ServiceID: "db",
		},
	}
	var out struct{}
	if err := client.Call("Catalog.Register", &arg, &out); err != nil {
		t.Fatalf("err: %v", err)
	}

	var out2 structs.IndexedHealthChecks
	node := structs.ServiceSpecificRequest{
		Datacenter:  "dc1",
		ServiceName: "db",
	}
	if err := client.Call("Health.ServiceChecks", &node, &out2); err != nil {
		t.Fatalf("err: %v", err)
	}

	checks := out2.HealthChecks
	if len(checks) != 1 {
		t.Fatalf("Bad: %v", checks)
	}
	if checks[0].Name != "db connect" {
		t.Fatalf("Bad: %v", checks)
	}
}

func TestHealth_ServiceNodes(t *testing.T) {
	dir1, s1 := testServer(t)
	defer os.RemoveAll(dir1)
	defer s1.Shutdown()
	client := rpcClient(t, s1)
	defer client.Close()

	testutil.WaitForLeader(t, client.Call, "dc1")

	arg := structs.RegisterRequest{
		Datacenter: "dc1",
		Node:       "foo",
		Address:    "127.0.0.1",
		Service: &structs.NodeService{
			ID:      "db",
			Service: "db",
			Tags:    []string{"master"},
		},
		Check: &structs.HealthCheck{
			Name:      "db connect",
			Status:    structs.HealthPassing,
			ServiceID: "db",
		},
	}
	var out struct{}
	if err := client.Call("Catalog.Register", &arg, &out); err != nil {
		t.Fatalf("err: %v", err)
	}

	arg = structs.RegisterRequest{
		Datacenter: "dc1",
		Node:       "bar",
		Address:    "127.0.0.2",
		Service: &structs.NodeService{
			ID:      "db",
			Service: "db",
			Tags:    []string{"slave"},
		},
		Check: &structs.HealthCheck{
			Name:      "db connect",
			Status:    structs.HealthWarning,
			ServiceID: "db",
		},
	}
	if err := client.Call("Catalog.Register", &arg, &out); err != nil {
		t.Fatalf("err: %v", err)
	}

	var out2 structs.IndexedCheckServiceNodes
	req := structs.ServiceSpecificRequest{
		Datacenter:  "dc1",
		ServiceName: "db",
		ServiceTag:  "master",
		TagFilter:   false,
	}
	if err := client.Call("Health.ServiceNodes", &req, &out2); err != nil {
		t.Fatalf("err: %v", err)
	}

	nodes := out2.Nodes
	if len(nodes) != 2 {
		t.Fatalf("Bad: %v", nodes)
	}
	if nodes[0].Node.Node != "foo" {
		t.Fatalf("Bad: %v", nodes[0])
	}
	if nodes[1].Node.Node != "bar" {
		t.Fatalf("Bad: %v", nodes[1])
	}
	if !strContains(nodes[0].Service.Tags, "master") {
		t.Fatalf("Bad: %v", nodes[0])
	}
	if !strContains(nodes[1].Service.Tags, "slave") {
		t.Fatalf("Bad: %v", nodes[1])
	}
	if nodes[0].Checks[0].Status != structs.HealthPassing {
		t.Fatalf("Bad: %v", nodes[0])
	}
	if nodes[1].Checks[0].Status != structs.HealthWarning {
		t.Fatalf("Bad: %v", nodes[1])
	}
}

func TestHealth_NodeChecks_FilterACL(t *testing.T) {
	dir, token, srv, client := testACLFilterServer(t)
	defer os.RemoveAll(dir)
	defer srv.Shutdown()
	defer client.Close()

	opt := structs.NodeSpecificRequest{
		Datacenter:   "dc1",
		Node:         srv.config.NodeName,
		QueryOptions: structs.QueryOptions{Token: token},
	}
	reply := structs.IndexedHealthChecks{}
	if err := client.Call("Health.NodeChecks", &opt, &reply); err != nil {
		t.Fatalf("err: %s", err)
	}
	found := false
	for _, chk := range reply.HealthChecks {
		switch chk.ServiceName {
		case "foo":
			found = true
		case "bar":
			t.Fatalf("bad: %#v", reply.HealthChecks)
		}
	}
	if !found {
		t.Fatalf("bad: %#v", reply.HealthChecks)
	}
}

func TestHealth_ServiceChecks_FilterACL(t *testing.T) {
	dir, token, srv, client := testACLFilterServer(t)
	defer os.RemoveAll(dir)
	defer srv.Shutdown()
	defer client.Close()

	opt := structs.ServiceSpecificRequest{
		Datacenter:   "dc1",
		ServiceName:  "foo",
		QueryOptions: structs.QueryOptions{Token: token},
	}
	reply := structs.IndexedHealthChecks{}
	if err := client.Call("Health.ServiceChecks", &opt, &reply); err != nil {
		t.Fatalf("err: %s", err)
	}
	found := false
	for _, chk := range reply.HealthChecks {
		if chk.ServiceName == "foo" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("bad: %#v", reply.HealthChecks)
	}

	opt.ServiceName = "bar"
	reply = structs.IndexedHealthChecks{}
	if err := client.Call("Health.ServiceChecks", &opt, &reply); err != nil {
		t.Fatalf("err: %s", err)
	}
	if len(reply.HealthChecks) != 0 {
		t.Fatalf("bad: %#v", reply.HealthChecks)
	}
}

func TestHealth_ServiceNodes_FilterACL(t *testing.T) {
	dir, token, srv, client := testACLFilterServer(t)
	defer os.RemoveAll(dir)
	defer srv.Shutdown()
	defer client.Close()

	opt := structs.ServiceSpecificRequest{
		Datacenter:   "dc1",
		ServiceName:  "foo",
		QueryOptions: structs.QueryOptions{Token: token},
	}
	reply := structs.IndexedCheckServiceNodes{}
	if err := client.Call("Health.ServiceNodes", &opt, &reply); err != nil {
		t.Fatalf("err: %s", err)
	}
	if len(reply.Nodes) != 1 {
		t.Fatalf("bad: %#v", reply.Nodes)
	}

	opt.ServiceName = "bar"
	reply = structs.IndexedCheckServiceNodes{}
	if err := client.Call("Health.ServiceNodes", &opt, &reply); err != nil {
		t.Fatalf("err: %s", err)
	}
	if len(reply.Nodes) != 0 {
		t.Fatalf("bad: %#v", reply.Nodes)
	}
}

func TestHealth_ChecksInState_FilterACL(t *testing.T) {
	dir, token, srv, client := testACLFilterServer(t)
	defer os.RemoveAll(dir)
	defer srv.Shutdown()
	defer client.Close()

	opt := structs.ChecksInStateRequest{
		Datacenter:   "dc1",
		State:        structs.HealthPassing,
		QueryOptions: structs.QueryOptions{Token: token},
	}
	reply := structs.IndexedHealthChecks{}
	if err := client.Call("Health.ChecksInState", &opt, &reply); err != nil {
		t.Fatalf("err: %s", err)
	}

	found := false
	for _, chk := range reply.HealthChecks {
		switch chk.ServiceName {
		case "foo":
			found = true
		case "bar":
			t.Fatalf("bad service 'bar': %#v", reply.HealthChecks)
		}
	}
	if !found {
		t.Fatalf("missing service 'foo': %#v", reply.HealthChecks)
	}
}
