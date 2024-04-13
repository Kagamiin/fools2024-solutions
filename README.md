# fools2024-solutions

This repository contains the source code and write-ups for Kagamiin's solutions for TheZZAZZGlitch April Fools Event 2024's Security Testing Program.

## Write-ups for each challenge

NOTE: I was only able to complete challenges 1, 2 and 4. I was unfortunately unable to complete challenge 3.

- [Hacking Challenge I - Hall of Fame Data Recovery (Red/Blue)](/challenge-1/README.md)
- [Hacking Challenge II - The Sus-file (Crystal)](/challenge-2/README.md)
- [Hacking Challenge IV - Pokémon Write-What-Where Version (Emerald)](/challenge-4/README.md)

## Usage

To use this repository, you'll likely need the following tools and assets:

- Go 1.22 (or later)
- mGBA (for challenge 4)
- Any good enough GB/GBC emulator (for challenges 1 and 2 - can be mGBA, too)
- Original Pokémon game ROMs:
  - Pokémon Blue (Game Boy, English)
  - Pokémon Crystal (Game Boy Color, English)
  - Pokémon Emerald (Game Boy Advance, English)

Each folder includes a small downloader script that will download the original challenge asset file from the fools2024 server, using **curl**. This is necessary because the assets cannot be embedded in this repository as their licensing status is dubious at best at the moment.

The challenge asset files may be needed in order to use the solution code to their fullest extent, though they aren't necessarily needed in order to understand the code or even run it (you may create and provide your own wherever needed).

## My thoughts on the challenge

It was so fun. But so tiresome.

It was my first CTF ever, after all. It's to be expected. I wasn't so prepared for it. I do think my approach for the challenges was pretty sane and logical, however, except for challenge 2 where I lost a lot of time at first trying to guess the map location that triggered the password to appear instead of sitting down and actually beginning to write the code to search for it.

Nevertheless, I had a lot of fun chatting with people in the #aprilfools24 channel over at the GCRI Discord, sharing thoughts and tips, supporting each other. Everybody was really cool there and the competition was very healthy. I was very happy to see this.

There was also a second challenge (Fight Simulation Program) consisting of a fangame where you had to earn achievements in order to gain points on the scoreboard, which comprised the majority of the points you could earn on the leaderboard. I didn't have time for it and was honestly not very interested in it, but from what people said about it it's a very interesting game with lots of funny references to Pokémon glitch science. I might play it later, after the challenge is done, just for funsies.

## Licensing

I hereby release the code for my solutions, as well as their documentation (write-ups), under the GNU General Public License v3.
