# UDP Log Receiver
This program was designed to receive UDP packets from a logging system, store a configurable number of them in memory, and return them via an HTTP. This can be very useful for debugging programs, or simply watching logs come in and display on a screen. Being UDP-based, if this program is not running logging frameworks should otherwise be unaffected. Having a maximum number of logged lines, memory is well controlled (at an approximate maximum of `KEEP_LOGS` x `UDP_BUFFER` bytes, doubled (for string building), plus some static overhead).

# Environment Variables and Command Line Overrides
Command line options will override their environment variable equivalent.
Environment variable | Command line arg | Default | Meaning
---|---|---|---
`HTTP_PORT` | `--http XXX` | 8080 | Port to respond to HTTP queries.
`UDP_PORT` | `--udp XXX` | 8080 | Port for UDP logs to receive.
`UDP_BUFFER` | `--buffer XXX` | 65000 | Size of UDP buffer for incoming packets in kilobytes. Minimum of 1024 (1K).
`KEEP_LOGS` | `--keep XXX` | 5000 | Number of log lines to keep.
`USE_CONSOLE` | `-c` | false/off | Whether to print log lines to the console

# HTTP API
Any path is supported - current recommendation is just `/logs` (root). The `q` query parameter allows you to filter your query to specific comma-delimited strings that at least one of must be present. The `mode` query parameter allows you to specify plain text response (`text`, the default), or `light` or `dark` HTML web page that will reload from the server every second.

# NLog Configuration
Add the following to your config's `targets`, adjusting `layout` to format your string however you see fit and adjusting the IP and port of the `address` attribute to point to the correct network location:
```
<target xsi:type="Network" name="myUDP" maxMessageSize="65000" address="udp://127.0.0.1:10000" 
  layout="${longdate}|${level}|${message} |${all-event-properties} ${exception:format=tostring}" />
```

Then under `rules`, add the following:
```
<logger name="*" minlevel="Trace" writeTo="myUDP" />
```
