// GOLANG HTML Parser
//
// Parse HTML very easy.
//
// Example:
//  HTML := "..."
//  doc, err := html.Parse(strings.NewReader(HTML))
//  // handle error ...
//
//  fmt.Println(doc.HTML())
//
//  element := doc.Find(&html.Match{
//  	Name: "p", Attributes: map[string]string{"id": "pid"},
//  })
//
//  fmt.Println(element.Text())
package html

import (
	"bytes"
	"io"
	"strings"

	"golang.org/x/net/html"
)

// HTML Parser
type HTMLParser struct {
	root *html.Node
}

// Parse html reader
func Parse(r io.Reader) (*HTMLParser, error) {
	h, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	return &HTMLParser{h}, err
}

// Find returns first match
func (p HTMLParser) Find(selection *Match) *Element {
	if selection == nil {
		return nil
	}

	return selectNode(p.root, selection)
}

// FindAll returns all matches
func (p HTMLParser) FindAll(selection *Match) []*Element {
	if selection == nil {
		return nil
	}

	return selectAllNodes(p.root, selection)
}

func (p HTMLParser) FindAllFunc(selection *Match, f func(*Element)) {
	if selection == nil {
		return
	}

	for _, e := range selectAllNodes(p.root, selection) {
		f(e)
	}
}

// HTML returns root as HTML string
func (p HTMLParser) HTML() string {
	buf := bytes.Buffer{}

	if err := html.Render(&buf, p.root); err != nil {
		return err.Error()
	}

	return buf.String()
}

type Element struct {
	Node *html.Node
}

func (elem Element) Text() string {
	var buf bytes.Buffer

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
		}
		if n.FirstChild != nil {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
	}
	f(elem.Node)

	return buf.String()
}

// Attr returns attribute value
//
// returns empty string if not found
func (elem Element) Attr(name string) string {
	for _, a := range elem.Node.Attr {
		if a.Key == name {
			return a.Val
		}
	}
	return ""
}

// Clear deletes the tag from the tree of a given HTML.
func (elem *Element) Clear() {
	elem.Node.Parent.RemoveChild(elem.Node)
}

// AppendChild adds a element c as a child of elem.
//
// It will panic if c already has a parent or siblings.
//
// it is shorthand for elem.Node.AppendChild(c.Node)
func (elem *Element) AppendChild(c *Element) {
	elem.Node.AppendChild(c.Node)
}

// Find returns first match
func (elem *Element) Find(selection *Match) *Element {
	if selection == nil {
		return nil
	}

	return selectNode(elem.Node, selection)
}

// FindAll returns all matches
func (elem *Element) FindAll(selection *Match) []*Element {
	if selection == nil {
		return nil
	}

	return selectAllNodes(elem.Node, selection)
}

// HTML returns node as HTML string
func (elem Element) HTML() string {
	buf := bytes.Buffer{}

	if err := html.Render(&buf, elem.Node); err != nil {
		return err.Error()
	}

	return buf.String()
}

type Match struct {
	// Element name (e.g. head, title, a, div, ...)
	Name string

	// Tag attributes
	//
	// Example:
	//  map[string]string{"id":"search"}
	//  	-> <element id="search" ...>
	//
	//  map[string]string{"class":"nme"}
	//  	-> <element class="nme" ...>
	//
	// If you want only include attribute, pass empty string:
	//  map[string]string{"href":""}
	//  	-> <element href="*" ...>
	Attributes map[string]string

	// Parent of element
	//  <parent>
	//  	...
	//  	<element>
	Parent *Match

	// First child of element
	//  <element>
	//  	<child>
	FirstChild *Match
}

// MatchNode returns true if selection want this node
func (s *Match) MatchNode(node *html.Node) bool {
	// check tag name
	if s.Name != "" && s.Name != node.Data {
		return false
	}

	// check tag attributes
	if s.Attributes != nil {
		if node.Attr == nil {
			return false
		}

		for k, v := range s.Attributes {
			attr := getAttr(node.Attr, k)
			if attr == nil {
				return false
			}

			if v != "" {
				if attr.Val == "" || (k != "class" && v != attr.Val) {
					return false
				}

				ok := false

				for _, className := range strings.Split(attr.Val, " ") {
					if v == className {
						ok = true
						break
					}
				}

				if !ok {
					return false
				}
			}
		}
	}

	// check parent
	if s.Parent != nil {
		if node.Parent == nil {
			return false
		}

		if !s.Parent.MatchNode(node.Parent) {
			return false
		}
	}

	// check first child
	if s.FirstChild != nil {
		if node.FirstChild == nil {
			return false
		}

		if !s.FirstChild.MatchNode(node.FirstChild) {
			return false
		}
	}

	return true
}

// selectNode returns first match
func selectNode(root *html.Node, selection *Match) *Element {

	var wnode *Element = nil

	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && selection.MatchNode(node) {
			wnode = &Element{node}
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}

	crawler(root)

	return wnode
}

// selectAllNodes returns all matches
func selectAllNodes(root *html.Node, selection *Match) []*Element {

	var wnodes []*Element = nil

	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && selection.MatchNode(node) {
			wnodes = append(wnodes, &Element{node})
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}

	crawler(root)

	return wnodes
}

func getAttr(attr []html.Attribute, name string) *html.Attribute {
	for _, a := range attr {
		if a.Key == name {
			return &a
		}
	}
	return nil
}
