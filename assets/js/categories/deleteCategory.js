// Button Delete Category
function deleteCategory(categoryId) {
    if (confirm("Are you sure you want delete this category")) {
        fetch(`/admin/categories/delete?id=${categoryId}`, {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify({id: categoryId})
        })
        .then(res => {
            if (res.ok) {
                window.location.reload();
            } else {
                alert("Delete failed")
            }
        })
        .catch(error => console.error("Error: ", error))
    }
}