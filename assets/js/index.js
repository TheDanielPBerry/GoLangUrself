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


		let thumbnail = target.closest('.thumbnail');
		if(thumbnail) {
			let open = thumbnail.getAttribute('focus') == 'focus';

			//Close other open previews
			let rows = document.getElementsByClassName('row-container');
			for(let i=0; i<rows.length; i++) {
				if(open || !rows[i].contains(thumbnail)) {
					let preview = rows[i].querySelector('.preview');
					preview.innerHTML = '';
					preview.classList.add('hide');
				}
			}

			//Clear out focus of all thumbnails
			let focused = document.querySelectorAll('.thumbnail[focus=focus]');
			if(focused.length > 0) {
				for(let i=0; i<focused.length; i++) {
					if(open || focused[i] != thumbnail) {
						focused[i].removeAttribute('focus');
					}
				}
			}

			if(!open) {
				thumbnail.setAttribute('focus', 'focus');

				let preview = target.closest('.row-container').querySelector('.preview');
				//Fetch preview html
				let videoId = thumbnail.getAttribute('data-video-id');
				let url = `/v/${videoId}/preview`;
				fetch(url, {
					method: "GET",
				}).then((resp) => resp.text())
				.then((text) => preview.innerHTML=text);
				preview.classList.remove('hide');
			}
		}


	});
});
