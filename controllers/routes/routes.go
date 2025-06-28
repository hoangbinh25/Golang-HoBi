package controllers

import (
	"net/http"

	admincontroller "github.com/Golang-Shoppe/controllers/adminController"
	cartcontroller "github.com/Golang-Shoppe/controllers/cartController"
	categorycontroller "github.com/Golang-Shoppe/controllers/categoryController"
	checkoutcontroller "github.com/Golang-Shoppe/controllers/checkoutController"
	"github.com/Golang-Shoppe/controllers/oauth"
	usercontroller "github.com/Golang-Shoppe/controllers/userController"
	"github.com/Golang-Shoppe/initializers"
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

	http.HandleFunc("/health", initializers.HealthCheck)

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
	http.HandleFunc("/admin/orders/confirm", admincontroller.ConfirmOrder)

	// Routes
	http.HandleFunc("/", admincontroller.AdminHomeHandler)
	http.HandleFunc("/product/detail", admincontroller.DetailProductHandler)

	// login routes
	http.HandleFunc("/login", usercontroller.LoginHandler)
	http.HandleFunc("/loginauth", usercontroller.LoginAuthHandler)
	http.HandleFunc("/auth/google", oauth.HandleGoogleLogin)
	http.HandleFunc("/auth/callback", oauth.HandleGoogleCallback)

	http.HandleFunc("/logout", usercontroller.LogoutHandler)

	// register routes
	http.HandleFunc("/register", usercontroller.RegisterHandler)
	http.HandleFunc("/registerauth", usercontroller.RegisterAuthHandler)

	// Email verification and password reset routes
	http.HandleFunc("/verify-email", usercontroller.VerifyEmailHandler)
	http.HandleFunc("/forgot-password", usercontroller.ForgotPasswordHandler)
	http.HandleFunc("/reset-password", usercontroller.ResetPasswordHandler)

	// profile user routes
	http.HandleFunc("/user/profile", usercontroller.ProfileUserHandler)
	http.HandleFunc("/user/profile/update", usercontroller.ProfileUserUpdateHandler)
	http.HandleFunc("/user/profile/updatepassword", usercontroller.ProfileUserUpdatePasswordHandler)
	http.HandleFunc("/user/orders", usercontroller.UserOrdersHandler)
	http.HandleFunc("/user/confirm-received", usercontroller.ConfirmReceivedHandler)

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
