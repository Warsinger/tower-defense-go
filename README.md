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

## TODOS
* ~~Click to place towers~~
* ~~Spawn creeps and move towards base~~
* ~~Budget for towers, or cooldowns~~
* Towers attack creeps
* Creeps attack towers
* Collision detection for creep movement
* Make player base more resilient
* Add money on creep death
* Multiple tower types
    * Player selects types
* Multiple creep types
* Tower art
* Creep art
* Stretch
    * Computer players
    * Networking players
    * Simulation for testing
    * Play simulated network opponent
