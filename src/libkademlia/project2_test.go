package libkademlia


import (
// "net"
"fmt"
"strconv"
"testing"
"net"
"bytes"

)

// Normal Case
func TestIterativeFindNode11(t *testing.T) {
	var id ID
	node1 := NewKademliaWithId("localhost:8000", id)

	// Generate 60 nodes
	node_set := make([][]*Kademlia, 3)
	for i := 0; i < 3; i++ {
		node_set[i] = make([]*Kademlia, 20)
		for j := 0; j < 20; j++ {
			address := "localhost:" + strconv.Itoa(8000+i*20+(j+1))
			tempID := node1.SelfContact.NodeID
			tempID[19] = uint8(32 + i * 20 + j)
			node_set[i][j] = NewKademliaWithId(address, tempID)
		}
	}
	// Initialize node1's k-bucket
	host, port, _ := StringToIpPort("localhost:8001")
	node1.DoPing(host, port)
	host, port, _ = StringToIpPort("localhost:8021")
	node1.DoPing(host, port)
	host, port, _ = StringToIpPort("localhost:8041")
	node1.DoPing(host, port)

	// Initialize three points' k bucket.
	for i := 0; i < 3; i++ {
		for j := 0; j < 20; j++ {
			address := "localhost:" + strconv.Itoa(8000+i*20+(j+1))
			host, port, _ = StringToIpPort(address)
			node_set[i][0].DoPing(host, port)
		}
	}

	keyID := node1.NodeID;
	keyID[19] = uint8(32)
	foundContacts, err := node1.DoIterativeFindNode(keyID);

	if err != nil {
		t.Error("Error: ", err)
	}

	isFound := false
	for i := 0; i < 20; i++ {
		if foundContacts[i].NodeID.Compare(node_set[0][0].NodeID) == 0{
			isFound = true
		}
	}
	if !isFound {
		t.Error("Incorrect result")
	}

	isFound = false
	for i := 0; i < 20; i++ {
		if foundContacts[i].NodeID.Compare(node_set[1][0].NodeID) == 0{
			isFound = true
		}
	}
	if !isFound {
		t.Error("Incorrect result")
	}

	isFound = false
	for i := 0; i < 20; i++ {
		if foundContacts[i].NodeID.Compare(node_set[2][0].NodeID) == 0{
			isFound = true
		}
	}
	if !isFound {
		t.Error("Incorrect result")
	}
	fmt.Println("Test 1 pass")
}

// One dead initial contact
// Contacts from k-buckects for first acive point should be returned.
func TestIterativeFindNode22(t *testing.T) {
	var id ID
	node1 := NewKademliaWithId("localhost:8200", id)

	// Generate one fake Node
	ipAddrStrings, _ := net.LookupHost("localhost")
	var host net.IP
	for i := 0; i < len(ipAddrStrings); i++ {
		host = net.ParseIP(ipAddrStrings[i])
		if host.To4() != nil {
			break
		}
	}
	tempID := node1.NodeID
	tempID[19] = uint8(32)
	fake_node := Contact{tempID, host, uint16(8201)}
	node1.Table.KBuckets[154].PushBack(fake_node)  //push it to node1's k-buckets

	// Generate 40 nodes
	node_set := make([][]*Kademlia, 3)
	for i := 0; i < 2; i++ {
		node_set[i] = make([]*Kademlia, 20)
		for j := 0; j < 20; j++ {
			address := "localhost:" + strconv.Itoa(8200+(i+1)*20+(j+1))
			tempID := node1.SelfContact.NodeID
			tempID[19] = uint8(32 + (i+ 1) * 20 + j)
			node_set[i][j] = NewKademliaWithId(address, tempID)
		}
	}

	// Initialize node1's k-bucket
	host, port, _ := StringToIpPort("localhost:8221")
	node1.DoPing(host, port)
	host, port, _ = StringToIpPort("localhost:8241")
	node1.DoPing(host, port)

	keyID := node1.NodeID;
	keyID[19] = uint8(32)
	for i := 0; i < 2; i++ {
		for j := 0; j < 20; j++ {
			address := "localhost:" + strconv.Itoa(8200+(i+1)*20+(j+1))
			host, port, _ = StringToIpPort(address)
			node_set[i][0].DoPing(host, port)
		}
	}

	foundContacts, err := node1.DoIterativeFindNode(keyID);

	if err != nil {
		t.Error("Error: ", err)
	}

	isFound := false
	for i := 0; i < 20; i++ {
		if foundContacts[i].NodeID.Compare(node_set[0][0].NodeID) == 0{
			isFound = true
		}
	}
	if !isFound {
		t.Error("Incorrect result")
	}

	isFound = false
	for i := 0; i < 20; i++ {
		if foundContacts[i].NodeID.Compare(node_set[1][0].NodeID) == 0{
			isFound = true
		}
	}
	if !isFound {
		t.Error("Incorrect result")
	}

	fmt.Println("Test 2 pass")
}

// Three dead initial contacts
// Empty contact list should be returned
func TestIterativeFindNode33(t *testing.T) {
	var id ID
	node1 := NewKademliaWithId("localhost:8300", id)

	// Generate 3 fake nodes
	ipAddrStrings, _ := net.LookupHost("localhost")
	var host net.IP
	for i := 0; i < len(ipAddrStrings); i++ {
		host = net.ParseIP(ipAddrStrings[i])
		if host.To4() != nil {
			break
		}
	}
	for i := 0; i < 3; i++ {
		tempID := node1.NodeID
		tempID[19] = uint8(32 + i)
		temp_node := Contact{tempID, host, uint16(8301 + i)}
		node1.Table.KBuckets[154].PushBack(temp_node)  //push it to node1's k-buckets
	}

	keyID := node1.NodeID;
	keyID[19] = uint8(32)
	foundContacts, err := node1.DoIterativeFindNode(keyID);

	if err == nil {
		t.Error("Error: ", err)
	}

	if len(foundContacts) != 0 {
		t.Error("Incorrect closest list of contacts.")
	}

	fmt.Println("Test 3 pass")
}

func TestIterativeStroeValue1(t *testing.T) {
	var id ID
	node1 := NewKademliaWithId("localhost:7000", id)

	// Generate 60 nodes
	node_set := make([][]*Kademlia, 3)
	for i := 0; i < 3; i++ {
		node_set[i] = make([]*Kademlia, 20)
		for j := 0; j < 20; j++ {
			address := "localhost:" + strconv.Itoa(7000+i*20+(j+1))
			tempID := node1.SelfContact.NodeID
			tempID[19] = uint8(32 + i*20 + j + 1)
			node_set[i][j] = NewKademliaWithId(address, tempID)
		}
	}
	// Initialize node1's k-bucket
	host, port, _ := StringToIpPort("localhost:7001")
	node1.DoPing(host, port)
	host, port, _ = StringToIpPort("localhost:7021")
	node1.DoPing(host, port)
	host, port, _ = StringToIpPort("localhost:7041")
	node1.DoPing(host, port)

	// Initialize three points' k bucket.
	for i := 0; i < 3; i++ {
		for j := 0; j < 20; j++ {
			address := "localhost:" + strconv.Itoa(7000+i*20+(j+1))
			host, port, _ = StringToIpPort(address)
			node_set[i][0].DoPing(host, port)
		}
	}

	key := NewRandomID()
	value := []byte("Hello world")
	// Store the value at the last node
	contacts, err := node1.DoIterativeStore(key, value)
	if err != nil {
		t.Error("Could not store value")
	}

	for i := 0; i < len(contacts); i++ {
		foundValue, _, err := node1.DoFindValue(&contacts[i], key)
		//fmt.Println(foundValue)
		if err != nil {
			t.Error("Could not find value in contacts")
		}
		if !bytes.Equal(foundValue, value) {
			t.Error("Stored value not match")
		}
	}

	fmt.Println("Test 4 pass")
}

// Normal case
// Value stored at the last node
func TestIterativeFindValue1(t *testing.T) {
	var id ID
	node1 := NewKademliaWithId("localhost:7100", id)

	// Generate 60 nodes
	node_set := make([][]*Kademlia, 3)
	for i := 0; i < 3; i++ {
		node_set[i] = make([]*Kademlia, 20)
		for j := 0; j < 20; j++ {
			address := "localhost:" + strconv.Itoa(7100+ i*20+(j+1))
			tempID := node1.SelfContact.NodeID
			tempID[19] = uint8(32 + i*20 + j + 1)
			node_set[i][j] = NewKademliaWithId(address, tempID)
		}
	}
	// Initialize node1's k-bucket
	host, port, _ := StringToIpPort("localhost:7101")
	node1.DoPing(host, port)
	host, port, _ = StringToIpPort("localhost:7121")
	node1.DoPing(host, port)
	host, port, _ = StringToIpPort("localhost:7141")
	node1.DoPing(host, port)

	// Initialize three points' k bucket.
	for i := 0; i < 3; i++ {
		for j := 0; j < 20; j++ {
			address := "localhost:" + strconv.Itoa(7100+ i*20+(j+1))
			host, port, _ = StringToIpPort(address)
			node_set[i][0].DoPing(host, port)
		}
	}

	key := NewRandomID()
	value := []byte("Hello world")
	// Store the value at the last node
	err := node_set[2][19].DoStore(&node_set[2][19].SelfContact, key, value)
	if err != nil {
		t.Error("Could not store value")
	}

	// Given the right keyID, it should return the value
	foundValue, err := node1.DoIterativeFindValue(key)
	//fmt.Println("found:", foundValue, err)
	if err != nil {
		t.Error("Cound not find value")
	}
	if !bytes.Equal(foundValue, value) {
		t.Error("Stored value did not match found value")
	}

	//Given the wrong keyID, it should return k nodes.
	wrongKey := NewRandomID()
	foundValue, err = node1.DoIterativeFindValue(wrongKey)
	if foundValue != nil {
		t.Error("Searching for a wrong ID did not return contacts")
	}

	fmt.Println("Test 5 pass")
}


func TestIterativeFindValue2(t *testing.T) {
	node1 := NewKademliaWithId("localhost:7200", NewRandomID())
	node2 := NewKademliaWithId("localhost:7201", NewRandomID())
	node3 := NewKademliaWithId("localhost:7202", NewRandomID())

	// Initialize node1's k-bucket
	host, port, _ := StringToIpPort("localhost:7201")
	node1.DoPing(host, port)
	host, port, _ = StringToIpPort("localhost:7202")
	node2.DoPing(host, port)

	key := node3.NodeID
	value := []byte("Hello wwwww")
	// Store the value at the last node
	err := node3.DoStore(&node3.SelfContact, key, value)
	if err != nil {
		t.Error("Could not store value")
	}

	// Given the right keyID, it should return the value
	foundValue, err := node1.DoIterativeFindNode(node2.NodeID)
	fmt.Println(foundValue)
	//fmt.Println("found:", foundValue, err)
	if err != nil {
		t.Error("Cound not find value")
	}
	//if !bytes.Equal(foundValue, value) {
	//	t.Error("Stored value did not match found value")
	//}
	fmt.Println("Test 6 pass")
}