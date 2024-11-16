window.addEventListener('load', () => {
	let resultList = document.getElementById('search-results');
	let search = document.getElementById('search');

	const ToggleResults = () => {
		resultList.classList.toggle('hide');
	};

	const NavigateResult = (li) => {
		let url = '/v/';
		let videoId = li.attributes['data-video-id'].value;
		if(videoId >= 0) {
			url += videoId;
			window.location.href = url;
		}
	};

	const ToggleLoader = () => {
		resultList.innerHTML = '<div class="loader"></div>';
		resultList.classList.remove('hide');
	};

	const UpdateSearchResults = (results) => {
		resultList.style.top = (search.offsetTop + 20) + 'px';
		resultList.style.left = (search.offsetLeft + 1) + 'px';

		resultList.classList.remove('hide');
		resultList.innerHTML = '';
		results.forEach((row, count) => {
			let li = document.createElement('li');
			li.innerText = row.Title.String;
			if(row.Year !== undefined) {
				li.innerText += ' (' + row.Year + ')';
			}

			let videoId = document.createAttribute('data-video-id');
			let focus = document.createAttribute('data-focus');

			videoId.value = row.VideoId;
			if(count == 0) {
				focus.value = true;
				li.ariaSelected = true;
			} else {
				focus.value = false;
				li.ariaSelected = false;
			}

			li.attributes.setNamedItem(videoId);
			li.attributes.setNamedItem(focus);

			resultList.append(li);
		});
	};

	var debounce = null;	//Track times between input events to search
	var lastSearch = null;
	search.addEventListener('input', (e) => {
		let target = e.target;

		const millis = Date.now();
		debounce = millis;

		setTimeout(() => {
			if(debounce == millis) {
				//Only continue with ajax request if there hasn't been a keyevent in the last quarter second
				let query = encodeURIComponent(target.value);
				if(query != '') {
					ToggleLoader();
					const time = Date.now();
					lastSearch = time;
					fetch(`/search/${query}`)
						.then((resp) => {
							if(lastSearch == time) {
								return resp;
							}
							throw new Error('Not the most recent request');
						})
						.catch((e) => null)
						.then((resp) => resp.json())
						.then(UpdateSearchResults)
						.catch((e) => console.error(e));
				}
			}
		}, 200);
	});

	search.addEventListener('focusout', (e) => {
		resultList.classList.add('hide');
	});

	search.addEventListener('focusin', (e) => {
		if(resultList.innerHTML !== '') {
			resultList.classList.remove('hide');
		}
	});

	resultList.addEventListener('mousemove', (e) => {
		let target = e.target;
		if(target.tagName === 'LI') {
			let currFocus = document.querySelector('li[data-focus=true]');
			currFocus.attributes['data-focus'].value = false;
			currFocus.ariaSelected = false;
			target.attributes['data-focus'].value = true;
			target.ariaSelected = true;
		}
	});

	resultList.addEventListener('mousedown', (e) => {
		let target = e.target;
		if(target.tagName == 'LI') {
			//target.style.backgroundColor = '#bbb';
			NavigateResult(target);
		}
	});

	/**
	* Use arrow keys to navigate focus between search elements
	*/
	search.addEventListener('keydown', (e) => {
		let key = e.key;
		let targetFocus = null;
		let currFocus = document.querySelector('li[data-focus=true]');
		if(key === 'ArrowUp') {
			e.preventDefault();
			targetFocus = currFocus.previousElementSibling;
			if(targetFocus == null) {
				targetFocus = document.querySelector('ul#search-results li:last-child');
			}
		} else if(key === 'ArrowDown') {
			e.preventDefault();
			targetFocus = currFocus.nextElementSibling;
			if(targetFocus == null) {
				targetFocus = document.querySelector('ul#search-results li:first-child');
			}
		}

		if(targetFocus != null) {
			currFocus.attributes['data-focus'].value = false;
			currFocus.ariaSelected = false;
			targetFocus.attributes['data-focus'].value = true;
			targetFocus.ariaSelected = true;
		}

		if(key === 'Enter') {
			NavigateResult(currFocus);
		}
		if(key === 'Escape') {
			ToggleResults();
		}

	});


});
