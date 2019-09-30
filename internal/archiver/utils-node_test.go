package archiver

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func Test_getElementsByTagName(t *testing.T) {
	htmlSource := `<div>
		<h1></h1>
		<h2></h2><h2></h2>
		<h3></h3><h3></h3><h3></h3>
		<p></p><p></p><p></p><p></p><p></p>
		<div></div><div></div><div></div><div></div><div></div>
		<div><p>Hey it's nested</p></div>
		<div></div>
		<img/><img/><img/><img/><img/><img/><img/><img/>
		<img/><img/><img/><img/>
	</div>`

	doc, err := parseHTMLSource(htmlSource)
	if err != nil {
		t.Errorf("getElementsByTagName(), failed to parse: %v", err)
	}

	tests := map[string]int{
		"h1":  1,
		"h2":  2,
		"h3":  3,
		"p":   6,
		"div": 7,
		"img": 12,
		"*":   31,
	}

	for tagName, count := range tests {
		t.Run(tagName, func(t *testing.T) {
			if got := len(getElementsByTagName(doc, tagName)); got != count {
				t.Errorf("getElementsByTagName() = %v, want %v", got, count)
			}
		})
	}
}

func Test_createElement(t *testing.T) {
	tests := []struct {
		name     string
		tagName  string
		tagCount int
	}{{
		name:     "3 headings1",
		tagName:  "h1",
		tagCount: 3,
	}, {
		name:     "4 headings2",
		tagName:  "h2",
		tagCount: 4,
	}, {
		name:     "5 headings3",
		tagName:  "h3",
		tagCount: 5,
	}, {
		name:     "10 paragraph",
		tagName:  "p",
		tagCount: 10,
	}, {
		name:     "6 div",
		tagName:  "div",
		tagCount: 6,
	}, {
		name:     "8 image",
		tagName:  "img",
		tagCount: 8,
	}, {
		name:     "22 custom tag",
		tagName:  "custom-tag",
		tagCount: 22,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &html.Node{}
			for i := 0; i < tt.tagCount; i++ {
				doc.AppendChild(createElement(tt.tagName))
			}

			if tags := getElementsByTagName(doc, tt.tagName); len(tags) != tt.tagCount {
				t.Errorf("createElement() = %v, want %v", len(tags), tt.tagCount)
			}
		})
	}
}

func Test_createTextNode(t *testing.T) {
	tests := []string{
		"hello world",
		"this is awesome",
		"all cat is good boy",
		"all dog is good boy as well",
	}

	for _, text := range tests {
		t.Run(text, func(t *testing.T) {
			node := createTextNode(text)
			if outerHTML := outerHTML(node); outerHTML != text {
				t.Errorf("createTextNode() = %v, want %v", outerHTML, text)
			}
		})
	}
}

func Test_getAttribute(t *testing.T) {
	tests := []struct {
		name       string
		htmlSource string
		attrName   string
		want       string
	}{{
		name:       "attr id from paragraph",
		htmlSource: `<p id="main-paragraph"></p>`,
		attrName:   "id",
		want:       "main-paragraph",
	}, {
		name:       "attr class from list",
		htmlSource: `<ul class="bullets"></ul>`,
		attrName:   "class",
		want:       "bullets",
	}, {
		name:       "attr style from paragraph",
		htmlSource: `<div style="display: none"></div>`,
		attrName:   "style",
		want:       "display: none",
	}, {
		name:       "attr doesn't exists",
		htmlSource: `<p id="main-paragraph"></p>`,
		attrName:   "class",
		want:       "",
	}, {
		name:       "node has no attributes",
		htmlSource: `<p></p>`,
		attrName:   "id",
		want:       "",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parseHTMLSource(tt.htmlSource)
			if err != nil {
				t.Errorf("getAttribute(), failed to parse: %v", err)
			}

			if got := getAttribute(node, tt.attrName); got != tt.want {
				t.Errorf("getAttribute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setAttribute(t *testing.T) {
	tests := []struct {
		name       string
		htmlSource string
		attrName   string
		attrValue  string
		want       string
	}{{
		name:       "set id of paragraph",
		htmlSource: `<p id="main-paragraph"></p>`,
		attrName:   "id",
		attrValue:  "txt-main",
		want:       `<p id="txt-main"></p>`,
	}, {
		name:       "set id from paragraph with several attrs",
		htmlSource: `<p id="main-paragraph" class="title"></p>`,
		attrName:   "id",
		attrValue:  "txt-main",
		want:       `<p id="txt-main" class="title"></p>`,
	}, {
		name:       "set new attr for paragraph",
		htmlSource: `<p></p>`,
		attrName:   "class",
		attrValue:  "title",
		want:       `<p class="title"></p>`,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parseHTMLSource(tt.htmlSource)
			if err != nil {
				t.Errorf("setAttribute(), failed to parse: %v", err)
			}

			setAttribute(node, tt.attrName, tt.attrValue)
			if outerHTML := outerHTML(node); outerHTML != tt.want {
				t.Errorf("setAttribute() = %v, want %v", outerHTML, tt.want)
			}
		})
	}
}

func Test_removeAttribute(t *testing.T) {
	tests := []struct {
		name       string
		htmlSource string
		attrName   string
		want       string
	}{{
		name:       "remove id of paragraph",
		htmlSource: `<p id="main-paragraph"></p>`,
		attrName:   "id",
		want:       `<p></p>`,
	}, {
		name:       "remove id from paragraph with several attrs",
		htmlSource: `<p id="main-paragraph" class="title"></p>`,
		attrName:   "id",
		want:       `<p class="title"></p>`,
	}, {
		name:       "remove inexist attr of paragraph",
		htmlSource: `<p id="main-paragraph"></p>`,
		attrName:   "class",
		want:       `<p id="main-paragraph"></p>`,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parseHTMLSource(tt.htmlSource)
			if err != nil {
				t.Errorf("removeAttribute(), failed to parse: %v", err)
			}

			removeAttribute(node, tt.attrName)
			if outerHTML := outerHTML(node); outerHTML != tt.want {
				t.Errorf("removeAttribute() = %v, want %v", outerHTML, tt.want)
			}
		})
	}
}

func Test_hasAttribute(t *testing.T) {
	tests := []struct {
		name       string
		htmlSource string
		attrName   string
		want       bool
	}{{
		name:       "attribute is exist",
		htmlSource: `<p id="main-paragraph"></p>`,
		attrName:   "id",
		want:       true,
	}, {
		name:       "attribute is not exist",
		htmlSource: `<p id="main-paragraph"></p>`,
		attrName:   "class",
		want:       false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parseHTMLSource(tt.htmlSource)
			if err != nil {
				t.Errorf("hasAttribute(), failed to parse: %v", err)
			}

			if got := hasAttribute(node, tt.attrName); got != tt.want {
				t.Errorf("hasAttribute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_textContent(t *testing.T) {
	tests := []struct {
		name       string
		htmlSource string
		want       string
	}{{
		name:       "ordinary text node",
		htmlSource: "this is an ordinary text",
		want:       "this is an ordinary text",
	}, {
		name:       "single empty node element",
		htmlSource: "<p></p>",
		want:       "",
	}, {
		name:       "single node with content",
		htmlSource: "<p>Hello all</p>",
		want:       "Hello all",
	}, {
		name:       "single node with content and unnecessary space",
		htmlSource: "<p>Hello all   </p>",
		want:       "Hello all   ",
	}, {
		name:       "nested element",
		htmlSource: "<div><p>Some nested element</p></div>",
		want:       "Some nested element",
	}, {
		name:       "nested element with unnecessary space",
		htmlSource: "<div><p>Some nested element</p>    </div>",
		want:       "Some nested element    ",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parseHTMLSource(tt.htmlSource)
			if err != nil {
				t.Errorf("textContent(), failed to parse: %v", err)
			}

			if got := textContent(node); got != tt.want {
				t.Errorf("textContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_outerHTML(t *testing.T) {
	tests := []struct {
		name       string
		htmlSource string
	}{{
		name:       "text node",
		htmlSource: "this is an ordinary text",
	}, {
		name:       "single element",
		htmlSource: "<h1>Hello</h1>",
	}, {
		name:       "nested elements",
		htmlSource: "<div><p>Some nested element</p></div>",
	}, {
		name:       "triple nested elements",
		htmlSource: "<div><p>Some <a>nested</a> element</p></div>",
	}, {
		name:       "mixed nested elements",
		htmlSource: "<div><p>Some <a>nested</a> element</p><p>and more</p></div>",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parseHTMLSource(tt.htmlSource)
			if err != nil {
				t.Errorf("outerHTML(), failed to parse: %v", err)
			}

			if got := outerHTML(node); got != tt.htmlSource {
				t.Errorf("outerHTML() = %v, want %v", got, tt.htmlSource)
			}
		})
	}
}

func Test_innerHTML(t *testing.T) {
	tests := []struct {
		name       string
		htmlSource string
		want       string
	}{{
		name:       "text node",
		htmlSource: "this is an ordinary text",
		want:       "",
	}, {
		name:       "single element",
		htmlSource: "<h1>Hello</h1>",
		want:       "Hello",
	}, {
		name:       "nested elements",
		htmlSource: "<div><p>Some nested element</p></div>",
		want:       "<p>Some nested element</p>",
	}, {
		name:       "mixed text and element node",
		htmlSource: "<div><p>Some element</p>with text</div>",
		want:       "<p>Some element</p>with text",
	}, {
		name:       "triple nested elements",
		htmlSource: "<div><p>Some <a>nested</a> element</p></div>",
		want:       "<p>Some <a>nested</a> element</p>",
	}, {
		name:       "mixed nested elements",
		htmlSource: "<div><p>Some <a>nested</a> element</p><p>and more</p></div>",
		want:       "<p>Some <a>nested</a> element</p><p>and more</p>",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parseHTMLSource(tt.htmlSource)
			if err != nil {
				t.Errorf("innerHTML(), failed to parse: %v", err)
			}

			if got := innerHTML(node); got != tt.want {
				t.Errorf("innerHTML() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_id(t *testing.T) {
	tests := []struct {
		name       string
		htmlSource string
		want       string
	}{{
		name:       "id exists",
		htmlSource: `<p id="main-paragraph"></p>`,
		want:       "main-paragraph",
	}, {
		name:       "id doesn't exist",
		htmlSource: `<p></p>`,
		want:       "",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parseHTMLSource(tt.htmlSource)
			if err != nil {
				t.Errorf("id(), failed to parse: %v", err)
			}

			if got := id(node); got != tt.want {
				t.Errorf("id() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_className(t *testing.T) {
	tests := []struct {
		name       string
		htmlSource string
		want       string
	}{{
		name:       "class doesn't exist",
		htmlSource: `<p></p>`,
		want:       "",
	}, {
		name:       "class exist",
		htmlSource: `<p class="title"></p>`,
		want:       "title",
	}, {
		name:       "multiple class",
		htmlSource: `<p class="title heading"></p>`,
		want:       "title heading",
	}, {
		name:       "multiple class with unnecessary space",
		htmlSource: `<p class="    title heading    "></p>`,
		want:       "title heading",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parseHTMLSource(tt.htmlSource)
			if err != nil {
				t.Errorf("className(), failed to parse: %v", err)
			}

			if got := className(node); got != tt.want {
				t.Errorf("className() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_children(t *testing.T) {
	tests := []struct {
		name       string
		htmlSource string
		want       []string
	}{{
		name:       "has no children",
		htmlSource: "<div></div>",
		want:       []string{},
	}, {
		name:       "has one children",
		htmlSource: "<div><p>Hello</p></div>",
		want:       []string{"<p>Hello</p>"},
	}, {
		name:       "has many children",
		htmlSource: "<div><p>Hello</p><p>I'm</p><p>Happy</p></div>",
		want:       []string{"<p>Hello</p>", "<p>I&#39;m</p>", "<p>Happy</p>"},
	}, {
		name:       "has nested children",
		htmlSource: "<div><p>Hello I'm <span>Happy</span></p></div>",
		want:       []string{"<p>Hello I&#39;m <span>Happy</span></p>"},
	}, {
		name:       "mixed text and element node",
		htmlSource: "<div><p>Hello I'm</p>happy</div>",
		want:       []string{"<p>Hello I&#39;m</p>"},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parseHTMLSource(tt.htmlSource)
			if err != nil {
				t.Errorf("children(), failed to parse: %v", err)
			}

			nodes := children(node)
			if len(nodes) != len(tt.want) {
				t.Errorf("children() count = %v, want = %v", len(nodes), len(tt.want))
			}

			for i, child := range nodes {
				wantHTML := tt.want[i]
				childHTML := outerHTML(child)
				if childHTML != wantHTML {
					t.Errorf("children() = %v, want = %v", childHTML, wantHTML)
				}
			}
		})
	}
}

func Test_childNodes(t *testing.T) {
	tests := []struct {
		name       string
		htmlSource string
		want       []string
	}{{
		name:       "has no children",
		htmlSource: "<div></div>",
		want:       []string{},
	}, {
		name:       "has one children",
		htmlSource: "<div><p>Hello</p></div>",
		want:       []string{"<p>Hello</p>"},
	}, {
		name:       "has many children",
		htmlSource: "<div><p>Hello</p><p>I'm</p><p>Happy</p></div>",
		want:       []string{"<p>Hello</p>", "<p>I&#39;m</p>", "<p>Happy</p>"},
	}, {
		name:       "has nested children",
		htmlSource: "<div><p>Hello I'm <span>Happy</span></p></div>",
		want:       []string{"<p>Hello I&#39;m <span>Happy</span></p>"},
	}, {
		name:       "mixed text and element node",
		htmlSource: "<div><p>Hello I'm</p>happy</div>",
		want:       []string{"<p>Hello I&#39;m</p>", "happy"},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parseHTMLSource(tt.htmlSource)
			if err != nil {
				t.Errorf("childNodes(), failed to parse: %v", err)
			}

			nodes := childNodes(node)
			if len(nodes) != len(tt.want) {
				t.Errorf("childNodes() count = %v, want = %v", len(nodes), len(tt.want))
			}

			for i, child := range nodes {
				wantHTML := tt.want[i]
				childHTML := outerHTML(child)
				if child.Type == html.TextNode {
					childHTML = textContent(child)
				}

				if childHTML != wantHTML {
					t.Errorf("childNodes() = %v, want = %v", childHTML, wantHTML)
				}
			}
		})
	}
}

func Test_firstElementChild(t *testing.T) {
	tests := []struct {
		name       string
		htmlSource string
		want       string
	}{{
		name:       "has no children",
		htmlSource: "<div></div>",
		want:       "",
	}, {
		name:       "has one children",
		htmlSource: "<div><p>Hey</p></div>",
		want:       "<p>Hey</p>",
	}, {
		name:       "has several children",
		htmlSource: "<div><p>Hey</p><b>bro</b></div>",
		want:       "<p>Hey</p>",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parseHTMLSource(tt.htmlSource)
			if err != nil {
				t.Errorf("firstElementChild(), failed to parse: %v", err)
			}

			if got := outerHTML(firstElementChild(node)); got != tt.want {
				t.Errorf("firstElementChild() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nextElementSibling(t *testing.T) {
	tests := []struct {
		name       string
		htmlSource string
		want       string
	}{{
		name:       "has no sibling",
		htmlSource: "<div></div>",
		want:       "",
	}, {
		name:       "has directly element sibling",
		htmlSource: "<div></div><p>Hey</p>",
		want:       "<p>Hey</p>",
	}, {
		name:       "has no element sibling",
		htmlSource: "<div></div>I'm your sibling, you know",
		want:       "",
	}, {
		name:       "has distant element sibling",
		htmlSource: "<div></div>I'm your sibling as well <p>only me matter</p>",
		want:       "<p>only me matter</p>",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := parseHTMLSource(tt.htmlSource)
			if err != nil {
				t.Errorf("nextElementSibling(), failed to parse: %v", err)
			}

			if got := outerHTML(nextElementSibling(node)); got != tt.want {
				t.Errorf("nextElementSibling() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_appendChild(t *testing.T) {
	// Child is from inside document
	t.Run("child from existing node", func(t *testing.T) {
		htmlSource := `<div><p>Lonely word</p><span>new friend</span></div>`
		want := `<div><p>Lonely word<span>new friend</span></p></div>`

		doc, err := parseHTMLSource(htmlSource)
		if err != nil {
			t.Errorf("appendChild(), failed to parse: %v", err)
		}

		p := getElementsByTagName(doc, "p")[0]
		span := getElementsByTagName(doc, "span")[0]

		appendChild(p, span)
		if got := outerHTML(doc); got != want {
			t.Errorf("appendChild() = %v, want %v", got, want)
		}
	})

	// Child is new element
	t.Run("child is new element", func(t *testing.T) {
		htmlSource := `<div><p>Lonely word</p><span>new friend</span></div>`
		want := `<div><p>Lonely word<span></span></p><span>new friend</span></div>`

		doc, err := parseHTMLSource(htmlSource)
		if err != nil {
			t.Errorf("appendChild(), failed to parse: %v", err)
		}

		p := getElementsByTagName(doc, "p")[0]
		newChild := createElement("span")

		appendChild(p, newChild)
		if got := outerHTML(doc); got != want {
			t.Errorf("appendChild() = %v, want %v", got, want)
		}
	})
}

func Test_prependChild(t *testing.T) {
	// Child is from inside document
	t.Run("child from existing node", func(t *testing.T) {
		htmlSource := `<div><p>Lonely word</p><span>new friend</span></div>`
		want := `<div><p><span>new friend</span>Lonely word</p></div>`

		doc, err := parseHTMLSource(htmlSource)
		if err != nil {
			t.Errorf("prependChild(), failed to parse: %v", err)
		}

		p := getElementsByTagName(doc, "p")[0]
		span := getElementsByTagName(doc, "span")[0]

		prependChild(p, span)
		if got := outerHTML(doc); got != want {
			t.Errorf("prependChild() = %v, want %v", got, want)
		}
	})

	// Child is new element
	t.Run("child is new element", func(t *testing.T) {
		htmlSource := `<div><p>Lonely word</p><span>new friend</span></div>`
		want := `<div><p><span></span>Lonely word</p><span>new friend</span></div>`

		doc, err := parseHTMLSource(htmlSource)
		if err != nil {
			t.Errorf("prependChild(), failed to parse: %v", err)
		}

		p := getElementsByTagName(doc, "p")[0]
		newChild := createElement("span")

		prependChild(p, newChild)
		if got := outerHTML(doc); got != want {
			t.Errorf("prependChild() = %v, want %v", got, want)
		}
	})
}

func Test_replaceNode(t *testing.T) {
	// new node is from existing element
	t.Run("new node from existing element", func(t *testing.T) {
		htmlSource := `<div><p>Lonely word</p><span>new friend</span></div>`
		want := `<div><span>new friend</span></div>`

		doc, err := parseHTMLSource(htmlSource)
		if err != nil {
			t.Errorf("replaceNode(), failed to parse: %v", err)
		}

		p := getElementsByTagName(doc, "p")[0]
		span := getElementsByTagName(doc, "span")[0]

		replaceNode(p, span)
		if got := outerHTML(doc); got != want {
			t.Errorf("replaceNode() = %v, want %v", got, want)
		}
	})

	// new node is new element
	t.Run("new node is new element", func(t *testing.T) {
		htmlSource := `<div><p>Lonely word</p><span>new friend</span></div>`
		want := `<div><span></span><span>new friend</span></div>`

		doc, err := parseHTMLSource(htmlSource)
		if err != nil {
			t.Errorf("replaceNode(), failed to parse: %v", err)
		}

		p := getElementsByTagName(doc, "p")[0]
		newNode := createElement("span")

		replaceNode(p, newNode)
		if got := outerHTML(doc); got != want {
			t.Errorf("replaceNode() = %v, want %v", got, want)
		}
	})
}

func Test_includeNode(t *testing.T) {
	htmlSource := `<div>
		<h1></h1><h2></h2><h3></h3>
		<p></p><div></div><img/><img/>
	</div>`

	doc, err := parseHTMLSource(htmlSource)
	if err != nil {
		t.Errorf("includeNode(), failed to parse: %v", err)
	}

	allElements := getElementsByTagName(doc, "*")
	h1 := getElementsByTagName(doc, "h1")[0]
	h2 := getElementsByTagName(doc, "h2")[0]
	h3 := getElementsByTagName(doc, "h3")[0]
	p := getElementsByTagName(doc, "p")[0]
	div := getElementsByTagName(doc, "div")[0]
	img := getElementsByTagName(doc, "img")[0]
	span := createElement("span")

	tests := []struct {
		name string
		node *html.Node
		want bool
	}{
		{"h1", h1, true},
		{"h2", h2, true},
		{"h3", h3, true},
		{"p", p, true},
		{"div", div, true},
		{"img", img, true},
		{"span", span, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := includeNode(allElements, tt.node); got != tt.want {
				t.Errorf("includeNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cloneNode(t *testing.T) {
	tests := []struct {
		name       string
		htmlSource string
		want       string
	}{{
		name:       "single div",
		htmlSource: "<div></div>",
	}, {
		name:       "div with one children",
		htmlSource: "<div><p>Hello</p></div>",
	}, {
		name:       "div with many children",
		htmlSource: "<div><p>Hello</p><p>I'm</p><p>Happy</p></div>",
		want:       "<div><p>Hello</p><p>I&#39;m</p><p>Happy</p></div>",
	}, {
		name:       "div with nested children",
		htmlSource: "<div><p>Hello I'm <span>Happy</span></p></div>",
		want:       "<div><p>Hello I&#39;m <span>Happy</span></p></div>",
	}, {
		name:       "div with mixed text and element node",
		htmlSource: "<div><p>Hello I'm</p>happy</div>",
		want:       "<div><p>Hello I&#39;m</p>happy</div>",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := tt.want
			if want == "" {
				want = tt.htmlSource
			}

			node, err := parseHTMLSource(tt.htmlSource)
			if err != nil {
				t.Errorf("cloneNode(), failed to parse: %v", err)
			}

			clone := cloneNode(node)
			if got := outerHTML(clone); got != want {
				t.Errorf("cloneNode() = %v, want %v", got, want)
			}
		})
	}
}

func Test_getAllNodesWithTag(t *testing.T) {
	htmlSource := `<div>
		<h1></h1>
		<h2></h2><h2></h2>
		<h3></h3><h3></h3><h3></h3>
		<p></p><p></p><p></p><p></p><p></p>
		<div></div><div></div><div></div><div></div><div></div>
		<div><p>Hey it's nested</p></div>
		<div></div>
		<img/><img/><img/><img/><img/><img/><img/><img/>
		<img/><img/><img/><img/>
	</div>`

	doc, err := parseHTMLSource(htmlSource)
	if err != nil {
		t.Errorf("getAllNodesWithTag(), failed to parse: %v", err)
	}

	tests := []struct {
		name string
		tags []string
		want int
	}{{
		name: "h1",
		tags: []string{"h1"},
		want: 1,
	}, {
		name: "h1,h2",
		tags: []string{"h1", "h2"},
		want: 3,
	}, {
		name: "h1,h2,h3",
		tags: []string{"h1", "h2", "h3"},
		want: 6,
	}, {
		name: "p",
		tags: []string{"p"},
		want: 6,
	}, {
		name: "p,span",
		tags: []string{"p", "span"},
		want: 6,
	}, {
		name: "div,img",
		tags: []string{"div", "img"},
		want: 19,
	}, {
		name: "span",
		tags: []string{"span"},
		want: 0,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(getAllNodesWithTag(doc, tt.tags...)); got != tt.want {
				t.Errorf("getAllNodesWithTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_removeNodes(t *testing.T) {
	htmlSource := `<div><h1></h1><h1></h1><p></p><img/></div>`

	tests := []struct {
		name   string
		want   string
		filter func(*html.Node) bool
	}{{
		name:   "remove all",
		want:   "<div></div>",
		filter: nil,
	}, {
		name: "remove one tag",
		want: "<div><p></p><img/></div>",
		filter: func(n *html.Node) bool {
			return tagName(n) == "h1"
		},
	}, {
		name: "remove several tags",
		want: "<div><img/></div>",
		filter: func(n *html.Node) bool {
			tag := tagName(n)
			return tag == "h1" || tag == "p"
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := parseHTMLSource(htmlSource)
			if err != nil {
				t.Errorf("removeNodes(), failed to parse: %v", err)
			}

			elements := getElementsByTagName(doc, "*")
			removeNodes(elements, tt.filter)

			if got := outerHTML(doc); got != tt.want {
				t.Errorf("removeNodes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setTextContent(t *testing.T) {
	textContent := "XXX"
	expectedResult := "<div>" + textContent + "</div>"

	tests := []struct {
		name       string
		htmlSource string
	}{{
		name:       "single div",
		htmlSource: "<div></div>",
	}, {
		name:       "div with one children",
		htmlSource: "<div><p>Hello</p></div>",
	}, {
		name:       "div with many children",
		htmlSource: "<div><p>Hello</p><p>I'm</p><p>Happy</p></div>",
	}, {
		name:       "div with nested children",
		htmlSource: "<div><p>Hello I'm <span>Happy</span></p></div>",
	}, {
		name:       "div with mixed text and element node",
		htmlSource: "<div><p>Hello I'm</p>happy</div>",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root, err := parseHTMLSource(tt.htmlSource)
			if err != nil {
				t.Errorf("setTextContent(), failed to parse: %v", err)
			}

			setTextContent(root, textContent)
			if got := outerHTML(root); got != expectedResult {
				t.Errorf("setTextContent() = %v, want %v", got, expectedResult)
			}
		})
	}
}

func parseHTMLSource(htmlSource string) (*html.Node, error) {
	doc, err := html.Parse(strings.NewReader(htmlSource))
	if err != nil {
		return nil, err
	}

	body := getElementsByTagName(doc, "body")[0]
	return body.FirstChild, nil
}
