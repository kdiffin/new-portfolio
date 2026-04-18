## Two categories?
Stuff we work on related to computer software is largely split into 2 categories:
- Computational / Mathematical
- Architectural / Software Engineering

## Computational

Breakthroughs in the field such as neural nets, raft, paxos, event loops, so on and so forth come from what I like to call **computational thinking** where your raw computational and logical power is used to solve the problem at hand.

A deep understanding of computer science is required to do such a thing. The algorithms are put to test via proofs, peer reviews of papers and empirical tests. Only after these stages have been gone through is the result adopted by the whole world (if it even sees the light of day!).

The feedback loop for developing and getting these algorithms to "production" (used in some area or field) is extremely lengthy and not-so-dopamine-inducing.

## Architectural

With architectural work, the thing we commonly call "software engineering" though, the **feedback loop is tight**. You write some code, have your fancy hot module reloading pick up the change, see the difference instantly, push it via your fancy CI / CD pipeline, have it deployed automatically, upgrade it to prod (if you even have that procedure) and boom. Your work is delivered instantly! If it is garbage your users / stakeholders will complain, if it is good you will see it in the stability / growth of your application.

We want to minimize the "bad feedback" we get, so when working in any big software project, we start developing opinions on *how* the code should be written, not just solving the specific functional problem itself, which is *computational* thinking, but *how* it should be solved so that other problems don't arise in the future! (whether that problem be in the implementation of the code, the structure of the code, the maintanability of it, anything that *could* cause poor performance.)

This leads us to having discourse around these topics, clashing with our stark opinions, reading and writing books on the matter such as [Clean Code](https://www.oreilly.com/library/view/clean-code-a/9780136083238/), [A Philosophy of Software Design](https://www.oreilly.com/library/view/a-philosophy-of/9781492091944/), [The Pragmatic Programmer](https://pragprog.com/titles/pragh/the-pragmatic-programmer/) and [Hypermedia Systems](https://hypermedia.systems/). 

This work is *Architectural work*. Us trying to solve not just for the "what" but for the how. And get the best version of the how, so that further "what"'s are solved easier.

Even the work of that of a humble web developer faces this problem, going on to manage the **complexity** of the codebase with regards to the implementation. Wrestling the never ending pursuit for a "cleaner" and more maintanable codebase. When the developer does this they are still dealing with the consequences of trying to keep the code clean. Is that one liner implementation too clever? Can it be read by another developer? Does the code being too readable hurt performance? (a classic example being .map and .sort's duplication of the input array). Is there a way to do this in such a way that is better. Both for the user and the developer. That is the battle we go on any time we touch the keyboard.

Saying that might make writing code seem too daunting though, but in reality these processes are usually done in the background, gained from years of experience. It is always worth it to look back on your codebase, figure out where bottlenecks in writing and shipping are, and really get to understand deeply why and how each module is working together with. 

While trying to balance the ergonomics and functionality of our codebase here again we find conflicting philosophies.

Discourse on languages (golang vs javascript), tooling (htmx vs react) and standards (flutter (an engine) vs react native / native (native components)). These problems seem to require *good communication skills* and maybe even managerial skills, beyond just computational know-how. This is notoriously hard to test under an interview circumstance. 

As of late it seems the way we develop these skills are being put under threat. With AI coding being used in most of the projects I see, whether be it hackathons, university mates or when prototyping. It's hard to understand what good architecture truly looks like to a junior nowadays. I feel blessed to have been able to code (if not much) before 2023 and the copilot craze, seeing the downfalls of bad architecture physically (as in literally having to physically type the unnecessary code out). Writing code by hand is a must to experience architectural downfalls. 

You really get to see where concepts such as modularization shine when you're forced to rewrite the same logic twice by hand, or forced to implement a function to a similar use case which could have been generalized. 

## So what now?

Nothing! Go and explore, study data strucures, create a large project, go find a job if you're fortunate enough in this job market. Live happily in the wonderful world which is the world of computation, devoid of any sort of vagueness or imprecision. 

A little bit more usefully though, try to figure out which type of work you like best. Maybe you don't care about your code's structure and best principles, maybe you just want to research the best algorithm, the **latest and fanciest AI model** (which seems to be a trending thing to want to create nowadays!) or maybe you just don't care and you want to build for users. Try out, explore, have a large breadth and get to see where you truly find yourself comfortable. Do not quit on first sign of discomfort though!

Having a "wide" understanding of the field will lead to you being a better programmer and a more valuable employee. You will be able to solve cross-departmental problems and really get to test yourself in fields you might not be fully comfortable with. It is definitely worth it to try, to try to improve and be the best that you could be.
