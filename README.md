# rka-control
Simple Go program for cross-platform control over the Roccat Kone AIMO mouse.

## Use
First, download the latest binary from [the releases page](https://github.com/Squ1dd13/rka-control/releases).

*The x86-64 binary provided has been tested on my Linux machine running Ubuntu 20.04. I have not provided builds
for Windows or macOS.*

You will also need to download the [`lights.yml`](https://raw.githubusercontent.com/Squ1dd13/rka-control/main/lights.yml) 
file, as this is how you specify the colours.

To test the program, run this command (modified to fit with your file paths):
```shell
'/path/to/rka-control-linux-x86-64' '/path/to/lights.yml'
```
If you get an error like `hidapi: failed to open device`, you probably need to run
the command as root (e.g. with `sudo`).

If the colours on your mouse change to a mixture of blues and greens,
the program works. (If not, it doesn't...) 

You can pick your own colours by changing the hex strings in `lights.yml`. They are
in typical RGBA hex format, but you may omit the alpha value and just use RGB (if you
do this, the alpha will be set to the maximum by default).

For `left_ribbon` and `right_ribbon`, there are four colours each. These 'ribbons' are
the lights that run down the sides of the mouse all the way to the back, and each one
has four configurable lights inside.

There is also a `brightness` value in the file, which controls the brightness for all
of the lights on the mouse. Its value must be between 0 (off) and 255 (full brightness).

## Disclaimer
This software should be regarded as **experimental** and as such, any damage caused to your hardware as an outcome
of the use of this software is **not my responsibility**. 

**Use at your own risk.**