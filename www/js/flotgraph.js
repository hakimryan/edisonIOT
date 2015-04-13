$(function() {

	// We use an inline data source in the example, usually data would
	// be fetched from a server

	var data = [],
		totalPoints = 300,
		x = 0;

	function getRandomData() {

		if (data.length > 0)
			data = data.slice(1);

		// Create sinus function

		while (data.length < totalPoints) {
			var y = Math.sin(x);
			data.push(y);
		}
		x = x + 0.1;

		// Zip the generated y values with the x values

		var res = [];
		for (var i = 0; i < data.length; ++i) {
			res.push([i, data[i]])
		}

		return res;
	}

	var updateInterval = 30;

	var plot = $.plot("#placeholder", [ getRandomData() ], {
		series: {
			shadowSize: 0	// Drawing is faster without shadows
		},
		yaxis: {
			min: -1,
			max: 1
		},
		xaxis: {
			show: false
		}
	});

	function update() {

		plot.setData([getRandomData()]);

		// Since the axes don't change, we don't need to call plot.setupGrid()

		plot.draw();
		setTimeout(update, updateInterval);
	}

	update();

});
 