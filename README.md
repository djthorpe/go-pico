# go-pico

Experiments with Raspberry Pi Pico Microcontroller. Use `make` command
to compile applications in the `cmd` folder, which is then placed in the `build`
folder. You can load any application (when Pico is in BOOLSEL state) using:

```bash
  sudo /opt/pico/bin/picotool load build/blink.uf2 
```

Minicom can be used for output:

```bash
  minicom --device /dev/serial0
```
