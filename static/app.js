var last_path = null;

var datachart = {
	datasets: [{
		label: 'Température',
		backgroundColor: 'transparent',
		borderColor: '#c33c54',
		yAxisID: 'temperature',
		data: []
	}, {
		label: 'Humiditée',
		backgroundColor: 'transparent',
		borderColor: '#254e70',
		yAxisID: 'humidite',
		data: []
	}] 
};

var ctx = document.getElementById("graphique").getContext('2d');

var graphique = new Chart(ctx, {
	type: 'line',
	data: datachart,

	options: {
		responsive: false,
		height: "250px",
		width: "250px",
		title: {
			display: true,
			text: "Température et Humiditée de la parcelle"
		},
		scales: {
			xAxes: [{
				type: 'time'
			}],
			yAxes: [{
				type: 'linear',
				display: true,
				position: 'left',
				id: 'temperature'
			}, {
				type: 'linear',
				display: true,
				position: 'right',
				id: 'humidite'
			}]
		}
	}
});

function loadParcelle(id) {
	// Definition de l'en-tête HTTP pour s'authentifier
	var config = {
		headers: {'Authorization': "Bearer " + sessionStorage.token}
	};

	/*
	console.log(config);
	axios.get('/api/capteur/1', config).then(function (response) {
		console.log(response.data);
	})
	*/


	axios.get('/api/capteur/' + id + '/mesure', config)
	.then(function (response) {
		//console.log(response.data);

		datachart.datasets[0].data = [];
		datachart.datasets[1].data = [];

		var mesures = response.data.mesures;
		for (var i = 0; i < mesures.length; i++) {
			datachart.datasets[0].data.push({
				x: mesures[i].date,
				y: mesures[i].temperature
			});

			datachart.datasets[1].data.push({
				x: mesures[i].date,
				y: mesures[i].humidite
			});
		}
		graphique.update();
	})
	.catch(function (error) {
		console.log(error);
	});
}

function clickPath(path, id) {
	if (last_path != null) {
		last_path.style.opacity = 0.6;
	}

	path.style.opacity = 1;
	last_path = path;
	loadParcelle(id);
}

function getSessionToken(username, password) {
	axios.get('/login?username='+username+'&password='+password )
	.then(function (response) {
		if (response.data.token !== 'undefined') {
			console.log("OK")
			console.log(response.data)
			sessionStorage.token = response.data.token;
			$('.ui.tiny.modal').modal('hide');
			return;
		}
		sessionStorage.token = null; 
	})
	.catch(function (error) {
		console.log(error.response)
	});
}

$('.ui.tiny.modal')
	.modal({
		closable: false,
		onApprove: function() {
			var username = document.getElementById("login-username").value; 
			var password = document.getElementById("login-password").value;
			getSessionToken(username, password);
			return false;
		}
	})
	.modal('show')
;