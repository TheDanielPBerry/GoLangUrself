window.addEventListener('load', () => {
	var videoFocus = null;
	document.body.addEventListener('click', (e) => {
		let target = e.target;

		//Close all the previews
		let preview = target.closest('.preview');
		if(preview) {
			//If clicked in a preview, do nothing
			return;
		}

		//Clear out focus of all thumbnails
		let focused = document.getElementsByClassName('thumbnail-focus');
		if(focused.length > 0) {
			for(let i=0; i<focused.length; i++) {
				focused[i].classList.remove('thumbnail-focus');
			}
		}

		let thumbnail = target.closest('.thumbnail');
		if(thumbnail) {
			thumbnail.classList.add('thumbnail-focus');
		}

		let rows = document.getElementsByClassName('row-container');
		for(let i=0; i<rows.length; i++) {
			if(!rows[i].contains(thumbnail)) {
				rows[i].querySelector('.preview').innerHTML = '';
			}
		}
	});
});
