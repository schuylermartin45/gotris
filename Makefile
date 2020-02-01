##
## File:        Makefile
##
## Author:      Schuyler Martin <schuylermartin45@gmail.com>
##
## Description: Builds Gotris.
##

# Directories
BIN = bin/
SRC = src/

# Compiler
CC = go

# Primary build directive
build:
	$(CC) -o $(BIN) $(SRC)*.go

# Clean directive
clean:
	rm -rf $(BIN)*
