# NativeWiki
NativeWIki provides a lightweight engine that servers markdown files from git
repositories alongside search functionality. The intent is to allow for
documentation not tied to any particular (i.e. our) wiki engine,
versioned in git, and editable from a terminal. To do this
a central server allows for the user to add wikiroots, which are directories
that contain markdown files and are part of a git repository. The server then
indexes those wikiroots and will display:
* A list of Wikiroots
* A global full text search against all roots
* Within a root:
    - A tree of files within that root
    - A search bar against the current wikiroot
    - The current file you are viewing or editing
    - In view mode, some helpful metadata
    - In edit mode, a markdown editor

The server is behind some simple authentication.

## Server

This repository is specifically for the server component

## Building
