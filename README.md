# Tower Defense
Simple tower defense game.

## Rules
* Towers can shoot and have health equal to their ammo.
    * When ammo or health runs out it disappears.
    * Players have a budget and can place towers with a mouse click
    * Players can select different kinds of towers with different effects (melee, ranged, AOE)
* Creeps spawn and move down the lane towards the base.
    * When a creep runs into a tower it lowers tower health by some ammount.
    * Or ranged creeps can fire from a distance
    * When a tower shoots a creep it lowers creep health
    * When creep health runs out it disappers
* Base is at end of lane and has health
* Single player to see how long you can go
    * Tracks score and level 
* Multi-player
    * Against computer to see who lasts longer
    * Against another player over network

## Game Play
* P or Spacebar to pause
* R to reset game
* Mouse left click to place a tower
* '+' or '-' to adjust game speed
* F for full screen
* L to display grid lines (10x10 cells)
* D to display debug info including range indicators and targeting lines
* Q to quit

## TODOS
* ~~Pass thru debug flag to renderers~~
* ~~Collision detection for creep movement, creeps stop when running into a tower~~
    * ~~Restrict creep and tower spawn on existing objects or partially off the board~~
* Make player base more resilient
* Multiple tower types
    * Player selects types
* Tower art
* Creep art
* ~~Multiple creep types~~
* ~~Click to place towers~~
* ~~Spawn creeps and move towards base~~
* ~~Budget for towers, or cooldowns~~
* ~~Towers attack creeps~~
    * ~~Show bullets~~
    * ~~Fix cooldown~~
* ~~Show text for health and cooldown of towers and creeps~~
* ~~Creeps attack towers~~
* ~~Add money on creep death~~
* ~~Scores~~
    * ~~save high scores~~
* Stretch
    * Computer players
    * Networking players
    * Simulation for testing
    * Play simulated network opponent
