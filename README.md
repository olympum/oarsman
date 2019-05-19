# SmarterRower: A WaterRower S4 Interface for Raspberry Pi

A simple and convenient to use Raspberry Pi logger for your WaterRower.


## How it works

1. Install the logger app to your Raspberry with the install script
2. Connect the Raspberry to WaterRower S4 via USB
3. Start to row
4. Hit 3x OK on the S4 when you are done. 

The logger automatically converts your workout data to TCX file and saves it to a Dropbox directory. You can use the tcx file to share it on Strava, Endomondo, etc. 


## Install

### Prerequerements: 
- Installed Raspberry OS
- SSH access to your Raspberry

### Run ./install.sh 
It will copy over the necessary files and installs everything you need on your Raspberry Pi. 
Please note: from now on when your Raspberry is turned on will contionously listen for the WaterRower S4. 


## Usage 

```
// workout listener
srower listen --export --dropbox

// upload to strava
srower strava 67636812.log --token=2876828678e87w6w87786x786x832

// upload to dropbox
srower dropbox 6726872.log --token=27832676786263asd786287378622

```

## Piping commands

srower listen --export --dropbox



You can provide tokens via ENV variables


## Trubleshooting

### My RPi doesn't work at all when I connect to S4.

Unfortunatelly the S4 can't output power so you need an external power source for your Pi. I use a small 0000 mAh powerbank which can keep my Zero alive for a day which means I can use it for weeks for my workouts.

### I have an error during the installation 

The installation script 


### How can I restart logging? 

After you hit 3x OK wait some seconds and give time to your Pi to save and export the workout. 
Disconnect the Pi from the power source and connect it again (reset).  




# Credits

This is a simplified and refactored version of olympum/oarsman. Many thanks for olympum.
