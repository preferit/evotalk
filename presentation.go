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

	d.Slide(
		H2("Improve first redesign"),

		P(`We now have multiple stakeholders depending on it and going
        forward they might get affected. Our goal is to be able to
        make changes as freely as possible without affecting the
        stakeholders. Before we release these changes can we improve
        the design?`),

		Table(
			Tr(
				Td(
					shell("$ tree rebel", "ex04.tree"),
					load("../ex/04/main.go"),
				),
				Td(
					load("../ex/04/phrase/phrase.go"),
				),
			),
		),
	)

	d.Slide(
		H2("Fix the repetition when used as a command"),

		Table(
			Tr(
				Td(

					shell("$ tree rebel", "ex05.tree"),
					load("../ex/05/main.go"),
				),
				Td(
					P(`What are the implications of our current design?`),

					Ul(
						Li("Lisa can import it with ",
							Code("import github.com/preferit/rebel/phrase"),
						),
						Li("Max can install it simply with ",
							Code("go install github.com/preferit/rebel@latest"),
						),
					),
					P(`It feels ok, current needs are met. You and Max quickly
        notice that the same phrase appears over and over and add a
        the feature of keeping last shouted phrase in a temporary file
        to minimize the repetition when you run your commands.`),
					load("../ex/05/phrase/phrase.go"),
				),
			),
		),
	)

	d.Slide(
		H2("The first deadend"),

		P(`Some time passes and Lisas swings by your office asking if
        you could implement something that doesn't repeat the same
        phrase every day?  But you just did, how do you now share it
        with Lisa?  At this point you'll realize that the initial
        redesign with the added phrase package, may need to changed
        and your own code updated aswell. This is one of those
        deadends where you have to turn around and find another
        path.`),
		Img(Src("firstdeadend.png")),
		P(`We can explore ways to use the current design, e.g. put
        some caching mechanism into the phrase package, but that seems
        out of place.  Also the caching for the command invocation
        locally on your computer may differ from what is available on
        the intranet site, you don't know at this point.  If we
        backtrack to the time where Max and Lisa came onboard, could
        we have taken a different route that would avoid this
        scenario?`),
	)

	d.Slide(
		H2("What did we do wrong?"),

		P(`Let's go through our reasoning and see if we can find the
        culprit before trying to solve it. Lisa and Max where not very
        specific about their needs just that they wanted to show a
        random phrase in different ways, in the terminal and on a
        webpage. We choose to extract the phrase generating logic for
        Lisa to use and told Max to use the command. Then we added
        more logic in the command, because it was You and Max that saw
        the problem. What Lisa wants now is the logic that is found in
        the command, but her current dependency is elsewhere.`),

		P(`Looking at the code more closely we placed func Shout
        inside the phrase package, why? <em>phrases don't shout, rebels
        do</em>. So if we want to keep func Shout in the package rebel And
        we want to share it with Lisa, how do we solve that?`),

		Table(
			Tr(
				Td(
					shell("$ tree rebel", "ex04.tree"),
					load("../ex/04/main.go"),
				),
				Td(
					load("../ex/04/phrase/phrase.go"),
				),
			),
		),

		B(`Move command, keep domain logic.`),
	)

	d.Slide(
		H2("Move command, keep domain logic"),
		Table(
			Tr(
				Td(
					shell("$ tree rebel", "ex06.tree"),
					load("../ex/06/cmd/rebel/main.go"),
				),
				Td(
					load("../ex/06/rebel.go"),
				),
			),
		),

		P(`This design, is as effortless as the first attempt at the
        decision point, but it makes it easier to evolve the rebel
        logic. It also forces you to design rebel features in such a
        way that they may be used by Lisa as well as Max.`),
	)

	d.Slide(
		H2("Implement a service"),

		P(`The project will grow and at some point Lisa comes by and
        says they are switching languages for the intranet
        implementation but they really want to have access to the
        logic of generating rebelious statements. Could you write a
        service that exposes it?<br> Easy peasy, but where to put it?
        we a few choices`),

		Table(Class("columns"),
			Tr(
				Td(
					"Add the service logic in existing command",
					shell("$ tree rebel", "ex07_1.tree"),
				),
				Td(
					Br(), Br(),
					"Add the service logic in package rebel",
					shell("$ tree rebel", "ex07_2.tree"),
				),
				Td("&nbsp;&nbsp;&nbsp;"),
				Td(
					Br(), Br(), Br(), Br(),
					"Add the service logic in package rebel/service",
					shell("$ tree rebel", "ex07_3.tree"),
				),
				Td("&nbsp;&nbsp;&nbsp;"),
				Td(
					Br(), Br(), Br(), Br(), Br(), Br(), Br(), Br(),
					"Combine option 3 with new command",
					shell("$ tree rebel", "ex07_4.tree"),
				),
			),
		),
	)

	d.Slide(
		H2("Add the service logic in existing command"),
		shell("$ tree rebel", "ex07_1.tree"),

		Table(
			Tr(
				Th("Pros"), Th("Cons"),
			),
			Td(
				Ul(
					Li(`Package rebel remains untouched`),
					Li(`API logic separated from domain logic`),
				),
			),
			Td(
				Ul(
					Li(`Existing command increases in complexity that Max does not need`),
					Li(`API logic mixed with command logic`),
				),
			),
		),
	)

	d.Slide(
		H2("Add the service logic in package rebel"),
		shell("$ tree rebel", "ex07_2.tree"),

		Table(
			Tr(
				Th("Pros"), Th("Cons"),
			),
			Td(
				Ul(
					Li(`I can't see any`),
				),
			),
			Td(
				Ul(
					Li(`Existing command increases in complexity that Max does not need`),
					Li(`API logic mixed with domain logic`),
				),
			),
		),
	)

	return d
}
