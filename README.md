# wiz-lights-kde-night-light

This is a go application that automatically changes the temperature of a smart LED ([Wiz](https://www.wizconnected.com/) lights) when night light activates on the system.

## Requirements

- A Smart Wi-Fi led that can be controlled using WiZ
- Desktop environment has to be KDE

## How to use

1. Pair the LED to your Wi-Fi network using the WiZ app
2. Note down the IP address & MAC address of the LED
3. In your router settings, under DHCP settings, assign a static IP for the LED's mac address so that the address does not change.
4. Run the command:

```sh
$ wiz-lights-kde-night-light -bulb-ip <your bulb ip address>
```

```
Usage of wiz-lights-kde-night-light:
  -bulb-ip string
        ip address of the bulb (default "192.168.1.140")
  -bulb-port string
        port of the bulb (default "38899")
```

## Behavior

The application monitors the KDE night light status and automatically adjusts the WiZ bulb temperature accordingly. It sets the bulb temperature to the monitor's temperature.

## How it works ?

This application monitors signals from the KDE night light service, and on receiving a property change signal, it sends an UDP packet to the bulb.

## Future enhancements

- Save & Name the bulbs
- Bulb discovery (through UDP broadcast)
- Change parameters of multiple bulbs at the same time
- Modify other parameters like color, brightness, etc