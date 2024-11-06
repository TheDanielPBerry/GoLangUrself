window.addEventListener('load', () => {
	const RecordRating = (rating) => {
		jsonData['video']['VideoId']
	};

	document.querySelector('[data-rating]').addEventListener('click', (e) => {
		let target = e.target;
		RecordRating(target.attributes['data-rating'].value);
	});
});
