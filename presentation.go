package evotalk

import (
	_ "embed"

	. "github.com/gregoryv/web"
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

	d.Slide(
		H2("First attempt at redesign"),

		P(`How do you convert the current state, a command, into an
        importable package while also keeping the command. First
        attempt; keep command in root and create a package with logic
        to generate the phrases.`),

		Table(
			Tr(
				Td(
					shell("$ tree rebel", "ex03.tree"),
					load("../ex/03/main.go"),
				),
				Td(
					load("../ex/03/phrase/phrase.go"),
				),
			),
		),
	)
	return d
}
