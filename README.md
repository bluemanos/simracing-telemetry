# Simracing Telemetry

[![Status](https://img.shields.io/badge/status-active-success.svg)]()
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](/LICENSE)

---

Record and display telemetry from:

* Forza Motorsport 2023
* Forza Motorsport 7

Fully configured. Written in Golang.

Plans to support: F1 2023, etc. And make a dashboard/cockpit view for them.

---

### Configuring Forza's UDP settings

1. Launch the game and head to the HUD options menu
2. Set `Data Out` to `ON`
3. Set `Data Out IP Address` to your computer's IP address
4. Set `Data Out IP Port` to `9999`
5. Set `Data Out Packet Format` to `CAR DASH`

### Running the App

#### Docker

1. `git clone https://github.com/bluemanos/simracing-telemetry.git`
2. `cd simracing-telemetry`
3. Set all the environment variables in the `.env` file
4. `docker compose up` which will also build the app

#### Local build

1. `git clone https://github.com/bluemanos/simracing-telemetry.git`
2. `cd simracing-telemetry`
3. Set all the environment variables in the `.env` file
4. `go build`
5. Run `./simracing-telemetry`

#### Use released binary from GitHub

1. Download the latest release from [Releases Page](https://github.com/bluemanos/simracing-telemetry/releases).
2. Get a `forzamotorsport` file from `src/telemetry/fms2023/` folder from the repository and save it in similar directory structure next to the binnary.
3. Set all the environment variables in the `.env` file
4. Run `./simracing-telemetry`

---

### Setup Adapters/Converters

Adapters are setup separately for every game.
The correct adapter setup schema is: `adapter-name:variable1:variable2:etc`.
For multiple adapters configuration use comma, eg: `adapter1:var1:var2,adapter2:var3:var4`.

Currently two adaters are supported:

1. [CSV](#csv-adapter)
2. [MySQL/MariaDB](#mysql-adapter)
3. [UDP forwarder](#udp-forwarder)

#### CSV Adapter
Example: `csv:./data/forzams2023:daily`
* `./data/forzams2023` a path to a directory or file where the CSV files will be saved
* `daily` a record interval. Possible values: `daily` and `none`. Daily retention need a path to directory, `none` retention need a path to file.

#### MySQL Adapter
Example: `mysql:user:password:host:3306:database`
* `user` a MySQL user
* `password` a MySQL password
* `host` a MySQL host
* `3306` a MySQL port
* `database` a MySQL database name

#### UDP forwarder
This adapter can forward the UDP packets to another IPs addresses.

Example: `udp:192.168.5.38:9999&192.168.5.26:9999`
* `ip` a MySQL user
* `port` a MySQL password

More IPs and ports can be added with `&` separator.