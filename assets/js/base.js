window.addEventListener('load', () => {
	document.getElementById('account-menu').addEventListener('click', (e) => {
		let target = e.target;
		target.nextElementSibling.classList.toggle('hide');
	});
});
