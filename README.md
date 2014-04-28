# WaterRower Interface

This is a very early stage hobbyist program to interface with the
Waterrower S4 monitor through the USB port. To be able to read from
the S4, the Prolific PL2303 USB to serial adapter driver must be
installed. There are versions for Mac and Windows on Prolific's site.

The program captures rowing data (distance, stroke rate, heart rate,
etc.) every 25 ms and dumps them onto stdout. The expectation is that
this output is saved to a log file, which can then be cat and piped
onto other programs that can consume either in real-time during the
workout, or once the workout completes. My particular setup consists
of this program, logstash-forwarder, logstash, ElasticSearch and
Kibana.

Running the Waterrower interface from source:

    go run *.go -duration=45m > logs/myworkout.log

This will start a workout in the S4 for 45 minutes. Other command line
options can be shown with `go run *.go -help`.

Pipe the workout onto logstash-forwarder (aka Lumberjack). This can be
done after the workout, which will dump the file onto elasticsearch,
or during the workout, which keep indexing one event at the time:

    cd config
    cat logs/myworkout.log | logstash-forwarder -config lumberjack.conf

Run logstash:

    cd config
    ~/Projects/be/logstash-1.4.0/bin/logstash -f logstash.conf

Run elasticsearch:

    ~/Projects/be/elasticsearch-1.1.0/bin/elasticsearch

Run kibana:

    ~/Projects/be/logstash-1.4.0/bin/logstash-web

And then visit the kibana interface:

    open http://localhost:9292/index.html

## TODO

* Allow interactive mode, where program starts capturing data without
  sending workout programming to S4. Probably detect SS (stroke start)
  as a way to know when to start requesting data.
* Allow CLI option to log to file, with and without a log file name.
  Use naming standard with log file is not provided.
* Detect end of workout, only possible for distance and duration
  workouts, not interactive. For distance it's straight forward, but
  for time it will likely require fetching the S4 display time, as the
  user can pause the workout.
* Refactor as a library instead of as an executable.
* Add interval workout programming.

## Real-Time Analysis

The current flow is to push the data to logstash, and all UI is driven
from ElasticSearch and Kibana. Although that's a fine tool for
post-workout analysis, it's not enough for real-time in-workout
feedback.

It should be possible to use a UI to define a workout structure
(interval, steady-state) with different targets (stroke, speed, HR)
and then monitor and alert status against targets (below, in-zone,
above). We need the event data to flow in real-time from the capturing
on the USB port all the way to the UI. On the capturing and transport
of data onto pre-processing, the flow is already there:

    S4 --- (USB) ---> rower.go -----> log_file -----> lumberjack

Lumberjack monitors the log file (or in fact a whole folder if it's
instructed to do so) and forwards messages onto logstash.

    lumberjack --- (msgpack TCP/UDP) ---> logstash

Eventually, we want to keep all the existing pre-processing in
logstash and send the events onto a queue (the data pushed to this
queue is transient since it's only used for real-time feedback on the
current workout program). For now we can use the Redis output since
it's simple, and eventually move to Kafka. Kafka will also allow
pulling into Hadoop, Storm, etc.

The Redis logstash output uses either RPUSH or PUBLISH command onto a
Redis list or channel/queue (respectively), and using PUBLISH allows
us to connect a simple client on SUBSCRIBE. We then push to the
browser using server-sent events (no need for WebSockets here).
Server-sent events are not supported in IE, so we need a polyfill.

Since we are going to use a pipe, for development we don't require
logstash, lumberjack, etc. to run the server, processing, etc. We can
simply implement this as:

1. Start basic HTTP server.
1. Browser client connects to HTTP server.
1. User enters workout and client sends commands to HTTP server.
1. HTTP server starts S4 interface and sends workout.
1. Captured events are logged.
1. HTTP server monitors logs and sends new lines to the HTTP client
   via SSE (we can keep loggin probably).

