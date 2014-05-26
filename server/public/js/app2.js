    function fillArrayWithNumbers(n) {
        var arr = Array.apply(null, Array(n));
        return arr.map(function (x, i) { return 0 });
    }

var n = 40, data = fillArrayWithNumbers(n), eventData = {};

var margin = {top: 20, right: 20, bottom: 20, left: 40},
    width = 960 - margin.left - margin.right,
    height = 500 - margin.top - margin.bottom;
 
var x = d3.scale.linear()
    .domain([1, n - 2])
    .range([0, width]);
 
var y = d3.scale.linear()
    .domain([0, 500])
    .range([height, 0]);
 
var line = d3.svg.line()
    .interpolate("basis")
    .x(function(d, i) { return x(i); })
    .y(function(d, i) { return y(d); });
 
var svg = d3.select("body").append("svg")
    .attr("width", width + margin.left + margin.right)
    .attr("height", height + margin.top + margin.bottom)
  .append("g")
    .attr("transform", "translate(" + margin.left + "," + margin.top + ")");
 
svg.append("defs").append("clipPath")
    .attr("id", "clip")
  .append("rect")
    .attr("width", width)
    .attr("height", height);
 
svg.append("g")
    .attr("class", "x axis")
    .attr("transform", "translate(0," + y(0) + ")")
    .call(d3.svg.axis().scale(x).orient("bottom"));
 
svg.append("g")
    .attr("class", "y axis")
    .call(d3.svg.axis().scale(y).orient("left"));
 
var path = svg.append("g")
    .attr("clip-path", "url(#clip)")
  .append("path")
    .datum(data)
    .attr("class", "line")
    .attr("d", line);

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
  if (typeof(eventData[obj.label]) === 'undefined') {
    eventData[obj.label] = obj.value;
    tick();
  } else {
    eventData[obj.label] = obj.value; 
    //((39 * eventData[obj.label] + obj.value) / 40).toFixed(0);
  }
};

function tick() {
  // push a new data point onto the back
  if (typeof(eventData['speed_cm_s']) === 'undefined') {
    return;
  }
  data.push(eventData['speed_cm_s']);
 
  // redraw the line, and slide it to the left
  path
      .attr("d", line)
      .attr("transform", null)
    .transition()
      .duration(1000)
      .ease("linear")
      .attr("transform", "translate(" + x(0) + ",0)")
      .each("end", tick);
 
  // pop the old data point off the front
  data.shift();
}
