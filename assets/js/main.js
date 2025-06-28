function increaseQuantity(button) {
    const group = button.closest(".input-group"); // dùng .input-group là cha gần nhất
    if (!group) {
        console.error("❌ Không tìm thấy .input-group");
        return;
    }

    const input = group.querySelector(".quantity");
    const stock = parseInt(group.querySelector(".stock-data")?.value || "0");
    const sold = parseInt(group.querySelector(".sold-data")?.value || "0");
    const productId = button.closest(".col")?.querySelector(".productId")?.value;
    const current = parseInt(input.value || "1");

    if (current + sold >= stock) {
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
    const productId = button.closest(".col")?.querySelector(".productId")?.value;
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

document.addEventListener("DOMContentLoaded", function () {
    // Bắt sự kiện cho nút "Mua ngay"
    const buyNowBtn = document.querySelector('.product-actions .btn.btn-danger:last-child');
    if (buyNowBtn) {
        buyNowBtn.addEventListener('click', function (e) {
            e.preventDefault();
            // Lấy số lượng hiện tại
            const quantityInput = document.querySelector('.input-group .quantity');
            const productIdInput = document.querySelector('.input-group .productId');
            const quantity = quantityInput ? parseInt(quantityInput.value) : 1;
            const productId = productIdInput ? productIdInput.value : null;

            if (!productId) {
                alert("Không tìm thấy sản phẩm!");
                return;
            }
            console.log(productId);


            // Chuyển hướng sang trang giỏ hàng/checkout kèm số lượng
            window.location.href = `/cart/add?product_id=${productId}&quantity=${quantity}`;
        });
    }
});