# xrandr-gotoggle

`xrandr-gotoggle` is a small tool for saving and loading xrandr configurations.

The tool saves xrandr configurations in a configuration file and references them
by calculating a sha256 sum over the ID's of all connected monitors, by using 
the following format:
```
<calculated sha256 from connected monitors>: ["<xrandr commandline args>"]
5a11328e98459635a78587a4f8b39dbe5f6618932f5247634fce60fe60a19921: ["--screen", "0", "--output", "eDP-1", "--mode", "1920x1080", "--pos", "0x0", "--primary"
```

## Usage

The tool provides the following commands:

| Command                                       | Function                                                                                                                |
|-----------------------------------------------|-------------------------------------------------------------------------------------------------------------------------|
| `xrandr-gotoggle apply`                       | Tries to apply the saved configuration if available.                                                                    |
| `xrandr-gotoggle config view`                 | Prints the current configuration from the configuration file.                                                           |
| `xrandr-gotoggle config print-current-config` | Prints the current configuration, which could be added to the configuration file, and the corresponding xrandr command. |
| `xrandr-gotoggle config set-current-config`   | Adds or overwrites the current configuration to the configuration file.                                                 |
| `xrandr-gotoggle print-checksum`              | Calculates and prints the screen checksum.                                                                              |

## Build

To build the binary you need [go](https://golang.org/) installed on your system and run the following command.

```
go build .
```

## Integration

To automatically apply a configuration on boot the tool should get run once after
starting the window manager.

E.g. for i3 the following line could be added to the configuration at `~/.config/i3/config`:
```
exec --no-startup-id /usr/local/bin/xrandr-gotoggle apply
```

To automatically apply a configuration on changes to the currently connected displays,
an udev-rule could get added (make sure to set the appropriate `DISPLAY` und `user`):
```
SUBSYSTEM=="drm", ACTION=="change", ENV{DISPLAY}=":1", RUN+="/bin/su - user -c '/usr/local/bin/xrandr-gotoggle apply'"
```
