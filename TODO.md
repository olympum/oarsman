# TODO

* Web interface read-only mode to see workouts and splits.
* Web interface workout mode, to monitor real-time the workout data.
* Allow importing .erg file to define workout, then visualize on the
  web interface during workout mode, ala TrainerRoad.
* Native app mode, via
  [Node-WebKit](https://github.com/jyapayne/Web2Executable)
* [Strava integration](https://github.com/strava/go.strava) for uploads
* Enable a C2 emulation mode using the power watts as the source of
  truth, which we then translate into pace (and speed) using the C2
  formula, and from speed we infer distance (create a difference
  C2Aggregator in order to achieve this).
* Rewrite the `collector.go` using a reactive functional style.
