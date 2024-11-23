window.addEventListener('load', () => {
	let accountMenu = document.getElementById('account-menu');
	if(accountMenu) {
		accountMenu.addEventListener('click', (e) => {
			let target = e.target;
			target.nextElementSibling.classList.toggle('hide');
		});
	}
});
