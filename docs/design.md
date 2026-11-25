# Own Base Game — Design

## Core Rules
- One base can host up to 4 players.
- Town Hall (TH) level is SHARED for the base: starts at TH1 and upgrades to max for everyone.
- Each player has their OWN resources and builders.
- Base owner chooses which player may build/upgrade on the base.
- Attacks use existing heroes/troops; levels are tied to the base (shared).

## Server State (per base/room)
- BaseID, OwnerPlayerID
- THLevel (shared, 1..max)
- Players: map[playerID] -> {Resources, Builders, Assigned bool}
- MaxPlayers = 4
