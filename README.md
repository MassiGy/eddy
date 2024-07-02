# Eddy - A minimal text editor


## Motivation & Aims

Eddy was a small project that helped me learn more about how games work. I wanted to explore this realm, and since I love text editors 
( I use vim btw ), I thought it would be great to create one.

**But why creating a text editor to learn how games work ?** Well they are very similar, basically you have a setup function that prepares 
everything that the program/interface might need as resources, then you'll have the main loop in which you catch user input and update
the current state (update & rerendering). 

So that is why a text editor is not much diffrent then a game. You draw the editor interface, you enter your main loop, you listen for
user input and you redraw to the screen the updated state.

## Screenshot

- Eddy v0.2.x screenshot, ( a very minimal editor, not designed for mainstream use, only experimental )

![Eddy screenshot](./eddy_screenshot.png "Eddy v0.2.x")


## Setup & Installation ( Without Go tools )

- First, download from Github the latest release. 
    - for Linux and OSX(MAC), download the eddy binary.
    - for Windows, download the eddy.exe binary.
- Add it to your PATH env variable, and run it.

**NOTE:** For those of you that are using windows and really want to change the code and build the program by themselves, it would be better to use MS-Powershell instead of DOS-CMD, since Powershell supports better POSIX commands (the makefile will run better on it).


## Uninstall

- Delete the binary file.
- Remove the binary from your PATH env variable.