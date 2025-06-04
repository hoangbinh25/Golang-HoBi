package controllers

import (
	"net/http"

	"github.com/Golang-Shoppe/controllers/admincontroller"
	cartcontroller "github.com/Golang-Shoppe/controllers/cartController"
	categorycontroller "github.com/Golang-Shoppe/controllers/categoryController"
	checkoutcontroller "github.com/Golang-Shoppe/controllers/checkoutController"
	"github.com/Golang-Shoppe/controllers/oauth"
)

func InitializersRoutes() {
	// Serve file css
	//Client
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("assets/uploads"))))

	http.HandleFunc("/web/assets/css/base.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "assets/css/base.css")
	})

	http.HandleFunc("/web/assets/css/main.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "assets/css/main.css")
	})

	http.HandleFunc("/web/assets/css/responsive.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "assets/css/responsive.css")
	})

	http.HandleFunc("/web/assets/css/common.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "assets/css/common.css")
	})

	// admin routes
	// Categories
	http.HandleFunc("/admin", admincontroller.AdminHomeHandler)
	http.HandleFunc("/admin/categories", admincontroller.Index)
	http.HandleFunc("/admin/categories/create", admincontroller.CreateCategory)
	http.HandleFunc("/admin/categories/update", admincontroller.UpdateCategory)
	http.HandleFunc("/admin/categories/delete", admincontroller.DeleteCategory)

	// Products
	http.HandleFunc("/admin/products", admincontroller.ProductHome)
	http.HandleFunc("/admin/products/create", admincontroller.CreateProduct)
	http.HandleFunc("/admin/products/update", admincontroller.UpdateProduct)
	http.HandleFunc("/admin/products/delete", admincontroller.DeleteProduct)

	// Orders
	http.HandleFunc("/admin/orders", admincontroller.OrderHome)

	// Routes
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/product/detail", admincontroller.DetailProductHandler)

	// login routes
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/loginauth", LoginAuthHandler)
	http.HandleFunc("/auth/google", oauth.HandleGoogleLogin)
	http.HandleFunc("/auth/callback", oauth.HandleGoogleCallback)

	http.HandleFunc("/logout", LogoutHandler)

	// register routes
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/registerauth", RegisterAuthHandler)

	// profile user routes
	http.HandleFunc("/user/profile", ProfileUserHandler)
	http.HandleFunc("/user/profile/update", ProfileUserUpdateHandler)
	http.HandleFunc("/user/profile/updatepassword", ProfileUserUpdatePasswordHandler)

	// Cart
	http.HandleFunc("/cart", cartcontroller.ViewCartHandler)
	http.HandleFunc("/cart/add", cartcontroller.AddToCartHandler)
	http.HandleFunc("/cart/data", cartcontroller.GetCartHandle)
	http.HandleFunc("/cart/update", cartcontroller.UpdateCartQuantityHandle)
	http.HandleFunc("/cart/delete", cartcontroller.DeleteCartItem)

	// Checkout
	http.HandleFunc("/checkout", checkoutcontroller.ShowCheckoutPage)
	http.HandleFunc("/checkout/submit", checkoutcontroller.CheckOutCartHandle)

	http.HandleFunc("/category", categorycontroller.ShowProductsCategory)

}
