# WaterRower Interface

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

It should be possible to define a workout structure (interval,
steady-state) with targets (stroke, speed, HR) and then monitor and
alert status against targets (below, in-zone, above). So we need the
data to flow in real-time from the capturing on the USB port connected
to the S4 monitor all the way to the UI. On the capturing and
transport of data onto pre-processing, the flow is already there:


    S4 --- (USB) ---> rower.go -----> log_file -----> lumberjack

Lumberjack monitors the log file (or in fact a whole folder if it's
instructed to do so) and forwards messages onto logstash.

    lumberjack --- (msgpack TCP/UDP) ---> logstash

Keep all the existing pre-processing in logstash and dump the events
onto a queue (the data pushed to this queue is transient since it's
only used for real-time feedback on the current workout program). For
now we can use the Redis output since it's simple, and eventually move
to Kafka. Kafka will also allow pulling into Hadoop, Storm, etc.

The Redis logstash output uses either RPUSH or PUBLISH command onto a
Redis list or channel/queue (respectively). A cursory search reveals
that what most people do is Node.js and socket.io to get the data from
a Redis queue (SUBSCRIBE) onto the browser using wss(I don't see a
good reason for WebSockets here and a long-poll request seems enough
... but anyway for the sake of simplicity).

