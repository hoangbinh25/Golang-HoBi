function showPassword(iconSpan, inputId) {
    const passwordInput = document.getElementById(inputId);
    const toggleIcon = iconSpan.querySelector('i');
    if (passwordInput.type === 'password') {
        passwordInput.type = 'text';
        toggleIcon.classList.remove('fa-eye-slash');
        toggleIcon.classList.add('fa-eye');
    } else {
        passwordInput.type = 'password';
        toggleIcon.classList.remove('fa-eye');
        toggleIcon.classList.add('fa-eye-slash');
    }
}
// Button Change Password
document.getElementById("changePasswordForm").addEventListener("submit", function (e) {
    e.preventDefault(); // Ngăn submit form gốc

    const currentPassword = document.getElementById('currentPassword').value;
    const newPassword = document.getElementById('newPassword').value;
    const confirmPassword = document.getElementById('confirmPassword').value;

    if (newPassword !== confirmPassword) {
        alert("Mật khẩu không khớp");
        return;
    } else if (newPassword.length < 8) {
        alert("Mật khẩu phải có ít nhất 8 ký tự");
        return;
    } else if (!/[A-Z]/.test(newPassword)) {
        alert("Mật khẩu phải có ít nhất 1 chữ cái viết hoa");
        return;
    } else if (!/\d/.test(newPassword)) {
        alert("Mật khẩu phải có ít nhất 1 chữ số");
        return;
    }

    const formData = new FormData();
    formData.append('currentPassword', currentPassword);
    formData.append('newPassword', newPassword);
    formData.append('confirmPassword', confirmPassword);

    fetch('/user/profile/updatepassword', {
        method: 'POST',
        headers: {
            "Content-Type": "application/x-www-form-urlencoded",
        },
        body: formData,
    })
        .then(response => {
            if (response.ok) {
                return response.text();
            } else {
                return response.text().then(text => { throw new Error(text) });
            }
        })
        .then(data => {
            alert('Mật khẩu đã được thay đổi thành công');
            window.location.href = "/user/profile";
        })
        .catch(error => {
            console.error('Lỗi:', error);
            alert(error.message);
        });
});

