# UDP Log Receiver
This program was designed to receive UDP packets from a logging system, store a configurable number of them in memory, and return them via an HTTP. This can be very useful for debugging programs, or simply watching logs come in and display on a screen. Being UDP-based, if this program is not running logging frameworks should otherwise be unaffected. Having a maximum number of logged lines, memory is well controlled (at an approximate maximum of `KEEP_LOGS` x `UDP_BUFFER` bytes, doubled (for string building), plus some static overhead).

## Environment Variables
- `HTTP_PORT`: Port to respond to HTTP queries. Defaults to 8080.
- `UDP_PORT`: Port for UDP logs to receive. Defaults to 10000.
- `UDP_BUFFER`: Size of UDP buffer for incoming packets. Minimum of 1024 (1K), defaults to 65000 (almost 65K).
- `KEEP_LOGS`: Number of log lines to keep. Minimum of 1, defaults to 5,000 logs.

## HTTP API
Any path is supported - current recommendation is just `/` (root). The `q` query parameter allows you to filter your query to specific comma-delimited strings that at least one of must be present.

## NLog Configuration
Add the following to your config's `targets`, adjusting layout to format your string however you see fit and adjusting the IP and port of the `address` attribute to point to the correct network location:
```
<target xsi:type="Network" name="myUDP" maxMessageSize="65000" address="udp://127.0.0.1:10000" 
  layout="${longdate}|${level}|${message} |${all-event-properties} ${exception:format=tostring}" />
```

Then under `rules`, add the following:
```
<logger name="*" minlevel="Trace" writeTo="myUDP" />
```
