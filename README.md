Secret Hitler Server
===

Secret Hitler Server is a simple go application that manages the game state for secret hitler.
It allows clients to connect to a server sent event channel to recieve events, and the current game state.
Players can post up their votes, actions, etc via a single rest endpoint.

Introduction
---

This package depends on github.com/murphysean/secrethitler.
Just run `go get -u` in your working directory to make sure you have the latest code.
Then `go build` to create the binary.
A note on authentication, for now we are allowing you to specify your user as a query parameter to accelerate development.
Just append ?playerID=$PLAYERID to any of the urls below to authenticate as that player for that request.

Building
---

Using docker you can build and run your app

Build binary for mac:

	docker run -it --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp -e GIT_TERMINAL_PROMPT=1 -e GOOS=darwin -e GOARCH=amd64 golang sh -c "git clone git://github.com/murphysean/secrethitler.git /go/src/github.com/murphysean/secrethitler; go build -v"
	./myapp

Build a docker image:

	docker build -t secret-h:1.0 .

Run the docker image:

	docker run -it --rm -p 8080:8080 --name my-running-secret-h secret-h:1.0

Creating a Player
---

	curl http://localhost:8080/api/players/ -H "Content-Type: application/json" -d '{"email":"murphysean84@gmail.com,"password":"abc123"}'
	{"id":"66543097","email":"murphysean84@gmail.com","username":"murphysean","name":"Sean Murphy","thumbnailUrl":"https://secure.gravatar.com/avatar/12c969a7728fe1bc2fb19c8627af81c9"}

Getting a Player
---

	curl http://localhost:8080/api/players/$PLAYERID
	{"id":"66543097","email":"murphysean84@gmail.com","username":"murphysean","name":"Sean Murphy","thumbnailUrl":"https://secure.gravatar.com/avatar/12c969a7728fe1bc2fb19c8627af81c9"}

Creating a Game
---

	curl http://localhost:8080/api/games/ -H "Content-Type: application/json" -X POST
	{"id":"3a37480d-5f65-433c-8f0d-82a3af1f5b59","eventID":1,"state":"","draw":[],"discard":[],"liberal":0,"facist":0,"failedVotes":0,"players":[],"round":{"id":0,"presidentID":"","chancellorID":"","state":"","votes":[],"policies":null,"enactedPolicy":"","executiveAction":""},"nextPresidentID":"","previousPresidentID":"","previousChancellorID":"","specialElectionRoundID":0,"specialElectionPresidentID":"","winningParty":""}

Listing Games
---

	curl http://localhost:8080/api/games/
	[{"id":"3a37480d-5f65-433c-8f0d-82a3af1f5b59","state":"","players":0}]

Getting a Game
---

	curl http://localhost:8080/api/games/$GAMEID
	{"id":"3a37480d-5f65-433c-8f0d-82a3af1f5b59","eventID":1,"state":"","draw":[],"discard":[],"liberal":0,"facist":0,"failedVotes":0,"players":[],"round":{"id":0,"presidentID":"","chancellorID":"","state":"","votes":[],"policies":null,"enactedPolicy":"","executiveAction":""},"nextPresidentID":"","previousPresidentID":"","previousChancellorID":"","specialElectionRoundID":0,"specialElectionPresidentID":"","winningParty":""}


Posting Events
---

Each authenticated player can post events via the rest api

	curl http://localhost:8080/api/games/$GAMEID/events/ -H "Content-Type: application/json" -d '{"type":"player.join","player":{"id":"a"}}'
	{"id":0,"type":"player.join","moment":"2018-04-11T20:51:45.625893491-06:00","player":{"id":"a","lastReaction":"0001-01-01T00:00:00Z"}}

Subscribing to the event stream
---

Anyone can subscribe to the event stream:

	curl http://localhost:8080/api/games/$GAMEID/events?includeState=true

If the query parameter includeState is set to true, the server will also include the current filtered game state for each event

Technical Details
---

The server is broken down into a few components, which can all be independently tested.
The idea was to create a log of events that could be applied to the game state.
The event log could be replayed at any time to discover the game state at any point in time.
Subscribers to the log could all independently keep their own state machine.

### Validate

Each player initiated action is first validated against the current game state.
The validation ensures that all incoming player events are allowed according to the rules.

### Apply

Apply takes the validated events, and blindly applies them to the current game state.
It also ensures the correct order of events, assigning each event an incrmenting event identifier.

### Engine

The engine is just another subscriber to events.
It will take the incoming event, and then produce additional events to advance the game state.

### Filter

Before any event, or the game state is sent to players it is filtered.
This is done to ensure that information is guarded while the game is in progress.

