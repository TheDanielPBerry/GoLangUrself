window.addEventListener('load', () => {
	let passwordField = document.getElementById('password');
	let passwordConfirmedField = document.getElementById('password_confirmation');
	let validateMsg = document.getElementById('validate-msg');

	const CheckPasswordConfirmation = () => {
		if(passwordField.value.length < 10) {
			return;
		}

		if(passwordField.value != passwordConfirmedField.value) {
			validateMsg.innerText = 'Password fields must match';
			validateMsg.classList.remove('hide');
		} else {
			validateMsg.classList.add('hide');
		}
	};

	passwordField.addEventListener('focusout', (e) => {
		let target = e.target;
		if(target.value.length < 10) {
			validateMsg.innerText = 'Password must be at least 10 characters in length';
			validateMsg.classList.remove('hide');
		} else {
			validateMsg.classList.add('hide');
			CheckPasswordConfirmation();
		}
	});

	passwordConfirmedField.addEventListener('focusout', (e) => 	CheckPasswordConfirmation());
});
