##
## File:        Makefile
##
## Author:      Schuyler Martin <schuylermartin45@gmail.com>
##
## Description: Builds Gotris.
##

# Directories
BIN = ./bin/
SRC = ./src/

# Go Compiler
GC = go
GFLAGS = build


# Primary build directive
build:
	$(GC) $(GFLAGS) -o $(BIN)gotris $(SRC)gotris

# Install dependencies
depend:
	$(GC) get github.com/gdamore/tcell

# Clean directive
clean:
	rm -rf $(BIN)*
