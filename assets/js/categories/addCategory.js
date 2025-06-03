// Button Add Category
document.addEventListener("DOMContentLoaded", function() {
    const addCategoryBtn = document.getElementById('addCategoryBtn');
    if (addCategoryBtn) {
        addCategoryBtn.addEventListener('click', function() {
            const categoryName = document.getElementById('categoryName').value;
            
            if (!categoryName.trim()) {
                alert('Vui lòng nhập tên danh mục');
                return;
            }
            
            fetch('/admin/categories/create', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ categoryName: categoryName })
            })
            .then(response => {
                if (response.ok) {
                    return response.text();
                }
                throw new Error('Lỗi khi tạo danh mục');
            })
            .then(data => {
                alert('Tạo danh mục thành công');
                window.location.reload();
            })
            .catch(error => {
                console.error("Error:", error);
                alert(error.message);
            });
        });
    } else {
    }
});