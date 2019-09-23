# Three Doors Trial
Golang simulator of the famous "Three Doors Trial" (three boxes trial as in the code)

## Scope
This is a simple simulator of the Three Doors Trial aka *The Monty Hall Problem*, which is fully explained elsewhere,
e.g. in this [Scientific American article](https://www.scientificamerican.com/article/the-3-door-monty-hall-problem/) and the references _ibid_.

## The approach
The approach is very simple:
 - generate random number from 0 to 2 (using "crypto/rand" and "math/rand" packages)
 - use this number as the person's under trial initial guess (we, programmers, number the things starting from zero, right?)
 - the treasure is always placed into the same box (random choice of the treasure placement changes nothing since the placement does not depend on the person's choice)
 - if the guess is a "lucky" one (i.e. guessed box number is equal to the number of the box with the treasure) then increment the "win if stay" counter
 - otherwise, increment the "win if change" counter
 
## Results
First, the program runs 10M tests to demonstrate the bit randomness - it calculates ones and zeroes statistics. Typical values are 49.97 to 50.04 percents.

Next, the program runs 10M trials and calculates and prints statistics to verify the randomness and demonstrate the trial results:
  - number of the random numbers generated so far
  - number and percent of the lucky guesses 
  - distribution of the initial guesses among the boxes  
  - number and percent of "wins if stay" (about 1/3)
  - number and percent of "wins if change" (about 2/3)
  
That's it!
