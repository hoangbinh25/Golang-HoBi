// Button Delete Product
function confirmDelete(productId) {
    if (confirm("Are you sure you want to delete this product?")) {
        fetch(`/admin/products/delete?id=${productId}`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ id: productId }),
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