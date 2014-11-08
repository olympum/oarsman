# Oarsman: A WaterRower S4 Interface

Simple interface to the WaterRower S4 monitor (USB version)

This is a very early stage hobbyist program to interface with the
WaterRower S4 monitor through the USB port. To be able to read from
the S4, the Prolific PL2303 USB to serial adapter driver may need to
be installed (OS specific). If required, there are versions for Mac
and Windows on Prolific's site.

The available commands are:

    version                   Print the version number
    workout                   Start a rowing workout
    export                    Export workout data from database
    import                    Import workout data from database
    help [command]            Help about any command
 
The program uses a SQLite3 database to store metadata about the
workout activities. The database file and all raw workout logs are
stored under the folder `.oarsman` in the user's home directory. The
database is created automatically if it does not exist the first time
the program is run.

To do a short 200m test workout:

    $ oarsman workout --distance=200

... now ... get rowing. Once done, come back to your computer and hit
RETURN (unfortunately that's the best I can do right now). The program
will save a log file with the raw activity event data, insert the
activity into the database, and export a TCX file. This is an example
of a full session:

    $ oarsman workout --duration=90m
    INFO: 2014/10/31 Using configuration defaults
    INFO: 2014/10/31 Working folder: /Users/brunofr/.oarsman
    INFO: 2014/10/31 Db folder: /Users/brunofr/.oarsman/db
    INFO: 2014/10/31 Workout folder: /Users/brunofr/.oarsman/workouts
    INFO: 2014/10/31 Temp folder: /var/folders/qv/g537wtg1543clytlpl0xn_tm0000gn/T/com.olympum.Oarsman
    INFO: 2014/10/31 Starting single duration workout: 5400 seconds
    INFO: 2014/10/31 Writing to /var/folders/qv/g537wtg1543clytlpl0xn_tm0000gn/T/com.olympum.Oarsman/2014-10-31T08:30:18Z.log
    INFO: 2014/10/31 >>> Press RETURN to end workout ... <<<
    INFO: 2014/10/31 WaterRower S4 02.10

    INFO: 2014/10/31 Workout completed successfully
    INFO: 2014/10/31 Importing activity from /var/folders/qv/g537wtg1543clytlpl0xn_tm0000gn/T/com.olympum.Oarsman/2014-10-31T08:30:18Z.log
    INFO: 2014/10/31 Reading from /var/folders/qv/g537wtg1543clytlpl0xn_tm0000gn/T/com.olympum.Oarsman/2014-10-31T08:30:18Z.log
    INFO: 2014/10/31 Writing to /var/folders/qv/g537wtg1543clytlpl0xn_tm0000gn/T/com.olympum.Oarsman/_m9Qlix_R8uN3zDNpkhW2RIs9OBEg548FwwjZA9vJ54=
    INFO: 2014/10/31 Parsed activity with start time 1414744219400
    INFO: 2014/10/31 Activity 1414744219400 saved to database
    INFO: 2014/10/31 Activity log saved in /Users/brunofr/.oarsman/workouts/2014-10-31T08:30:19Z.log
    INFO: 2014/10/31 Reading from /Users/brunofr/.oarsman/workouts/2014-10-31T08:30:19Z.log
    INFO: 2014/10/31 Writing to /var/folders/qv/g537wtg1543clytlpl0xn_tm0000gn/T/com.olympum.Oarsman/HU2iF3wiDQzJvC2XkKXr2FCMa6je_unsNOR0Zhpnfmk=.log
    INFO: 2014/10/31 Writing aggregate data to
    /var/folders/qv/g537wtg1543clytlpl0xn_tm0000gn/T/com.olympum.Oarsman/2014-10-31T08:30:19Z.tcx

If you did not save the TCX file, you can always export individual
activities as TCX (Garmin Training Center). To find out the workout
activity id, first list all available workouts using the `list`
command:

    $ oarsman list
    id,start_time,distance,ave_speed,max_speed
    1404553035100,2014-07-05T09:37:15Z,16408,4.214744,5.950000
    1414596607600,2014-10-29T15:30:07Z,200,3.174603,5.600000

We see the id for the 200m workout we just did is `1414596607600`,
so now we can export it with the `export` command:

    $ oarsman export --id=1414596607600
    2014/10/29 17:29:40 Writing aggregate data to /var/folders/qv/g537wtg1543clytlpl0xn_tm0000gn/T/com.olympum.Oarsman/1414596607600.tcx

Note that whilst the activity data events (distance, stroke rate,
heart rate, etc.) are captured from the S4 every 25 ms, the TCX export
uses a 100ms resolution. The standard usage in TCX files is 1s in most
activity software, although the schema allows any RFC3339 datetime,
which we are using to provide higher frequency.
