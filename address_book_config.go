package polity

// type addressMapConfig map[string]AddressString
// type addressBookConfig map[string]addressMapConfig

// func addressToString(addr net.Addr) AddressString {
// 	str := fmt.Sprintf("%s://%s", addr.Network(), addr.String())
// 	return AddressString(str)
// }

// func (am AddressMap) Config() addressMapConfig {
// 	conf := make(addressMapConfig, len(am))
// 	for ns, addr := range am {
// 		conf[ns] = addr
// 	}
// 	return conf
// }

// func (ab AddressBook) Config() addressBookConfig {
// 	conf := make(addressBookConfig, len(ab))
// 	for p, addrmap := range ab {
// 		hexkey := fmt.Sprintf("%x", p[:])
// 		conf[hexkey] = addrmap.Config()
// 	}
// 	return conf
// }
