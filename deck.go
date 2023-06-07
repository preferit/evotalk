package evotalk

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	. "github.com/gregoryv/web"
	"github.com/gregoryv/web/files"
)

// NewDeck returns a Deck with default styling and navigation.
func NewDeck() *Deck {
	return &Deck{
		Slides: make([]*Element, 0),
		Styles: []*CSS{DefaultTheme(), HighlightColors()},
		nav:    newNavbar(),
	}
}

// Deck represents a deck of slides, ie. a presentation.
type Deck struct {
	Title  string // header title
	Slides []*Element
	Styles []*CSS

	nav *navbar
}

// Slide appends a new slide to the deck. elements can be anything
// supported by package gregoryv/web.
func (d *Deck) Slide(elements ...interface{}) {
	d.Slides = append(d.Slides, Wrap(elements...))
	d.nav.max = len(d.Slides)
}

func (d *Deck) GroupEnd() {
	d.nav.groupEnd(len(d.Slides) + 1)
}

// Page returns a web page ready for use.
func (d *Deck) Page() *Page {
	styles := Style()
	for _, s := range d.Styles {
		styles.With(s)
	}
	body := Body()
	for i, content := range d.Slides {
		j := i + 1
		id := fmt.Sprintf("%v", j)
		slide := Div(Class("slide"), Id(id))

		slide.With(
			Header(content.Children[0]),
		)
		div := Div(Class("content"))
		div.With(content.Children[1:]...)
		slide.With(div, d.nav.next())
		body.With(slide)
	}
	body.With(Script(enhancejs))
	return NewFile("index.html",
		Html(
			Head(
				Title(d.Title),
				Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
				styles,
			),
			body,
		),
	)
}

// ----------------------------------------
// layouts

func LayoutTwoCol(v ...interface{}) *Element {
	tr := Tr()
	td := Td()
	for _, v := range v {
		switch v := v.(type) {
		case string:
			if v == "----" {
				tr.With(td)
				td = Td()
				tr.With(td)
			} else {
				td.With(v)
			}
		default:
			td.With(v)
		}
	}

	return Table(Class("layout twocol"), tr)
}

// ----------------------------------------
// navbar

func newNavbar() *navbar {
	return &navbar{
		current:      1,
		groupDivider: make(map[int]bool),
	}
}

type navbar struct {
	max     int // number of slides
	current int // current slide

	groupDivider map[int]bool
}

func (b *navbar) groupEnd(v int) {
	b.groupDivider[v] = true
}

// BuildElement is called at time of rendering
func (b *navbar) next() *Element {
	ul := Ul()

	for i := 0; i < b.max; i++ {
		j := i + 1
		hash := fmt.Sprintf("#%v", j)
		li := Li(A(Href(hash), j))
		if j == b.current {
			li.With(Class("current"))
		}
		if b.groupDivider[j] {
			ul.With(Li(" | "))
		}
		ul.With(li)
	}
	b.current++
	return Nav(ul)
}

// ----------------------------------------
// styling
func DefaultTheme() *CSS {
	css := NewCSS()
	css.Import("https://fonts.googleapis.com/css?family=Inconsolata|Source+Sans+Pro")

	css.Style("html, body",
		"font-family: 'Source Sans Pro', sans-serif",
		"margin: 0 0",
		"padding: 0 0",
	)
	css.Style(".slide",
		"margin: 0 0",
		"padding: 0 0",
		"text-align: center",
		"height: 100vh",
	)
	bg := "#cde9e9"
	css.Style(".slide header",
		"display: block",
		"border: 1px solid "+bg, // needed to make it without margins ?!
		"margin: 0 0",
		"background-color: "+bg,
		"vertical-align: center",
		"height: 10%",
	)

	css.Style(".slide .content",
		// "border: 3px solid red",
		"margin: 0 0",
		"padding: 1% 1%",
		"height: 80%",
		"overflow: auto",
	)
	css.Style(".slide ul, p, pre",
		"text-align: left",
	)

	// ----------------------------------------
	// layout
	css.Style(".layout",
		"width: 100%",
	)
	css.Style(".twocol tr td",
		"width: 50%",
	)

	// ----------------------------------------
	// navbar
	css.Style(".slide nav",
		"display: block",
		"margin: 0 0",
		"padding: 0 0",
		"text-align: center",
		"clear: both",
		"width: 100%",
		"height: 5%",
	)
	css.Style(".slide nav ul",
		"list-style-type: none",
		"margin: 0 0",
		"padding: 0 0",
		"text-align: center",
	)
	css.Style(".slide nav ul li",
		"margin: 0 1px",
		"display: inline",
	)
	css.Style(".slide nav ul li.current a, .slide nav ul li:hover a",
		"background-color: #e2e2e2",
		"border-radius: 5px",
	)
	css.Style("nav a:link, nav a:visited",
		"color: #727272",
		"padding: 3px 5px",
		"text-decoration: none",
		"cursor: pointer",
	)

	// ----------------------------------------
	// Go:ish colors
	css.Style("a:link, a:visited",
		"color: #007d9c",
		"text-decoration: none",
	)
	css.Style("a:hover",
		"text-decoration: underline",
	)
	css.Style(".header",
		"width: 100%",
		"border-bottom: 1px solid #727272",
		"text-align: right",
		"margin-top: -2em",
		"margin-bottom: 1em",
	)
	css.Style("h1, h2, h3",
		"text-align: center",
	)
	css.Style(".twocolumn",
		"width: 100%",
	)
	// ----------------------------------------
	// source code
	css.Style(".srcfile",
		"margin-top: 1.6em",
		"margin-bottom: 1.6em",
		"background-image: url(\"printout-whole.png\")",
		"background-repeat: repeat-y",
		"padding-left: 36px",
		"background-color: #fafafa",
		"tab-size: 4",
		"-moz-tab-size: 4",
		"min-width: 35vw",
	)
	css.Style(".srcfile code",
		"padding: .6em 0 2vh 0",
		"background-image: url(\"printout-whole.png\")",
		"background-repeat: repeat-y",
		"background-position: right top",
		"display: block",
		"text-align: left",
		"font-family: Inconsolata",
		"overflow: wrap",
	)
	css.Style(".srcfile code .line", // each line
		"display: block",
		"width: 98%",
		"clear: both",
		"margin-bottom: -1.5vh",
	)
	css.Style(".srcfile code span:hover", // each line
		"background-color: #b4eeb4",
	)

	css.Style(".srcfile code i",
		"font-style: normal",
		"color: #a2a2a2",
	)
	css.Style(".shell",
		"padding: 1em",
		"border-radius: 10px",
		"min-width: 40vw",
		"overflow: wrap",
	)

	// ----------------------------------------
	// table of contents
	css.Style("li.h3",
		"margin-left: 2em",
	)
	css.Style("table.columns",
		"width: 100%",
	)
	css.Style(".columns td",
		"padding: 0 1em",
		"width: 19%",
	)
	css.Style(".columns .shell",
		"min-width: auto",
	)
	css.Style(".dark",
		"background-color: #2e2e34",
		"color: aliceblue",
	)
	css.Style(".light",
		"background-color: #ffffff",
		"color: #3b2616",
	)
	css.Style("img.center",
		"display: block",
		"margin-left: auto",
		"margin-right: auto",
	)
	css.Style("img.left",
		"float: left",
		"margin-right: 2em",
	)
	css.Style(".group",
		"text-align: left",
		"margin-right: 3em",
	)
	css.Style(".group:first-child",
		"margin-left: 5vw",
	)
	css.Style("td",
		"vertical-align: top",
	)
	css.Style("td:nth-child(2)",
		"padding-left: 2em",
	)
	css.Style("ul.left, ol.left",
		"text-align: left",
	)
	return css
}

// ----------------------------------------
// helpers

func load(src string) *Element {
	v := mustLoad(src)
	v = strings.ReplaceAll(v, "\t", "    ")
	v = highlight(v)
	v = numLines(v, 1)
	return Div(
		Class("srcfile"),
		Pre(Code(v)),
	)
}

func loadFunc(file, src string) *Element {
	v := files.MustLoadFunc(file, src)
	v = strings.ReplaceAll(v, "\t", "    ")
	v = highlight(v)
	return Div(
		Class("srcfile"),
		Pre(Code(v)),
	)
}

func numLines(v string, n int) string {
	lines := strings.Split(v, "\n")
	for i, l := range lines {
		lines[i] = fmt.Sprintf("<span class=line><i>%3v</i> %s</span>", n, l)
		n++
	}
	return strings.Join(lines, "\n")
}

func godoc(pkg, short string) *Element {
	var out []byte
	if short != "" {
		out, _ = exec.Command("go", "doc", short, pkg).Output()
	} else {
		out, _ = exec.Command("go", "doc", pkg).Output()
	}
	v := string(out)
	v = strings.ReplaceAll(v, "\t", "    ")
	v = highlightGoDoc(v)
	return Wrap(
		Pre(v),
		A(Attr("target", "_blank"),
			Href("https://pkg.go.dev/"+pkg),
			"pkg.go.dev/", pkg,
		),
	)
}

func shell(cmd, filename string) *Element {
	v := mustLoad(filename)
	return Pre(Class("shell dark"), cmd, Br(), v)
}

func mustLoad(src string) string {
	data, err := os.ReadFile(src)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

func mustLoadLines(filename string, from, to int) *Element {
	v := files.MustLoadLines(filename, from, to)

	v = strings.ReplaceAll(v, "\t", "    ")
	v = highlight(v)
	v = numLines(v, from)
	return Div(
		Class("srcfile"),
		Pre(Code(v)),
	)
}

// ----------------------------------------

func HighlightColors() *CSS {
	css := NewCSS()
	css.Style(".keyword", "color: darkviolet")
	css.Style(".type", "color: dodgerblue")
	css.Style(".comment, .comment>span", "color: darkgreen")
	return css
}

// highlight go source code
func highlight(v string) string {
	v = keywords.ReplaceAllString(v, `$1<span class="keyword">$2</span>$3`)
	v = types.ReplaceAllString(v, `$1<span class="type">$2</span>$3`)
	v = comments.ReplaceAllString(v, `<span class="comment">$1</span>`)
	return v
}

// highlightGoDoc output
func highlightGoDoc(v string) string {
	v = docKeywords.ReplaceAllString(v, `$1<span class="keyword">$2</span>$3`)
	v = types.ReplaceAllString(v, `$1<span class="type">$2</span>$3`)
	v = comments.ReplaceAllString(v, `<span class="comment">$1</span>`)
	return v
}

var types = regexp.MustCompile(`(\W)(\w+\.\w+)(\)|\n)`)
var keywords = regexp.MustCompile(`(\W?)(package|continue|break|const|select|defer|import|for|func|range|return|go|var|switch|if|case|label|type|struct|interface)(\W)`)
var docKeywords = regexp.MustCompile(`(\W?)(^package|func|type|struct|interface)(\W)`)
var comments = regexp.MustCompile(`(//[^\n]*)`)

// ----------------------------------------

//go:embed enhance.js
var enhancejs string
