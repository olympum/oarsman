(function() {
  var tv = 100;

  // instantiate our graph!
  var graph = new Rickshaw.Graph( {
	  element: document.getElementById("chart"),
	  width: 900,
	  height: 500,
	  renderer: 'line',
	  series: new Rickshaw.Series.FixedDuration([{ name: 'one' }], undefined, {
		  timeInterval: tv,
		  maxDataPoints: 200,
		  timeBase: new Date().getTime() / 1000
	  })
  } );

  graph.render();

  // some data every so often
  var i = 0;
  var data = {};  
  var iv = setInterval( function() {
    for (var key in data) {
      if (data.hasOwnProperty(key)) {
        var element = document.getElementById(key);
        if (element) {
          element.innerHTML = data[key];
        }
      }
    }

    graph.series.addData(data);
    graph.render();
  }, tv );

  var source = new EventSource('/events');
  source.onopen = function (event) {
    console.log("eventsource connection open");
  };
  source.onerror = function() {
    if (event.target.readyState === 0) {
      console.log("reconnecting to eventsource");
    } else {
      console.log("eventsource error");
    }
  };
  source.onmessage = function(e) {
    console.log("event: " + e.data);
    var obj = JSON.parse(e.data);
    if (obj.label === "calories") {
      return;
    }
    if (obj.label === 'total_distance_meters') {
      return;
    }
    if (obj.value === 0) {
      return;
    }
    if (typeof(data[obj.label]) === 'undefined') {
      data[obj.label] = obj.value;
    } else {
      data[obj.label] = ((3 * data[obj.label] + obj.value) / 4.0).toFixed(0);
    }

    console.log(data);    
  };
})();
