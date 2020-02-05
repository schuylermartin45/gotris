# Gotris

## An implementation of Tetris in Go following MVC design principles

## Description
This is a simple text (and maybe GUI in the future)-based version of Tetris
written in Go.

## Personal Goals with this Project
* Learn Go
* Implement something that is close to being Tetris with some personal flair
* Use the principles of MVC design to have multiple rendering options for
  the game. To start with, there will be a simple text debugging mode and
  then later there will be a fancier text mode and maybe a GUI.
* Have fun

## Dependencies
Go lacks a lot of old-school terminal control abilities (`clear`, `getch()`,
etc) so this project uses a 3rd party, cross platform library used by a number
of other text-based Go games, [tcell](https://github.com/gdamore/tcell).

### To Install:
#### Automatic
```bash
make depend
```
#### Manual
```bash
go get "github.com/gdamore/tcell"
```

## Build Intstructions
```bash
make
```

## Usage
```bash
./bin/gotris [render mode]
```
Where `[render mode]` is one of these options:
### `debug`
![Early debug mode screenshot](/media/gotris_early_debug_mode.png)
### `text`
[Work in progress]
