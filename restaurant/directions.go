package restaurant

type Directions struct {
	Landmark           string `json:"landmark,omitempty"`
	CrossStreet1       string `json:"crossStreet1,omitempty"`
	CrossStreet2       string `json:"crossStreet2,omitempty"`
	PickupInstructions string `json:"pickupInstructions,omitempty"`
}
