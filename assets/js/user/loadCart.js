document.addEventListener("DOMContentLoaded", function () {
  const cartIcon = document.getElementById("cart-icon");
  const cartPopup = document.getElementById("cart-popup");
  const cartItems = document.getElementById("cart-items");
  const cartCount = document.getElementById("cart-count");
  const cartSummary = document.getElementById("cart-summary");

  // Hiển thị popup khi hover
  cartIcon.addEventListener("mouseenter", () => {
    cartPopup.style.display = "block";
    fetch("/cart/data")
      .then(res => res.json())
      .then(data => {
        cartItems.innerHTML = "";
        let totalQty = 0;

        data.forEach(item => {
          totalQty += item.quantity;

          const itemEl = document.createElement("div");
          itemEl.className = "cart-item";
          itemEl.innerHTML = `
            <img src="/${item.image}" alt="${item.name}">
            <div class="flex-grow-1">
              <div class="text-truncate">${item.name}</div>
              <div class="text-danger fw-bold">${item.price.toLocaleString()} ₫</div>
            </div>
          `;
          cartItems.appendChild(itemEl);
        });

        cartCount.textContent = totalQty;
        cartSummary.textContent = `${totalQty} Thêm Hàng Vào Giỏ`;
      });
  });

  // Ẩn popup khi chuột rời đi
  cartIcon.addEventListener("mouseleave", () => {
    setTimeout(() => {
      cartPopup.style.display = "none";
    }, 500);
  });
});