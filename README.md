# DraWFH

Draw from home, online. This is a online drawing game intended to have fun with friends and colleagues.

![build and test](https://github.com/jeffbean/drawfh/workflows/build%20and%20test/badge.svg?branch=master)

## Contributing

See the CONTRIBUTING file for more details on setting up a local or even remote deployment of the application.

## Inspiration

During this time of quarantine, work from home has become the new norm. With this drawfh.online : drawing from home online was born.

This coding project will serve as a nice distraction.

### learning

Some technical inspiration is using websockets, the Google App Engine, and maintaining a web application for cost.

- Learning [HTML canvas](https://www.w3schools.com/tags/ref_canvas.asp).
- Learning [websockets](https://godoc.org/golang.org/x/net/websocket).
  - If we feel adventurous [grpc-web](https://github.com/grpc/grpc-web) bi-directional streaming.
- Design and maintain cold stored images and historical statistics for minimal cost.
- Collaborate at small scale and designing a CI/CD system to accommodate contributors.

## terminology

These terms are to help code organization, and help design the internals, as well as give common language when working on the system.

- Player : is a participant that may or may not have login information, this is to say we want to support anonymous players, that will get a random identifier in the system. If a Player configures an authenticated profile, we can track game stats for these players over time.
-- participant : an active drawer, guesser.
-- observer/spectator : a player that can see the game, but will not be drawing or guessing. (future)

- Game Template : a game template is able to create games. A game template can have rules to create games. A game template will allows games to be repeatably created. The game template holds settings and one or more dictionaries of words for a game. These templates will be stored and

- Game : a game is the active games of drawing. The game has X number of rounds per the settings. Each round is complete when each player has drawn at least once. It consists of many players and many words to draw. It has a single drawing player, and many guessers. The game tracks the results of each round played and any meta information we desire.

- Round : a round is when all players in a game have drawn once. If a player joins the game during a round, we append them to the end of the list of required round players and gate the round completion on them drawing. If a Player leaves during a round, the round will be considered complete if they have not gone, or if they already participated.

- Guessing/Drawing : each player will guess the active drawer's picture. When the player is up to draw, they choose from the list of words. These words should be a list of 5 with some variety of lengths.

- Chat Hub : the messages interface that serves both the player input, and status messages from the game state. The chat hub will be created and active at the game template level. This allows general chatting before a game is created or running.

- Words : the dictionary of words will by a unique list of words. These will be case insensitive for guessing.

## design

The idea is that like other image guessing games, we have a set of words that one person is drawing, and the others try to guess. This will be designed for private games for creating and playing games quickly. The design should account for the goal of also having saved and permanent game templates for customized games.

We should strive to support input of different languages.

game templates are intended to persist - we will need authenticated users to manage game templates, so a notion of a game template administrator or owner. This is all in the future, but the design should account for this future. Unique Identifiers should have a human slug to allow better short links, or we have a url shortener later using the UUIDs.

The drawer's canvas will have some configurable input settings, like color and size of the pencil. This will be using the HTML Canvas API and expose more and more of this API as we go along. We can also add an Undo button for undoing the last X time/inputs. Maybe a fun way to limit this is to limit the number of Undos per Player per game, or per drawing.

Sharable links for a game template will be able to be created with the embedded password. game templates without passwords can be joined by any and all players with a game template link.

### game modes

To start we will have a single game type, where guessing happens for points based on some criteria. Other game modes, for example battling, could include using points to manipulate others drawing or guessing systems.

## minimal viable

- single ephemeral game template
- memory database
- static dictionary
- anonymous players (no auth)
- no draw settings (you get a simple black pencil)
- english words

## goals

- Have the ability to provide custom word dictionary (github gist, etc...)
- Private game templates have saved settings
  - ideally comes with a permalink. ( eg. drawfh.online/game template/beans-private-game template )
- Able to save and recall all previous drawn images.
- game template's have history of games, and can view stats from past games.
  - Slideshow of prior drawings etc...
- Track stats for authenticated players
- game templates have password protection option.
- GSuite integration

## user flow

Some high level flows to play the game.

### quick game

With no login and the least amount of input, create a game and start playing.

1. Anonymous or authenticated player creates a "quick game".
    - This will use the server default game template.
    - This anonymous game is created with complete functional default settings, and a default (maybe random) dictionary.
    - can change settings while waiting to start.
2. A sharable game link is created.
3. Any player anonymous or authenticated can then join the game.
4. Game start.

### custom game

If you have a curated list of words and some settings you enjoy.

1. Authenticated player can create a game template.
2. The player adds a new named dictionary for use in games created from the game template.
3. Player then creates a new game using these custom settings.
4. A sharable game link is created.
5. Any player anonymous or authenticated can then join the game.
6. Game start.

## Credits

- your name here
