package menu

type Menu struct {
	RestaurantID int           `json:"restaurantId,omitempty"`
	Entrees      []Entree      `json:"entrees,omitempty"`
	Sides        []Side        `json:"sides,omitempty"`
	Drinks       []Drink       `json:"drinks,omitempty"`
	NonFoodItems []NonFoodItem `json:"nonFoodItems,omitempty"`
}
