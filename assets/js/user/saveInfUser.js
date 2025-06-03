// Button Save Information User
document.addEventListener('DOMContentLoaded', () => {
    const btnSaveInfo = document.getElementById('btnSaveInfo');
    if (btnSaveInfo) {
        btnSaveInfo.addEventListener('click', function() {
            const name = document.getElementById('name').value;
            const email = document.getElementById('email').value;
            const phone = document.getElementById('phone').value;
            const gender = document.getElementById('gender').value;
            const address = document.getElementById('address').value;

            const formData = new FormData();
            formData.append('name', name);
            formData.append('email', email);
            formData.append('phone', phone);
            formData.append('gender', gender);
            formData.append('address', address);

            fetch('/user/profile/update', {
                method: 'POST',
                body: formData,
            })
            .then(response => {
                if (response.redirected) {
                    window.location.href = response.url;
                    return;
                } else {
                    throw new Error('Error updating user information')
                }
            })
            .then(data => {
                if (data.success) {
                    alert('User information updated successfully');
                    window.location.href = "/user/profile";
                }
                    
            })
            .catch(error => {
                console.error('Error:', error);
                alert(error.message);
            });
        });
    }
})

