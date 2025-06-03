// Button Add Product
window.onload = function () {
    const addBtn = document.getElementById('addProductBtn');
    if (addBtn) {
        addBtn.addEventListener('click', function (e) {
            e.preventDefault();

            const productName = document.getElementById('productName').value;
            const description = document.getElementById('description').value;
            const oldPrice = document.getElementById('oldPrice').value;
            const currentPrice = document.getElementById('currentPrice').value;
            const quantity = document.getElementById('quantity').value;
            const fileInput = document.querySelector("#uploadImage");
            const category_id = document.getElementById('category_id').value;

            // create object data
            const formData = new FormData();
            formData.append('productName', productName);
            formData.append('description', description);
            formData.append('quantity', quantity);
            formData.append('uploadImage', fileInput.files[0]);
            formData.append('oldPrice', oldPrice)
            formData.append('currentPrice', currentPrice);
            formData.append('category_id', category_id);

            fetch('/admin/products/create', {
                method: 'POST',
                body: formData
            })
                .then(response => {
                    if (!response.ok) {
                        return response.text().then(err => {
                            throw new Error("Server error: " + err);
                        });
                    }
                    return response.text();
                    
                })
                .then(data => {
                    alert('Product created successfully');
                    window.location.href = "/admin/products";
                })
                .catch(error => {
                    alert(error.message)
                });
        });
    }
}

// Upload Image Product
const uploadImage = document.getElementById('uploadImage');
if (uploadImage) {
    uploadImage.addEventListener('change', function () {
        const file = this.files[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = function () {
                const img = document.getElementById("previewImage");
                img.src = reader.result;
            }
            reader.readAsDataURL(file);
        }
    })
}


// Preview Image Product
function PreviewImage(event) {
    var reader = new FileReader();
    reader.onload = function () {
        var output = document.getElementById('previewImage');
        output.src = reader.result;
        output.style.display = "block"
    };
    reader.readAsDataURL(event.target.files[0]);
}