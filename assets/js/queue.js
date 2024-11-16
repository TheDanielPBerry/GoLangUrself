window.addEventListener('load', () => {

	const AddToQueue = (videoId, target) => {
		let url = `v/${videoId}/queue`;

		fetch(url, {
			method: "POST",
		})
		.then((resp) => {
			if(resp.status === 201) {
				target.innerHTML = '✓ Added to Queue';
				target.setAttribute('active', true);
				return resp.json();
			} else if(resp.status === 204) {
				target.innerHTML = 'Watch Later';
				target.removeAttribute('active');
			} else {
				return resp.json();
			}
		})
		.then((json) => {
			if(json?.error != null) {
				window.location.replace('/login');
			}
		});
	};

	document.body.addEventListener('click', (e) => {
		let target = e.target;
		if(target.classList.contains('queue-video')) {
			e.preventDefault();
			
			AddToQueue(target.attributes['data-video-id'].value, target);
		}
	});
});
