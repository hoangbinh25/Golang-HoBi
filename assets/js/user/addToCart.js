document.addEventListener("DOMContentLoaded", () => {
  const cartCount = document.getElementById("cart-count");

  document.querySelectorAll(".add-to-cart").forEach(btn => {
    btn.addEventListener("click", async () => {
      const productId = btn.getAttribute("data-id");

      try {
        const res = await fetch("/cart/add", {
          method: "POST",
          headers: {
            "Content-Type": "application/x-www-form-urlencoded"
          },
          body: `product_id=${productId}`
        });

        if (res.ok) {
          const data = await res.json();

          // Cập nhật số lượng trên giỏ
          cartCount.textContent = data.total_quantity;

          // (Tùy chọn) alert hoặc toast
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
