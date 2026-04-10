# Upstream Parity Snapshot

- Framework tag: `v13.3.0`
- Skeleton tag: `v13.1.2`
- Generated: `2026-04-05T08:24:37.461Z`

## Framework Components

| Component     | Package                    | Framework Dependencies                                                                                                                                                                   | PHP Files | PHP Lines |
| ------------- | -------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------: | --------: |
| Auth          | `framework/auth`          | framework/collections, framework/contracts, framework/http, framework/macroable, framework/queue, framework/support                                                                 |        49 |      6259 |
| Broadcasting  | `framework/broadcasting`  | framework/bus, framework/collections, framework/container, framework/contracts, framework/queue, framework/support                                                                  |        22 |      2577 |
| Bus           | `framework/bus`           | framework/collections, framework/contracts, framework/pipeline, framework/support                                                                                                     |        18 |      3416 |
| Cache         | `framework/cache`         | framework/collections, framework/contracts, framework/macroable, framework/support                                                                                                    |        64 |      9134 |
| Collections   | `framework/collections`   | framework/conditionable, framework/contracts, framework/macroable                                                                                                                      |        11 |      8419 |
| Concurrency   | `framework/concurrency`   | framework/console, framework/contracts, framework/process, framework/support                                                                                                          |         6 |       379 |
| Conditionable | `framework/conditionable` | -                                                                                                                                                                                         |         2 |       183 |
| Config        | `framework/config`        | framework/collections, framework/contracts                                                                                                                                              |         1 |       306 |
| Console       | `framework/console`       | framework/collections, framework/contracts, framework/macroable, framework/support, framework/view                                                                                   |        88 |      9021 |
| Container     | `framework/container`     | framework/contracts, framework/reflection                                                                                                                                               |        22 |      2805 |
| Contracts     | `framework/contracts`     | -                                                                                                                                                                                         |       150 |      6031 |
| Cookie        | `framework/cookie`        | framework/collections, framework/contracts, framework/macroable, framework/support                                                                                                    |         5 |       619 |
| Database      | `framework/database`      | framework/collections, framework/container, framework/contracts, framework/macroable, framework/support                                                                              |       247 |     59004 |
| Encryption    | `framework/encryption`    | framework/contracts, framework/support                                                                                                                                                  |         3 |       513 |
| Events        | `framework/events`        | framework/bus, framework/collections, framework/container, framework/contracts, framework/macroable, framework/support                                                              |         7 |      1570 |
| Filesystem    | `framework/filesystem`    | framework/collections, framework/contracts, framework/macroable, framework/support                                                                                                    |        10 |      3199 |
| Foundation    | `framework/foundation`    | -                                                                                                                                                                                         |       245 |     31872 |
| Hashing       | `framework/hashing`       | framework/contracts, framework/support                                                                                                                                                  |         6 |       679 |
| Http          | `framework/http`          | framework/collections, framework/macroable, framework/session, framework/support                                                                                                      |        68 |     11049 |
| JsonSchema    | `framework/json-schema`   | framework/contracts                                                                                                                                                                      |        10 |       673 |
| Log           | `framework/log`           | framework/contracts, framework/support                                                                                                                                                  |        11 |      2106 |
| Macroable     | `framework/macroable`     | -                                                                                                                                                                                         |         1 |       135 |
| Mail          | `framework/mail`          | framework/collections, framework/container, framework/contracts, framework/macroable, framework/support                                                                              |        39 |      6352 |
| Notifications | `framework/notifications` | framework/broadcasting, framework/bus, framework/collections, framework/container, framework/contracts, framework/filesystem, framework/mail, framework/queue, framework/support |        25 |      2708 |
| Pagination    | `framework/pagination`    | framework/collections, framework/contracts, framework/support                                                                                                                          |        18 |      3035 |
| Pipeline      | `framework/pipeline`      | framework/contracts, framework/macroable, framework/support                                                                                                                            |         3 |       462 |
| Process       | `framework/process`       | framework/collections, framework/contracts, framework/macroable, framework/support                                                                                                    |        14 |      2489 |
| Queue         | `framework/queue`         | framework/collections, framework/console, framework/container, framework/contracts, framework/database, framework/filesystem, framework/pipeline, framework/support               |       112 |     12577 |
| Redis         | `framework/redis`         | framework/collections, framework/contracts, framework/macroable, framework/support                                                                                                    |        16 |      2691 |
| Reflection    | `framework/reflection`    | framework/collections, framework/contracts                                                                                                                                              |         3 |       439 |
| Routing       | `framework/routing`       | framework/collections, framework/container, framework/contracts, framework/http, framework/macroable, framework/pipeline, framework/session, framework/support                    |        62 |     11162 |
| Session       | `framework/session`       | framework/collections, framework/contracts, framework/filesystem, framework/support                                                                                                   |        16 |      2981 |
| Support       | `framework/support`       | framework/collections, framework/conditionable, framework/contracts, framework/macroable, framework/reflection                                                                       |       106 |     20245 |
| Testing       | `framework/testing`       | framework/collections, framework/contracts, framework/macroable, framework/support                                                                                                    |        30 |      6771 |
| Translation   | `framework/translation`   | framework/collections, framework/contracts, framework/filesystem, framework/macroable, framework/support                                                                             |        11 |      1798 |
| Validation    | `framework/validation`    | framework/collections, framework/container, framework/contracts, framework/macroable, framework/support, framework/translation                                                      |        46 |     11750 |
| View          | `framework/view`          | framework/collections, framework/container, framework/contracts, framework/events, framework/filesystem, framework/macroable, framework/support                                    |        53 |      8852 |

## Local Components

| Component     | Files | go.mod | package.json | doc.go |
| ------------- | ----: | ------ | ------------ | ------ |
| auth          |    33 | yes    | yes          | yes    |
| broadcasting  |     5 | yes    | yes          | yes    |
| bus           |     5 | yes    | yes          | yes    |
| cache         |     5 | yes    | yes          | yes    |
| collections   |     5 | yes    | yes          | yes    |
| concurrency   |     5 | yes    | yes          | yes    |
| conditionable |     5 | yes    | yes          | yes    |
| config        |    13 | yes    | yes          | yes    |
| console       |     6 | yes    | yes          | yes    |
| container     |     5 | yes    | yes          | yes    |
| contracts     |     5 | yes    | yes          | yes    |
| cookie        |     5 | yes    | yes          | yes    |
| database      |     5 | yes    | yes          | yes    |
| encryption    |    10 | yes    | yes          | yes    |
| events        |     5 | yes    | yes          | yes    |
| filesystem    |     5 | yes    | yes          | yes    |
| foundation    |    10 | yes    | yes          | yes    |
| hashing       |    13 | yes    | yes          | yes    |
| http          |     6 | yes    | yes          | yes    |
| json-schema   |     5 | yes    | yes          | yes    |
| log           |     5 | yes    | yes          | yes    |
| macroable     |     5 | yes    | yes          | yes    |
| mail          |     5 | yes    | yes          | yes    |
| notifications |     5 | yes    | yes          | yes    |
| pagination    |     5 | yes    | yes          | yes    |
| pipeline      |     5 | yes    | yes          | yes    |
| process       |     5 | yes    | yes          | yes    |
| queue         |     5 | yes    | yes          | yes    |
| redis         |     5 | yes    | yes          | yes    |
| reflection    |     5 | yes    | yes          | yes    |
| routing       |     6 | yes    | yes          | yes    |
| session       |     5 | yes    | yes          | yes    |
| support       |    12 | yes    | yes          | yes    |
| testing       |     5 | yes    | yes          | yes    |
| translation   |     5 | yes    | yes          | yes    |
| validation    |     5 | yes    | yes          | yes    |
| view          |     6 | yes    | yes          | yes    |
