// Create a new HTML5 EventSource
var source = new EventSource('/events/');

// Create a callback for when a new message is received.
source.onmessage = function(e) {

	var y = e.data * 500 / 1024;
    // console.log(y);
    $("#temperature").text(y);

};


var datapoints = [],
	totalPoints = 300,
	updateInterval = 30,
	i = 0;

function getDatapoints() {

	if (datapoints.length > 0) {
		datapoints = datapoints.slice(1);
	}

	while (datapoints.length < totalPoints) {
		var y = $("#temperature").text();
		datapoints.push(y);
	}

	var res = [];
	for (var i = 0; i < datapoints.length; ++i) {
		res.push([i, datapoints[i]])
	}

	return res;
}

var plot = $.plot("#placeholder", [ getDatapoints() ], {
	series: {
		shadowSize: 0	// Drawing is faster without shadows
	},
	yaxis: {
		min: 15,
		max: 50
	},
	xaxis: {
		show: false
	}
});

function update() {
	plot.setData([getDatapoints()]);
	plot.draw();
	setTimeout(update, updateInterval);
}

update();
