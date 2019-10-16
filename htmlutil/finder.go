package htmlutil

import "golang.org/x/net/html"

// GetElementByID returns the html.Node for a specific id tag
func GetElementByID(id string, n *html.Node) (element *html.Node, ok bool) {
	for _, a := range n.Attr {
		if a.Key == "id" && a.Val == id {
			return n, true
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if element, ok = GetElementByID(id, c); ok {
			return
		}
	}
	return
}

// ValueForAttribute is a utility function returning a value given a key
func ValueForAttribute(key string, n *html.Node) (string, bool) {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val, true
		}
	}
	return "", false
}
