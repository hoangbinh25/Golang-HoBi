function increaseQuantity(button) {
    const group = button.closest(".input-group"); // dùng .input-group là cha gần nhất
    if (!group) {
        console.error("❌ Không tìm thấy .input-group");
        return;
    }

    const input = group.querySelector(".quantity");
    const stock = parseInt(document.getElementById("stockQuantity")?.value || "0");
    const sold = parseInt(document.getElementById("soldQuantity")?.getAttribute("data-sold") || "0");
    const productId = button.closest(".row")?.querySelector(".productId")?.value;
    const current = parseInt(input.value || "1");

    if (current >= (stock - sold)) {
        alert("⚠️ Vượt quá số lượng tồn kho!");
        return;
    }

    input.value = current + 1;
    updateQuantity(productId, current + 1);
}

function decreaseQuantity(button) {
    const group = button.closest(".input-group");
    if (!group) return;

    const input = group.querySelector(".quantity");
    const productId = button.closest(".row")?.querySelector(".productId")?.value;
    const current = parseInt(input.value || "1");


    input.value = current - 1;
    updateQuantity(productId, current - 1);

}

function setAction(type) {
    document.getElementById("action_type").value = type;
}


function updateQuantity(productId, quantity) {
    fetch("/cart/update", {
        method: "POST",
        headers: {
            "Content-Type": "application/x-www-form-urlencoded",
        },

        body: `productId=${productId}&quantity=${quantity}`,

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

document.addEventListener("DOMContentLoaded", function () {
    const addToCartBtn = document.querySelector('.product-actions .add-to-cart');
    if (addToCartBtn) {
        addToCartBtn.addEventListener('click', function (e) {
            e.preventDefault();
            const quantityInput = document.querySelector('.input-group .quantity');
            const productIdInput = document.getElementById('productId');
            const quantity = quantityInput ? parseInt(quantityInput.value) : 1;
            const productId = productIdInput ? productIdInput.value : null;

            if (!productId) {
                alert("Không tìm thấy sản phẩm!");
                return;
            }

            // Gửi request AJAX thay vì chuyển trang
            fetch('/cart/add', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                },
                body: `product_id=${productId}&quantity=${quantity}`
            })
                .then(res => res.json())
                .then(data => {
                    if (data.success) {
                        alert('Đã thêm vào giỏ hàng!');
                    } else {
                        alert(data.message || 'Thêm vào giỏ hàng thất bại!');
                    }
                })
                .catch(() => {
                    alert('Lỗi khi thêm vào giỏ hàng!');
                });
        });
    }
});