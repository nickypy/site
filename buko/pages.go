package server

type NavBar struct {
	Items []NavBarItem
}

type NavBarItem struct {
	Key  string
	Link string
}
