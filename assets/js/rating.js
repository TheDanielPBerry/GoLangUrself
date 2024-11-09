window.addEventListener('load', () => {
	const RecordRating = (rating) => {
		jsonData['video']['VideoId']

		const myHeaders = new Headers();
		myHeaders.append("Content-Type", "application/json");

		let videoId = jsonData['video']['VideoId'];
		let url = `/v/${videoId}/rate/${rating}`;

		if(!jsonData?.rating?.RatingId) {
			fetch(url, {
				method: "POST",
			})
			.then((resp) => resp.json())
			.then((json) => {
				if(json?.error != null) {
					window.location.replace('/login');
				} else {
					jsonData['rating'] = json;
				}
			});
		} else {
			fetch(url, {
				method: "PUT",
			})
			.then((resp) => resp.json())
			.then((json) => {
				if(json?.error != null) {
					window.location.replace('/login');
				} else {
					jsonData['rating'] = json;
				}
			});
		}
	};

	document.querySelectorAll('[data-rating]').forEach((button) => button.addEventListener('click', (e) => {
		let target = e.target;
		RecordRating(target.attributes['data-rating'].value);
	}));
});
