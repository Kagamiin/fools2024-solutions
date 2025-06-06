# fools2024-solutions

This repository contains the source code and write-ups for Kagamiin's solutions for TheZZAZZGlitch April Fools Event 2024's Security Testing Program.

## Write-ups for each challenge

- [Hacking Challenge I - Hall of Fame Data Recovery (Red/Blue)](/challenge-1/README.md)
- [Hacking Challenge II - The Sus-file (Crystal)](/challenge-2/README.md)
- [Hacking Challenge III - gbhttp](/challenge-3/README.md)
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

Each folder includes a small downloader script that will download the original challenge asset file from [zzazzdzz's fools2024 Git repository](https://github.com/zzazzdzz/fools2024), using **curl**. I am not embedding the assets here because I want to keep this repository under a single license.

The challenge asset files are © TheZZAZZGlitch 2024 and are licensed under the MIT license. They may be needed in order to use the solution code to their fullest extent, though they aren't necessarily needed in order to understand the code or even run it.

## My thoughts on the challenge

It was so fun. But so tiresome.

It was my first CTF ever, after all. It's to be expected. I wasn't so prepared for it. I do think my approach for the challenges was pretty sane and logical, however, except for challenge 2 where I lost a lot of time at first trying to guess the map location that triggered the password to appear instead of sitting down and actually beginning to write the code to search for it.

Nevertheless, I had a lot of fun chatting with people in the #aprilfools24 channel over at the GCRI Discord, sharing thoughts and tips, supporting each other. Everybody was really cool there and the competition was very healthy. I was very happy to see this.

There was also a second challenge (Fight Simulation Program) consisting of a fangame where you had to earn achievements in order to gain points on the scoreboard, which comprised the majority of the points you could earn on the leaderboard. It's a very interesting and well-made game with lots of funny references to Pokémon glitch science. I did play it after I was done with the CTF challenges, and I'd definitely recommend everybody go try it out too, even after the leaderboards aren't up anymore.

## Greetings

- To TheZZAZZGlitch, for being so awesome, for hosting these challenges, for all of your glitch research, for all of your content everywhere, for everything!
- To pret, for all of your research and disassemblies, without which my work wouldn't be possible!
- To GCRI, for being such a wholesome glitch science community! Love y'all!
- To the GBDev and mGBA communities!
- And everybody who supported me during these 16 days of pure madness!

## Licensing

(C) 2024 Kagamiin~

I hereby release the code for my solutions, as well as their documentation (write-ups; see exceptions below), under the [GNU General Public License v3-or-later](https://github.com/Kagamiin/ssdpcm/blob/main/COPYING) (also available online, [here](https://www.gnu.org/licenses/gpl-3.0.html)).

The current document, however, is openly licensed via [CC0](https://creativecommons.org/publicdomain/zero/1.0/).

### Exceptions

The following files are part of the repository's documentation, but are derivative works of material not owned by the repository's copyright holder (Kagamiin~). Such files are not covered by the GNU GPL v3-or-later license and the rights to the original material are (C) 1996, 1998, 2000-2001, 2004-2005 GAMEFREAK inc. and reserved by their respective owners.

Such content is used without permission from its copyright holders, but in good faith of it being fair use/fair dealing due to its purely demonstrative purposes:

- challenge1/docs/flag1.png
- challenge1/docs/tile5d.png
- challenge2/docs/flag2.png
- challenge2/docs/mapblocks.png
- challenge2/docs/plaintext_flag.png
- challenge2/docs/start.png
- challenge4/docs/credits.png
