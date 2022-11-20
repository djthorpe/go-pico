# go-pico

This repository contains a tinygo implementation of the Raspberry Pi Pico (RP2040) Microcontroller and everything else you need in order to start developing golang on the RP2040  with [tinygo](https://tinygo.org/). The support for some SDK features are missing in tinygo, so this is an alternative which more closely tracks the [Pico SDK](https://github.com/raspberrypi/pico-sdk).

  * For dependencies and installation, please see [INSTALLATION.md](doc/INSTALLATION.md).
  * Once installed, you can compile and run the [HELLOWORLD.md](doc/HELLOWORLD.md) example.

## Versions

This repository currently tracks the following versions:

   * [Pico SDK 1.4.0](https://github.com/raspberrypi/pico-sdk/tree/1.4.0)
   * [tinygo 0.26.0](https://github.com/tinygo-org/tinygo/tree/v0.26.0)

## Documentation

  * General Purpose IO [GPIO](GPIO.md)
  * Pulse Width Modulation [PWM](PWM.md)
  * Analog to Digital Converter [ADC](ADC.md)

## Contributing & Distribution

__This repository is currently in development and subject to change.__

Please do file feature requests and bugs [here](https://github.com/djthorpe/go-pico/issues).

The license is [Apache 2](LICENSE) so feel free to redistribute. Redistributions in either source
code or binary form must reproduce the copyright notice, and please link back to this
repository for more information:

> ### RP2040 SDK for tinygo
> https://github.com/djthorpe/go-pico/
>
> Copyright (c) 2022, David Thorpe, All rights reserved.


