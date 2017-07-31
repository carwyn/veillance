# Setting up an Intel NUC for Exhibition

## Step 1: Cabling and Startup

Connect a monitor, keyboard and power wource to the NUC then with the NUC powered **off** plugin in **either**:

* an ethernet cable into the network port on the NUC. The other end should be connected
  to a broadband router or other internet connected host network running DHCP IP address assignment.

OR

* a compatible 3G or 4G USB dongle into one of the USB ports on the back of the NUC.

Power on the NUC and wait for the login prompt on the monitor.

## Step 2: Start Wireless Network

Log into the NUC with the username `root` and the password. Next type the following to start up
the wireless network.

`nmcli c up veillance`

Test the wireless network is working properly by connecting a client device and browse an internet web site.

## Step 3: Start Interceptor Software

On the NUC run the following command to start the network data interceptor:

```
interceptor -i eth0 "port 80 or port 53"
```

The network interface to use above (`eth0`) will vary when a 3G or 4G dongle is used.
To find the name of the interface to use use the following command:

```
ip a
```
This will list all the netork interfaces on the NUC.

To stop the interceptor at any point hold down the `CTRL` key and press the `\` (backslash) key.

## Step 4: Connect Graphical Front End

The simplest way to connect the visualisation front end to the network is to connect it
wirelessly in the same way as any other client. Connect to the `Veillance` Wifi network
with the password `veillance`.

Next point the OpenFrameworks software at the IP address `10.42.0.1` and start the graphical
front end.

Note that any web broswing done on the front end host along with and DNS lookups it triggers
will also be intercepted and sent to the visualisation front end.
