# Keylogger

This is a simple keylogger that is written in Golang. That is intended to work for Ubuntu.


# Techinical Specifications. 

Under Ubuntu all keymaps are stored in `/usr/share/ibus/keymaps` so first we try to read all the keys and it's mappings in each of the files

Then we use the `showkey` command to get a list of keypresses. 

This is then interrupted by the code.

One Ceveat of this is that `showkey` will exit if there is no key event for 10 seconds. But this I think is acceptable and we can always restart the process. If someone is not typing for 10 seconds then It is highly unlikely they will do so in the next fraction of the second. 

But as soon as we see that this command has stopped we start one back again so that we will not loose many keystrokes. 

# Parsing. 

A lot of the parsing is done using Regexes for now _unitl I find a better way to do this_

# Logging. 

For now only a file based logger is implemented. 

The program tries to log the keys to a file in `/tmp/$ENAME_keys.log` where `$ENAME` is a environment variable and has to be set in `/etc/environment`

# Running. 

`sudo ./keylogger &`
`bg`
`disown`


# Roadmap.

* Copy logs to a FTP server. 
* TCP to some server and lot the results there as opposed to keeping then on the hostmachine. 
* Document on how to auto start the keylogger on boot. _(This is trivial, but still....)
* A better Parse that tracks all keys. 
* A Viewer that will format the logs for better viewability.