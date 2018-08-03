# Underbot
## Details
The bot that can play independently play Undertale (and some fangames) in it's entirety (that's the goal anyways)

## Build Instructions
### Linux (X11)
1. If you don't already have ```dep``` installed, use this command to do so: ```curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh```
2. Build OpenCV from the [GoCV project](https://github.com/hybridgroup/gocv) (OpenCV 3+ might work, but please try the GoCV version before reporting issues here)
3. Get this project with ```git clone https://gitlab.com/256/Underbot.git```
4.  Run ```dep ensure``` to get the needed packages in the root directory
5.  Finally, use ```go build``` and run the executable named ```Underbot```
### Linux (Wayland)
Until Wayland provides a method to interact with other windows, as it is designed to limit interaction between windows, this is unlikely to ever be in the future of this project.
### Other platforms
This is also not supported yet as the project relies on the X11 window system. However, porting is possible in the case of Windows and is being tracked in [Issue #7](https://gitlab.com/256/Underbot/issues/7). If you think you are up to the task, go ahead and give it a shot at porting it over.

## TODO
View tasks that need to be done at [the issue boards](https://gitlab.com/256/Underbot/boards?=)