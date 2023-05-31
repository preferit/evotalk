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

func Presentation() *deck {
	d := newDeck()
	d.Title = "Go; Design for change"
	d.Styles = append(d.Styles,
		theme(),
		highlightColors(),
	)

	d.Slide(
		H2("Go; Design for change"),
		A(Href("#2"), Img(Src("evotalk.png"))),
		Br(), Br(), Br(),
		Span("Gregory Vinčić, 2023"),
		Br(), Br(), Br(),
	)

	d.Slide(
		H2("Story"),
		Img(Src("people.png")),

		P(mustLoad("sl01.txt")),
		/*Pre(Class("shell dark"),
			"$ git clone git@github.com:preferit/cotalk.git\n",
			"$ cd cotalk",
		),*/
	)

	d.Slide(
		H2("func main()"),

		P(`Let's kick off our project. We'll name it "rebel" and use the repository domain
"github.com/preferit/rebel".

We choose to create a command, ie. package main. The implications for
our ability to evolve
`),
		Ol(
			Li(`You Cannot share any logic within it with others (including
yourself), Go disallows the import of main packages`),
			Li(`API documentation is hidden when using e.g. go doc, as a result of
(1)`),
		),

		shell("$ tree rebel", "ex01.tree"),
		load("../ex/01/main.go"),
	)

	d.Slide(
		H2("You code along.."),

		P(`This is a starting point; you have no intention to share
        any logic you only want to use the command yourself. You are
        happy and code along the nice feature of randomizing rebelious
        statements.`),

		Table(
			Tr(
				Td(
					load("../ex/02/main.go"),
				),
				Td(
					shell("$ tree rebel", "ex02.tree"),

					P(`At this point your coworkers Max and Lisa see
                    the work and you end up in a discussion;`),
				),
			),
		),
	)

	d.Slide(
		H2("Share with friends"),

		Table(
			Tr(
				Td(
					Img(Src("youmaxlisa.png")),
				), Td(

					Pre(`
- Max:  Would be nice to see that phrase when I login as the message of the day
- You:  Easy peasy, just go install ... and run it.
- Lisa: Can we include it on the intranet?
- You:  Hmm.. (you pause and start thinking)
- You:  Not really; you could use the binary, but it would be slow with all the
        traffic we have`),

					P(`Here you are faced with a decision on how to share
                the logic of generating a random rebelious
                statement.`),

					Ol(
						Li(`Redesign the logic as an importable package or`),
						Li(`Write a small service with an API.`),
						Li(`Share the data only and let them figure it out`),
					),

					P(`The first and second option will both require some
               effort. As the consumers are your friiends the first
               seems more fitting and much easier to do. The third
               option, though viable, does not help this presentation
               :-).`),
				),
			),
		),
	)

	return d
}

// ----------------------------------------
// helpers
// ----------------------------------------

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

//go:embed docs/enhance.js
var enhancejs string

// ----------------------------------------

// newDeck returns a deck with default styling and navigation on bottom
func newDeck() *deck {
	return &deck{
		Styles: []*CSS{},
	}
}

// Had this idea of a deck of slides; turned out to be less
// useful. Leaving it here for now.
type deck struct {
	Title  string // header title
	Slides []*Element
	Styles []*CSS // first one is the deck default styling
}

// Slide appends a new slide to the deck. elements can be anything
// supported by the web package.
func (d *deck) Slide(elements ...interface{}) {
	d.Slides = append(d.Slides, Wrap(elements...))
}

// Page returns a web page ready for use.
func (d *deck) Page() *Page {
	styles := Style()
	for _, s := range d.Styles {
		styles.With(s)
	}
	body := Body()
	nav := &navbar{current: 1, max: len(d.Slides)}
	for i, content := range d.Slides {
		j := i + 1
		id := fmt.Sprintf("%v", j)
		slide := Div(Class("slide"), Id(id))

		slide.With(
			Header(content.Children[0]),
		)
		div := Div(Class("content"))
		div.With(content.Children[1:]...)
		slide.With(div, nav.next())
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

type navbar struct {
	max     int // number of slides
	current int // current slide
}

// BuildElement is called at time of rendering
func (b *navbar) next() *Element {
	ul := Ul()
	groupDivider := map[int]bool{
		9:  true,
		13: true, // concurrent
		17: true, // channels
	}

	for i := 0; i < b.max; i++ {
		j := i + 1
		hash := fmt.Sprintf("#%v", j)
		li := Li(A(Href(hash), j))
		if j == b.current {
			li.With(Class("current"))
		}
		if groupDivider[j] {
			ul.With(Li(" | "))
		}
		ul.With(li)
	}
	b.current++
	return Nav(ul)
}

// ----------------------------------------

func theme() *CSS {
	css := NewCSS()
	css.Import("https://fonts.googleapis.com/css?family=Inconsolata|Source+Sans+Pro")

	css.Style("html, body",
		"font-family: 'Source Sans Pro', sans-serif",
		"margin: 0 0",
		"padding: 0 0",
	)
	css.Style(".slide",
		//"border: 3px solid magenta",
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

	// navbar
	css.Style(".slide nav",
		// "border: 3px solid green",
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

	// goish colors
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
		// "border: 3px solid brown",
		"width: 100%",
	)
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

	// toc
	css.Style("li.h3",
		"margin-left: 2em",
	)
	css.Style(".shell",
		"padding: 1em",
		"border-radius: 10px",
		"min-width: 40vw",
		"overflow: wrap",
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
var keywords = regexp.MustCompile(`(\W?)(^package|const|select|defer|import|for|func|range|return|go|var|switch|if|case|label|type|struct|interface)(\W)`)
var docKeywords = regexp.MustCompile(`(\W?)(^package|func|type|struct|interface)(\W)`)
var comments = regexp.MustCompile(`(//[^\n]*)`)

func highlightColors() *CSS {
	css := NewCSS()
	css.Style(".keyword", "color: darkviolet")
	css.Style(".type", "color: dodgerblue")
	css.Style(".comment, .comment>span", "color: darkgreen")
	return css
}
