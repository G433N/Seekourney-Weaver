# Overview

Seekourney Weaver is a search engine for local files and for websites. 
It consists of two main parts; the front-end and the back-end.   
The front-end is where you search for words, function signatures, etc 
and recive a list of files and websites that contain your search expressions. 
It has a browser version and a TUI version.   
The back-end consists of a server, 
which handles requests from all parts of the system, 
a database that stores data about files and websites 
and indexers which retrive words from files and websites.   
Information on these two parts can be found in READMEs 
in their respective folders.

# Required programs

Docker and Go version 1.24.2 need to be installed to run.
Docker must be runnable without using `sudo`.
See a [Docker documentation guide](https://docs.docker.com/engine/install/linux-postinstall/) on this.
