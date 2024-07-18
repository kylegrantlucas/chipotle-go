package restaurant

type Marketing struct {
	OperationsMarket                  string `json:"operationsMarket,omitempty"`
	SpecialMenuPanelInstructions      string `json:"specialMenuPanelInstructions,omitempty"`
	FeatureMenuPanel                  string `json:"featureMenuPanel,omitempty"`
	KidsMenuPanel                     string `json:"kidsMenuPanel,omitempty"`
	CaloriesOnMenuPanel               string `json:"caloriesOnMenuPanel,omitempty"`
	FoodWithIntegrityMenuBoardWidthID string `json:"foodWithIntegrityMenuBoardWidthId,omitempty"`
	MenuBoardPanelHeightID            string `json:"menuBoardPanelHeightId,omitempty"`
	MenuPanelTypeID                   string `json:"menuPanelTypeId,omitempty"`
	AlcoholCategory                   string `json:"alcoholCategory,omitempty"`
	AlcoholCategoryDescription        string `json:"alcoholCategoryDescription,omitempty"`
	MarketingAlcoholDescription       string `json:"marketingAlcoholDescription,omitempty"`
}
