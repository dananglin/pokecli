# pokecli

## Overview

**pokecli** is a simple CLI application that uses the [PokéAPI](https://pokeapi.co/) for exploring the Pokémon world and capturing Pokémon.

### Repository mirrors

- **Code Flow:** https://codeflow.dananglin.me.uk/apollo/pokecli
- **GitHub:** https://github.com/dananglin/pokecli

## Requirements

- **Go:** A minimum version of Go 1.23.1 is required for building the pokecli. Please go [here](https://go.dev/dl/) to download the latest version.

## Build the application

Clone this repository to your local machine.

```
git clone https://github.com/dananglin/pokecli.git
```

Build the application.

- Build with go.
   ```
   go build -o pokecli ./cmd/pokecli
   ```

- Or build with [mage](https://magefile.org/) if you have it installed.
   ```
   mage clean build
   ```
   
## Example Usage

- Run the application and view the help menu
   ```
   $ ./pokecli

   Welcome to the Pokemon world!
   pokecli > help

   Commands:

   catch   Catch a Pokemon and add it to your Pokedex
   exit    Exit the Pokedex
   explore List all the Pokemon in a given area
   help    Display the help message
   inspect Inspect a Pokemon from your Pokedex
   map     Display the next 20 locations in the Pokemon world
   mapb    Display the previous 20 locations in the Pokemon world
   pokedex List the names of all the Pokemon in your Pokedex
   release Release a Pokemon back into the wild
   visit   Visit a location area
   ```

- Use `map` to page through the location areas in the Pokemon world.
   ```
   pokecli > map
   mturnback-cave-pillar-1
   turnback-cave-pillar-2
   turnback-cave-pillar-3
   turnback-cave-before-pillar-1
   turnback-cave-between-pillars-1-and-2
   turnback-cave-between-pillars-2-and-3
   turnback-cave-after-pillar-3
   snowpoint-temple-1f
   snowpoint-temple-b1f
   snowpoint-temple-b2f
   snowpoint-temple-b3f
   snowpoint-temple-b4f
   snowpoint-temple-b5f
   wayward-cave-1f
   wayward-cave-b1f
   ruin-maniac-cave-0-9-different-unown-caught
   ruin-maniac-cave-10-25-different-unown-caught
   maniac-tunnel-26-plus-different-unown-caught
   trophy-garden-area
   iron-island-area
   ```

- Let's use the `visit` command to visit **iron-island-area**.
   ```
   pokecli > visit iron-island-area
   You are now visiting iron-island-area
   ```

- Use the `explore` command to discover all the Pokémon in this location area.
   ```
   pokecli > explore
   Exploring iron-island-area...
   (using data from cache)
   Found Pokemon:
   - tentacool
   - tentacruel
   - magikarp
   - gyarados
   - qwilfish
   - wingull
   - pelipper
   - finneon
   - lumineon
   ```

- Use the `catch` command to throw a Pokéball at a Pokémon. Currently you have a 50% chance to capture each one.
   ```
   pokecli > catch qwilfish
   Throwing a Pokeball at qwilfish...
   qwilfish escaped!

   pokecli > catch wingull
   Throwing a Pokeball at wingull...
   wingull escaped!

   pokecli > catch lumineon
   (using data from cache)
   Throwing a Pokeball at lumineon...
   lumineon was caught!
   You may now inspect it with the inspect command.

   pokecli > catch gyarados
   Throwing a Pokeball at gyarados...
   gyarados was caught!
   You may now inspect it with the inspect command.
   ```

- Use the `pokedex` command to list the names of all the Pokémon that you've caught.
   ```
   pokecli > pokedex
   Your Pokedex:
     - corsola
     - lumineon
     - gyarados
     - gastly
     - bidoof
     - wobbuffet
     - lunatone
     - corphish
   ```

- Use the `inspect` command to inspect one of the Pokémon that you've caught.
   ```
   pokecli > inspect lunatone
   Name: lunatone
   Height: 10
   Weight: 1680
   Stats:
     - hp: 90
     - attack: 55
     - defense: 65
     - special-attack: 95
     - special-defense: 85
     - speed: 70
   Types:
     - rock
     - psychic
   ```

- If you want to release a Pokémon back into the wild use the `release` command.
   ```
   pokecli > release lunatone
   lunatone was released back into the wild.

   pokecli > pokedex
   Your Pokedex:
     - corphish
     - corsola
     - lumineon
     - gyarados
     - gastly
     - bidoof
     - wobbuffet
   ```
