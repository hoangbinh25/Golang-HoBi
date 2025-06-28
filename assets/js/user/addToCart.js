document.addEventListener("DOMContentLoaded", () => {
  const cartCount = document.getElementById("cart-count");

  document.querySelectorAll(".add-to-cart").forEach(btn => {
    btn.addEventListener("click", async (e) => {
      e.preventDefault(); // tránh reload nếu trong form

      const productId = btn.getAttribute("data-id");
      const quantityInput = document.querySelector(".quantity");
      const quantity = quantityInput ? quantityInput.value : 1;

      try {
        const res = await fetch("/cart/add", {
          method: "POST",
          headers: {
            "Content-Type": "application/x-www-form-urlencoded"
          },
          body: `product_id=${productId}&quantity=${quantity}`
        });

        const text = await res.text();
        let data = {};
        try {
          data = JSON.parse(text);
        } catch (err) {
          console.warn("Không phải JSON:", text);
        }

        if (res.ok) {
          cartCount.textContent = data.total_quantity || "?";
          alert("✅ Đã thêm vào giỏ hàng!");
        } else if (res.status === 401) {
          alert("⚠️ Bạn cần đăng nhập để thêm vào giỏ.");
          window.location.href = "/login";
        } else {
          alert("❌ Lỗi khi thêm vào giỏ hàng.");
        }
      } catch (err) {
        alert("❌ Lỗi kết nối tới server.");
        console.error(err);
      }
    });
  });
});
