package main

import "golang.org/x/exp/slices"

func (node Node) Friends() []NodeAddress {
	allFriends := node.friends
	allFriends = slices.Compact(allFriends)
	node.friends = allFriends
	return node.friends
}

func (node Node) SyncFriends() error {
	config := node.GetConfig()
	config.Friends = node.Friends()
	config.Friends = slices.Compact(config.Friends)
	fileLocation := "./test/data/" + node.Nickname() + ".config.json"
	return config.Save(fileLocation)
}

func (node Node) AddFriend(newFriend NodeAddress) {
	if node.address != newFriend {
		if !slices.Contains(node.friends, newFriend) {
			node.friends = append(node.friends, newFriend)
			node.friends = slices.Compact(node.friends)
			node.SyncFriends()
		}
	}
}
