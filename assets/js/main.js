function increaseQuantity(button) {
    const group = button.closest(".row");
    const quantityInput = group.querySelector(".quantity");
    const productIdInput = group.querySelector(".productId");
    if (!productIdInput) {
        console.error("❌ Không tìm thấy .productId trong dòng");
        return;
    };
    const productId = productIdInput.value;

    const stock = parseInt(group.querySelector(".stock-data").value || "0");
    const sold = parseInt(group.querySelector(".sold-data").value || "0");
    const current = parseInt(quantityInput.value || "1");

    if (current + sold >= stock) {
        alert("Vượt quá số lượng tồn kho!");
        return;
    }

    updateQuantity(productId, current + 1)
    quantityInput.value = current + 1;
}

function decreaseQuantity(button) {
    const group = button.closest(".row");
    const quantityInput = group.querySelector(".quantity");
    const current = parseInt(quantityInput.value || "1");
    const productIdInput = group.querySelector(".productId");
    if (!productIdInput) return;

    const productId = productIdInput.value;
    if (current > 1) {
        quantityInput.value = current - 1;
        updateQuantity(productId, current - 1)
    }
}

function updateQuantity(productId, quantity) {
    fetch("/cart/update", {
        method: "POST",
        headers: {
            "Content-Type": "application/x-www-form-urlencoded",
        },

        body: `product_id=${productId}&quantity=${quantity}`,

    }).then(res => {
        if (!res.ok) {
            alert('Không thể cập nhật số lượng')
        }
    });
}

function deleteCartItem(button) {
    const row = button.closest(".row");
    const productIdInput = row.querySelector(".productId");

    const productId = productIdInput.value;

    if (confirm("Bạn có chắc muốn xóa sản phẩm khỏi giỏ hàng?")) {
        fetch("/cart/delete", {
            method: "POST",
            headers: {
                "Content-Type": "application/x-www-form-urlencoded",
            },
            body: `product_id=${productId}`
        }).then(res => {
            if (res.ok) {
                row.remove();
            } else {
                alert("Xóa sản phẩm không thành công");
            }
        });
    }
}

function checkoutCartItems() {

}
