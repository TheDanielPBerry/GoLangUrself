window.addEventListener('load', () => {
	var paused = false;
	var timer = jsonData?.watchEvent?.ProgressSeconds ?? 0;

	const ConvertToReadableRuntime = (seconds) => {
		let hours = Math.floor(seconds / 3600);
		let leftOver = seconds % 3600;

		let minutes = Math.floor(leftOver / 60)
		minutes = minutes.toString().padStart(2, '0');
		leftOver = leftOver % 60;
		seconds = leftOver.toString().padStart(2, '0');

		return `${hours}:${minutes}:${seconds}`;
	};

	const RecordWatchEvent = () => {
		if(jsonData?.video?.TotalRuntimeSeconds > timer) {
			const myHeaders = new Headers();
			myHeaders.append("Content-Type", "application/json");

			let videoId = jsonData['video']['VideoId'];
			let url = `/v/${videoId}/watch/${timer}`;

			if(!jsonData?.watchEvent?.WatchEventId) {
				fetch(url, {
					method: "POST",
				}).then((resp) => resp.json())
				.then((json) => jsonData['watchEvent'] = json);
			} else {
				fetch(url, {
					method: "PUT",
				}).then((resp) => resp.json())
				.then((json) => {
					if(json?.error != null) {
						window.location.replace('/login');
					} else {
						jsonData['watchEvent'] = json;
					}
				});
			}
		}
	};

	//Initialize Total Runtime Display
	let TotalRuntime = document.getElementById('TotalRuntime');
	let ProgressTime = document.getElementById('progress');
	let ProgressSlider = document.getElementById('progress-slider');

	TotalRuntime.innerText = ConvertToReadableRuntime(parseInt(TotalRuntime.innerText));
	ProgressTime.innerText = ConvertToReadableRuntime(timer);
	ProgressSlider.value = timer;

	//When user clicks progress select
	ProgressSlider.addEventListener('input', (e) => {
		paused = true;
		ProgressTime.innerText = ConvertToReadableRuntime(e.target.value);
	});

	//When user mouseups from slider
	ProgressSlider.addEventListener('change', (e) => {
		paused = false;
		timer = e.target.value;
		RecordWatchEvent();
	});
	
	setInterval(() => {
		if(!paused) {
			if(timer < jsonData['video'].TotalRuntimeSeconds) {
				timer++
				ProgressTime.innerText = ConvertToReadableRuntime(timer);
				ProgressSlider.value = timer;
				if(timer >= 10 && timer % 10 == 0) {
					//Only update every ten seconds
					RecordWatchEvent();
				}
			} else {
				paused = true;
			}
		}
	}, 1000);

	window.addEventListener('beforeunload', function (e) {
		//Save the watch time when leaving the page
		RecordWatchEvent();
	});


	const togglePause = () => {
		paused = !paused;
		document.getElementById('pause').children[0].innerText = paused ? 'play_arrow' : 'pause';
	};
	
	document.getElementById('video').addEventListener('click', togglePause);
	document.getElementById('pause').addEventListener('click', togglePause);

	var mouseMoveTime = 0;
	document.getElementById('video-player').addEventListener('mousemove', () => {
		let videoControls = document.getElementById('video-controls');
		videoControls.style.opacity = 1;
		let myTime = Date.now();
		mouseMoveTime = myTime;
		setTimeout(() => {
			if(mouseMoveTime == myTime && !paused) {
				videoControls.style.opacity = 0;
			}
		}, 2000);
	});
});
