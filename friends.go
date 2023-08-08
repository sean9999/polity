package main

func (node Node) Friends() []NodeAddress {
	return node.friends
}

func (node Node) SyncFriends() error {
	config := node.GetConfig()
	config.Friends = node.Friends()
	fileLocation := "./test/data/" + node.Nickname() + ".config.json"
	return config.Save(fileLocation)
}

func (node Node) AddFriend(f NodeAddress) {
	node.friends = append(node.friends, f)
	node.SyncFriends()
}
