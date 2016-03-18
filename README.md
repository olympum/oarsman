# Oarsman: A WaterRower S4 Interface

Simple interface to the WaterRower S4 monitor (USB version)

This is a very early stage hobbyist program to interface with the
WaterRower S4 monitor through the USB port. To be able to read from
the S4, the Prolific PL2303 USB to serial adapter driver may need to
be installed (OS specific). If required, there are versions for Mac
and Windows on Prolific's site.

The available commands are:

    version                   Print the version number
    train                     Start a rowing workout activity
    export                    Export workout data from database
    import                    Import workout data from database
    list                      List all workout activities in the database
    remove                    Remove an activity from the database
    help [command]            Help about any command

The program uses a SQLite3 database to store metadata about the
workout activities. The database file and all raw workout logs are
stored under the folder `.oarsman` in the user's home directory. The
database is created automatically if it does not exist the first time
the program is run.

To do a short 200m test workout:

    $ oarsman train --distance=200

... now ... get rowing. Once done, come back to your computer and hit
RETURN (unfortunately that's the best I can do right now). The program
will save a log file with the raw activity event data, insert the
activity into the database, and export a TCX file. This is an example
of a full 110-minute session:

    $ oarsman train --duration=110m
    INFO: 2014/11/10 Using configuration defaults
    INFO: 2014/11/10 Working folder: /Users/brunofr/.oarsman
    INFO: 2014/11/10 Db folder: /Users/brunofr/.oarsman/db
    INFO: 2014/11/10 Workout folder: /Users/brunofr/.oarsman/workouts
    INFO: 2014/11/10 Temp folder: /var/folders/qv/g537wtg1543clytlpl0xn_tm0000gn/T/com.olympum.Oarsman
    INFO: 2014/11/10 Starting single duration workout: 6600 seconds
    INFO: 2014/11/10 Writing to /var/folders/qv/g537wtg1543clytlpl0xn_tm0000gn/T/com.olympum.Oarsman/2014-11-10T09:28:56Z.log
    INFO: 2014/11/10 >>> Press RETURN to end workout ... <<<
    INFO: 2014/11/10 WaterRower S4 02.10

    INFO: 2014/11/10 Workout completed successfully
    INFO: 2014/11/10 Importing activity from /var/folders/qv/g537wtg1543clytlpl0xn_tm0000gn/T/com.olympum.Oarsman/2014-11-10T09:28:56Z.log
    INFO: 2014/11/10 Reading from /var/folders/qv/g537wtg1543clytlpl0xn_tm0000gn/T/com.olympum.Oarsman/2014-11-10T09:28:56Z.log
    INFO: 2014/11/10 Writing to /var/folders/qv/g537wtg1543clytlpl0xn_tm0000gn/T/com.olympum.Oarsman/iWnsTDJtTgJXV_35f8aOlnVPu5zAsdWq71LRkOBxTW4=
    INFO: 2014/11/10 Parsed activity with start time 1415611737000
    INFO: 2014/11/10 Activity 1415611737000 saved to database
    INFO: 2014/11/10 Activity log saved in /Users/brunofr/.oarsman/workouts/2014-11-10T09:28:57Z.log
    INFO: 2014/11/10 Reading from /Users/brunofr/.oarsman/workouts/2014-11-10T09:28:57Z.log
    INFO: 2014/11/10 Writing to /var/folders/qv/g537wtg1543clytlpl0xn_tm0000gn/T/com.olympum.Oarsman/4zd1VuBbntMg_HOsHU7MZte2vnhmvnnEmZmzR2jegDE=.log
    INFO: 2014/11/10 Writing aggregate data to
    /var/folders/qv/g537wtg1543clytlpl0xn_tm0000gn/T/com.olympum.Oarsman/2014-11-10T09:28:57Z.tcx

If you did not save the TCX file, you can always export individual
activities as TCX (Garmin Training Center). To find out the workout
activity id, first list all available workouts using the `list`
command:

    $ oarsman list
    id,start_time,distance,duration,ave_speed,max_speed,ave_cadence,max_cadence,ave_power,max_power,calories,ave_hr,max_hr
    1397805238100,2014-04-18T07:13:58Z,10000,2307,4.334633723450368,0,20.037261698440233,24,0,0,0,135.1585788561523,149
    1397807779100,2014-04-18T07:56:19Z,10000,2312,4.325259515570934,0,20.91915261565068,24,0,0,0,136.65067012537824,143
    1415685752200,2014-11-11T06:02:32Z,15467,3686,4.196147585458491,5.95,19.970519317748337,27,149.9281783009095,221,799,134.4238188654578,155

The `id` for the 110' workout we just did is `1415685752200`, which we
can export with the `export` command:

    $ oarsman export --id=1415685752200
    INFO: 2014/11/11 Writing aggregate data to /var/folders/qv/g537wtg1543clytlpl0xn_tm0000gn/T/com.olympum.Oarsman/2014-11-11T06:02:32Z.tcx

Note that the activity data events (distance, stroke rate, heart rate,
etc.) are captured from the S4 every 25 ms in the raw log, alongside
with pulse and stroke events. The exports, in TCX and CSV, are done at
a 1000ms resolution (1Hz), i.e. using a track point every second.

All workout activity files follow the RFC3339 for naming based on date
and time.

## Vendoring ##

This project uses vendoring and govendor. To install govendor:

```
go get -u github.com/kardianos/govendor
```

To restore after a checkout:

```
govendor sync
```

To add a dependency to latest, to a tag, or to a specific commit:

```
govendor fetch golang.org/x/net/context
govendor fetch golang.org/x/net/context@v1
govendor fetch golang.org/x/net/context@a4bbce9fcae005b22ae5443f6af064d80a6f5a55
```
