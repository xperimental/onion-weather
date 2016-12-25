# onion-weather

This is a program for getting weather data from a [NetAtmo](https://www.netatmo.com) account and displaying on the OLED display of an [Onion Omega](https://onion.io/) device.

## Installation

### NetAtmo client credentials

This application tries to get data from the NetAtmo API. For that to work you will need to create an application in the [NetAtmo developer console](https://dev.netatmo.com/dev/myaccount), so that you can get a Client ID and secret.

### Build/download the binary

If you want to run this on the Onion Omega you will need a Go installation that can compile to MIPS32 (Go 1.8 or newer). If you have that version installed you can just get the code using `go get` and then use the provided build script to build a binary:

```bash
go get -d github.com/xperimental/onion-weather
./build-mips.sh
```

If you do not have a matching Go installation you can download a pre-built binary from the releases page.

### Install on device

Now you need to copy the binary to your Omega, for example using `scp`. There is also an [init-script](_contrib/onion-weather) and an [empty configuration file](_contrib/onion-weather.json) in the repository. Place the script in `/etc/init.d`.

Put your client credentials and your username and password into the configuration file and place it in `/etc/onion-weather.json` on the device.

Now you should be able to start the service and see your weather data:

```bash
/etc/init.d/onion-weather start
# To stop
/etc/init.d/onion-weather stop
```

If you want to run the service when the device boots run this command:

```bash
/etc/init.d/onion-weather enable
# To disable
/etc/init.d/onion-weather disable
```

## Debug

For debugging you can enable "dummy display" output on the console using the `--dummy` command-line parameter or using the `useDummy` configuration option.

Debug output will look something like this:

```
#######################
#--- N e t A t m o ---#
#Main room @ 18:42:13 #
#  19.3 C - 47 %      #
#  402 ppm            #
#Outside @ 18:42:56   #
#  1.7 C - 87 %       #
#######################
```
