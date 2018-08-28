# ssm2logger

A fast (like, I mean, SUPER FAST) cross-platform, headless Subaru Select Monitor 2 (SSM2) logging tool written in go.

# Why?
I built a small Raspberry Pi system to live in my 2005 Subaru Outback XT to do data logging automatically every time I drive. Most of the software out there didn't offer the solution(s) I was hoping for. They either required a UI, or were built for another embedded system with features/hardware I didn't have or care to have.

I finally found PiMonitor, which I found I could modify easily enough to run headless. Sadly, it too fell short in that is SUPER slow for datalogging (~1 sample per second).

So, I set-out to build this.

# Credits
I drew inspiration, and copied quite a lot of code from (https://github.com/src0x/LibSSM2), the .NET C# library for SSM2. In fact, I started down the path of trying to use it for my solution, but realized pretty quickly that writing cross-platform .NET Core that talks to serial ports could be quite difficult.

# TODO
* Once I get a handle on the actual SSM2 "library", break it out into it's own project, or make it easily consumable from this one
* Add tests
* Add a cobra CLI
* Finish the MVP functionality
  * Consume a RomRaider XML definition file for parameters
  * Log specific PIDs to a log file
* Allow "plugins" for things that aren't SSM. Specifically my ADS_1256 DAC (https://www.waveshare.com/wiki/High-Precision_AD/DA_Board)

# Stuff I'll probably need later
* https://github.com/janne/bcm2835 - To talk to the Raspberry Pi GPIO

# Useful Stuff
* https://subdiesel.wordpress.com/2011/07/13/ssm2-via-serial-at-10400-baud/ - 10400 baud bump!
