package main

func barfOn(err error) {
	if err != nil {
		panic(err)
	}
}
