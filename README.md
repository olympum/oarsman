# WaterRower Interface

* Allow interactive mode, where program starts capturing data without
  sending workout programming to S4. Probably detect SS (stroke start)
  as a way to know when to start requesting data.
* Allow CLI option to log to file, with and without a log file name.
  Use naming standard with log file is not provided.
* Detect end of workout, only possible for distance and duration
  workouts, not interactive.
* Refactor as a library instead of an executable.
