package evotalk

import (
	_ "embed"

	. "github.com/gregoryv/web"
)

func Presentation() *Deck {
	d := NewDeck()
	d.Title = "Go; Design for change"

	keywords := func(v ...string) *Element {
		ul := Ul()
		for _, v := range v {
			ul.With(Li(v))
		}
		return ul
	}

	d.Slide(
		Wrap(
			H2("Go; Design for change"),
			Span("Gregory Vinčić, 2023"),
		),
		A(Href("#2"), Img(Src("evotalk.png"))),

		keywords(
			"awareness",
			"change related to sharing",
			"pros and cons",
		),
	)

	d.Slide(
		H2("A story"),

		P(`From the simplest "func main" to a large service. What
        choices to consider when your project grows and how to keep it
        on track for future changes.`, Br(),

			`We begin with a small project and then evolve it.  Along
        the way depending on the choices we make, the <b>design</b>
        will change. <em>I'll refer to directory layout, not as
        structure but as design.</em>`),

		Img(Src("people.png")),

		keywords(
			"design, not structure",
			"focus on package depencencies",
		),
	)

	d.Slide(
		H2("func main()"),

		LayoutTwoCol(

			P(`Let's kick off our project. We'll name it
                    "rebel" and use the repository domain
                    "github.com/preferit/rebel". We choose to create
                    a <em>command</em>, ie. package main. The implications for
                    our ability to evolve `),

			Ol(Class("left"),

				Li(`You Cannot share any logic within it with
                        others (including yourself), Go disallows the
                        import of main packages`),

				Li(`API documentation is hidden when using
                        e.g. go doc, as a result of (1)`),
			),
			"----",
			shell("$ tree rebel", "ex01.tree"),
			load("../ex/01/main.go"),
		),
	)

	d.Slide(
		H2("You code along..."),

		P(`You have no intention to share any logic, the command is
        for you alone. You are happy and code along the nice feature
        of randomizing rebelious statements.`),

		LayoutTwoCol(
			load("../ex/02/main.go"),
			"----",
			shell("$ tree rebel", "ex02.tree"),

			P(`At this point your coworkers Max and Lisa see the work
            and you end up in a discussion;`),
		),
	)

	d.Slide(
		H2("Share with friends"),

		LayoutTwoCol(
			Img(Src("youmaxlisa.png")),
			"----",

			Pre(`
- Max:  Would be nice to see that phrase when I login as the message of the day
- You:  Easy peasy, just go install ... and run it.
- Lisa: Can we include it on the intranet?
- You:  Hmm.. (you pause and start thinking)
- You:  Not really; you could use the binary, but it would be slow with all the
        traffic we have`),

			P(`Here you are faced with a decision on how to share the
            logic of generating a random rebelious statement.`),

			Ol(
				Li(`Redesign the logic as an importable package or`),
				Li(`Write a small service with an API.`),
				Li(`Share the data only and let them figure it out`),
			),

			P(`The first and second option will both require some
               effort. As the consumers are your friends the first
               seems more fitting and much easier to do. The third
               option, though viable, does not help this presentation
               :-).`),
		),
	)

	d.Slide(
		H2("First attempt at redesign"),

		P(`How do you convert the current state, a command, into an
        importable package while also keeping the command. First
        attempt; keep command in root and create a package with logic
        to generate the phrases.`),

		LayoutTwoCol(
			shell("$ tree rebel", "ex03.tree"),
			load("../ex/03/main.go"),
			"----",
			load("../ex/03/phrase/phrase.go"),
		),
	)

	d.Slide(
		H2("Improve first redesign"),

		P(`We now have multiple stakeholders depending on it and going
        forward they might get affected. Our goal is to be able to
        make changes as freely as possible without affecting the
        stakeholders. Before we release these changes can we improve
        the design?`),

		LayoutTwoCol(
			shell("$ tree rebel", "ex04.tree"),
			load("../ex/04/main.go"),
			"----",
			load("../ex/04/phrase/phrase.go"),
			P(`What are the implications of our current design?`),

			Ul(
				Li("Lisa can import it with ",
					Code("import github.com/preferit/rebel/phrase"),
				),
				Li("Max can install it simply with ",
					Code("go install github.com/preferit/rebel@latest"),
				),
			),
		),
	)

	d.Slide(
		H2("Minimize repetition"),

		LayoutTwoCol(

			P(`It feels ok, current needs are met. You and Max quickly
            notice that the same phrase appears over and over and add
            a the feature of keeping last shouted phrase in a
            temporary file to minimize the repetition when you run
            your commands.`),

			load("../ex/05/main.go"),
			"----",
			load("../ex/05/phrase/phrase.go"),
		),
	)

	d.Slide(
		H2("The first deadend"),

		P(`Shortly after, Lisa swings by your office asking if you
        could implement something that doesn't repeat the same phrase
        every day?  But you just did, how do you now share it with
        Lisa?  At this point you'll realize that the initial redesign
        with the added phrase package, <em class="deadend">may need to
        change and your own code updated</em>. This is one of those
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

		LayoutTwoCol(
			shell("$ tree rebel", "ex04.tree"),
			load("../ex/04/main.go"),
			"----",
			load("../ex/04/phrase/phrase.go"),
		),

		B(`Move command, keep domain logic.`),
	)

	d.Slide(
		H2("Move command, keep domain logic"),
		LayoutTwoCol(
			shell("$ tree rebel", "ex06.tree"),
			load("../ex/06/cmd/rebel/main.go"),
			"----",
			load("../ex/06/rebel.go"),
		),

		P(`This design, is as effortless as the first attempt at the
        decision point, but it makes it easier to evolve the rebel
        logic. It also forces you to design rebel features in such a
        way that they may be used by Lisa as well as Max.`),

		P(`With this design, if Lisa came along with the same request;
		you could easily tell her to set MinimizeRepetition to
		true.`),

		Ul(
			Li("Lisa import ",
				Code("import github.com/preferit/rebel"),
			),
			Li("Max install",
				Code("go install github.com/preferit/rebel/cmd/rebel@latest"),
			),
		),
	)
	// ----------------------------------------
	d.GroupEnd()

	serviceOpt := func(show int) {

		opt := func(n int, e *Element) *Element {
			if n == 0 {
				return Wrap()
			}
			if n <= show {
				return e
			}
			return Wrap()
		}

		d.Slide(
			H2("Lisa wants a service"),

			P(`The project will grow and at some point Lisa comes by and
        says they are switching languages for the intranet
        implementation but they really want to have access to the
        logic of generating rebelious statements. Could you write a
        service that exposes it?<br> Easy peasy, but where to put it?
        we a few choices`),

			Table(Class("columns"),
				Tr(
					Td(
						H3(1),
						"Add the service logic in existing command",

						opt(1,
							Wrap(
								shell("$ tree rebel", "ex07_1.tree"),
								H3("Pros"),
								Ul(
									Li(`Package rebel remains untouched`),
									Li(`API logic separated from domain logic`),
								),
								H3("Cons"),
								Ul(
									Li(`Existing command increases in complexity that Max does not need`),
									Li(`API logic mixed with command logic`),
								),
							),
						),
					),
					Td(
						H3(2),
						"Add the service logic in package rebel",
						opt(2,
							Wrap(
								shell("$ tree rebel", "ex07_2.tree"),
								H3("Pros"),
								Ul(
									Li(`I can't see any`),
								),
								H3("Cons"),
								Ul(
									Li(`Existing command increases in complexity that Max does not need`),
									Li(`API logic mixed with domain logic`),
								),
							),
						),
					),
					Td(
						H3(3),
						"Add the service logic in package rebel/service",
						opt(3,
							Wrap(
								shell("$ tree rebel", "ex07_3.tree"),
								H3("Pros"),
								Ul(
									Li(`Package rebel remains untouched`),
									Li(`API logic separated from domain logic`),
								),
								H3("Cons"),
								Ul(
									Li(`Existing command increases in complexity that Max does not need`),
								),
							),
						),
					),
					Td(
						H3(4),
						"New command only",
						opt(4,
							Wrap(
								shell("$ tree rebel", "ex07_4.tree"),
								H3("Pros"),
								Ul(
									Li(`Package rebel remains untouched`),
									Li(`API logic separated from domain logic`),
									Li(`Existing command remains intact`),
								),
								H3("Cons"),
								Ul(
									Li(`Build complexity increases, each stakeholder has their own command`),
								),
							),
						),
					),
					Td(
						H3(5),
						"Combine option 3 and 4",
						opt(5,
							Wrap(
								shell("$ tree rebel", "ex07_5.tree"),
								P("Let's go with this one. Lisa got her service and Max is unaffected by the change."),
							),
						),
					),
				),
			),
		)
	}
	serviceOpt(0)
	serviceOpt(1)
	serviceOpt(2)
	serviceOpt(3)
	serviceOpt(4)
	serviceOpt(5)

	// ----------------------------------------
	d.GroupEnd()

	crawlOpt := func(show int) {

		// opt returns the given e if n is less or equal to show
		opt := func(n int, e *Element) *Element {
			if n == 0 || n > show {
				return Wrap()
			}
			return e
		}

		d.Slide(
			H2("The crawl feature"),

			P(`The service is up and running and you get an idea of adding
		a feature to package rebel for scanning the web for statements
		that the rebel can shout. It doesn't feel like the feature
		fits in the rebel package directly but will be used by
		it. Let's call the feature crawl, but where to put it?`),

			Table(Class("columns"), Tr(
				Td(
					H3(1),
					"Add it directly in package rebel",
					opt(1,
						Wrap(
							shell("$ tree rebel", "ex08_1.tree"),
							H4("Pros"),
							Ul(
								Li(Em("Easy to edit initially")),
							),
							H4("Cons"),
							Ul(
								Li("Tightly coupled with the rebel domain logic"),
							),
						),
					),
				),
				Td(
					H3(2),
					"Add it to an internal/crawl package",
					opt(2,
						Wrap(
							shell("$ tree rebel", "ex08_2.tree"),
							H4("Pros"),
							Ul(
								Li("Separated from domain logic"),
								Li("Cannot be imported by modules outside of the rebel module"),
							),
							H4("Cons"),
							Ul(
								Li("Can be imported by packages within this modules that you might want to move"),
							),
						),
					),
				),
				Td(
					H3(3),
					"Add to a sub package",
					opt(3,
						Wrap(
							shell("$ tree rebel", "ex08_3.tree"),
							H4("Pros"),
							Ul(
								Li("Decoupled from domain logic"),
							),
							H4("Cons"),
							Ul(
								Li("You immediately need to be aware of it being importable"),
							),
						),
					),
				)), // end TR
			),
		)
	}

	crawlOpt(0)
	crawlOpt(1)
	crawlOpt(2)
	crawlOpt(3)

	// ----------------------------------------
	d.GroupEnd()

	d.Slide(
		H2("The team grows"),

		P(`You've come a long way since that early `, Code(`func
		main`), `.  The service is growing as Lisa finds new features
		to add. At some point you need to bring in more people to work
		on the service side. Initially you might start working in the
		same repository but down the road it may make more sense to
		split the service into it's own. The latter is what we're
		interested in. `),

		Table(Class("columns"),
			Tr(
				Td(
					B("Let's go with option 2"),
					shell("$ tree rebel", "ex08_2.tree"),
				),
				Td(

					Img(Src("rebelsrv.png")),
				),
				Td(
					Img(Src("rebelsrv2.png")),
				),
			),
		),
	)

	d.Slide(
		H2("Dependencies as layers"),

		P(`The design for change requires awareness of where we might
		end up in the future and selecting a route that enables it as
		frictionless as possible. Ie. if we'd selected `, A(Href("#17"), "option 4"), ` when
		adding the service, the move at the end would have been even
		easier at the expense of mixing service logic with command
		logic. Depending on the service this may be a good tradeof.`),

		Img(Src("rebelsrvok.png")),
	)

	d.Slide(
		H2("Final thoughts"),
		Ul(Class("summary"),
			Li("Keep domain logic in root"),
			Li("Add exposing layers ABOVE the domain layer"),
			Li("Expose minimal API for your stakeholders"),
			Li("Add internal technical layers BELOW in internal/"),
			Li("Only import internal packages from parent"),
		),
	)

	return d
}
