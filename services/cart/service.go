package cart

import (
	"fmt"

	"github.com/AhmedRabea0302/ecom/types"
)

func getCartItemsIDs(items []types.CartItem) ([]int, error) {
	productIDs := make([]int, len(items))

	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("Invalid quantity for the product: %d", item.ProductID)
		}

		productIDs[i] = item.ProductID
	}
	return productIDs, nil
}

func (h *Handler) createOrder(ps []types.Product, items []types.CartItem, userID int) (int, float64, error) {

	productMap := make(map[int]types.Product)

	for _, product := range ps {
		productMap[product.ID] = product
	}

	// Check if all products are already in stock
	if err := checkIfCartsInStock(items, productMap); err != nil {
		return 0, 0, nil
	}

	// Calculate The Total Price
	totalPrice := calculatePrice(items, productMap)

	// Reduce Quantity of products in our db
	for _, item := range items {
		product := productMap[item.ProductID]
		product.Quantity -= item.Quantity

		h.productStore.UpdateProduct(product)
	}
	// Create Order
	orderID, err := h.store.CreateOrder(types.Order{
		UserID:  userID,
		Total:   totalPrice,
		Status:  "pending",
		Address: "Some Address",
	})

	if err != nil {
		return 0, 0, err
	}

	// Create Order Items
	for _, item := range items {
		h.store.CreateOrderItem(types.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productMap[item.ProductID].Price,
		})
	}

	return orderID, totalPrice, nil
}

func checkIfCartsInStock(cartItems []types.CartItem, products map[int]types.Product) error {
	if len(cartItems) == 0 {
		return fmt.Errorf("cart is empty")
	}

	for _, item := range cartItems {
		product, ok := products[item.ProductID]
		if !ok {
			return fmt.Errorf("product %d is not available in the stock, please refresh your cart", item.ProductID)
		}

		if product.Quantity < item.Quantity {
			return fmt.Errorf("product %s of id %d is out of stock, only %d available", product.Name, item.ProductID, product.Quantity)
		}
	}

	return nil
}

func calculatePrice(cartItems []types.CartItem, products map[int]types.Product) float64 {
	var total float64

	for _, item := range cartItems {
		product := products[item.ProductID]
		total += float64(item.Quantity) * product.Price
	}

	return total
}
