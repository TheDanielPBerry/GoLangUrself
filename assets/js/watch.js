window.onload = () => {
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
			let url = `/v/${videoId}/watch/${timer}/`;

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

	ProgressSlider.addEventListener('input', (e) => {
		paused = true;
		ProgressTime.innerText = ConvertToReadableRuntime(e.target.value);
	});
	
	ProgressSlider.addEventListener('change', (e) => {
		paused = false;
		timer = e.target.value;
		RecordWatchEvent();
	});
	
	setInterval(() => {
		if(!paused) {
			timer++
			ProgressTime.innerText = ConvertToReadableRuntime(timer);
			ProgressSlider.value = timer;
			if(timer >= 10 && timer % 10 == 0) {
				//Only update every ten seconds
				RecordWatchEvent();
			}
		}
	}, 1000);

	window.addEventListener('beforeunload', function (e) {
		RecordWatchEvent();
	});
};
