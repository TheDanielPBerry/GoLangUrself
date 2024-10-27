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

				let rowContainer = target.closest('.row-container');
				let preview = rowContainer.querySelector('.preview');
				//Fetch preview html
				let videoId = thumbnail.getAttribute('data-video-id');
				let url = `/v/${videoId}/preview`;
				fetch(url, {
					method: "GET",
				}).then((resp) => resp.text())
				.then((text) => {
					preview.innerHTML=text;
					rowContainer.previousElementSibling.scrollIntoView({behavior: 'smooth'});
				});
				preview.classList.remove('hide');
			}
		}


	});
	window.addEventListener('keydown', (e) => {
		let focused = document.querySelectorAll('.thumbnail[focus=focus]');
		if(focused.length > 0) {
			focused = focused[0];
			let row = focused.closest('.row');
			let nextThumbnail = null;
			let offset = 0;
			if(e.key === 'ArrowRight') {
				nextThumbnail = focused.nextElementSibling;
				if(nextThumbnail !== null) {
					offset = nextThumbnail.getBoundingClientRect().width - 10;
				}
			} else if(e.key === 'ArrowLeft') {
				nextThumbnail = focused.previousElementSibling;
				if(nextThumbnail !== null) {
					offset = -nextThumbnail.getBoundingClientRect().width + 4;
				}
			}
			if(nextThumbnail !== null) {
				e.preventDefault();
				nextThumbnail.click();
				setTimeout(() => row.scrollBy({ left: offset, top: 0, behavior: 'smooth'}), 0);
				return;
			}

			if(e.key === 'ArrowDown') {
				nextThumbnail = focused.closest('.row-container').nextElementSibling?.nextElementSibling?.querySelector('.thumbnail');
			} else if(e.key === 'ArrowUp') {
				nextThumbnail = focused.closest('.row-container').previousElementSibling?.previousElementSibling?.querySelector('.thumbnail');
			}
			if(nextThumbnail !== undefined && nextThumbnail !== null) {
				e.preventDefault();
				nextThumbnail.click();
				return;
			}

			if(e.key === 'Escape') {
				focused.click();
			}
		}
	});
});
