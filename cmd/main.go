package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/kylegrantlucas/chipotle-go"
	"github.com/kylegrantlucas/chipotle-go/menu"
	"github.com/kylegrantlucas/chipotle-go/restaurant"
	"github.com/kylegrantlucas/chipotle-go/search"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	client := chipotle.NewClient("INSERT_YOUR_API_KEY_HERE")

	query := search.Query{
		Latitude:           38.495693700000004,
		Longitude:          -121.19452040000002,
		Radius:             9046700,
		RestaurantStatuses: []string{"OPEN", "LAB"},
		ConceptIds:         []string{"CMG"},
		OrderBy:            "distance",
		OrderByDescending:  false,
		PageSize:           4000,
		PageIndex:          0,
		Embeds: search.Embeds{
			AddressTypes:   []string{"MAIN"},
			RealHours:      true,
			Directions:     true,
			Catering:       true,
			OnlineOrdering: true,
			Timezone:       true,
			Marketing:      true,
			Chipotlane:     true,
			Sustainability: true,
			Experience:     true,
		},
	}

	fmt.Println("Searching for restaurants...")
	result, err := client.Search(query)
	if err != nil {
		log.Fatal(err)
	}

	// some stats
	fmt.Printf("Total restaurants: %d\n", len(result.Restaurants))

	// drop the old database, we don't care if it doesn't exist, so ignore that class of error
	err = os.Remove("./chipotle.db")
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}

	// Open the database connection
	fmt.Println("Opening database connection...")
	db, err := sql.Open("sqlite3", "./chipotle.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Set some DB pragmas to speed this up
	fmt.Println("Setting database pragmas to speed up insertion...")
	_, err = db.Exec("PRAGMA synchronous = OFF")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("PRAGMA journal_mode = OFF")
	if err != nil {
		log.Fatal(err)
	}

	// Create tables
	fmt.Println("Creating tables...")
	createTables(db)

	// Configurable running thread limit for fetching menus
	fetchThreadLimit := 75

	// Create channels for restaurant processing
	restaurantChan := make(chan restaurant.Restaurant)

	menus := []*menu.Menu{}

	// Wait groups for the two stages
	var fetchWg sync.WaitGroup

	// Mutex for database operations
	var menuMutex sync.Mutex
	var dbMutex sync.Mutex

	// log the start
	fmt.Println("Fetching menus and inserting restaurants...")

	// Start goroutines for fetching menus
	for i := 0; i < fetchThreadLimit; i++ {
		fetchWg.Add(1)
		go func() {
			defer fetchWg.Done()
			for r := range restaurantChan {
				// insert the restaurant into the database
				dbMutex.Lock()
				err := insertRestaurant(db, r)
				if err != nil {
					log.Fatal(err)
				}
				dbMutex.Unlock()

				m, err := client.GetMenu(r.RestaurantNumber)
				if err != nil {
					log.Printf("Error getting menu for restaurant %s: %v\n", r.RestaurantName, err)
					continue
				}

				menuMutex.Lock()
				menus = append(menus, m)
				menuMutex.Unlock()
			}
		}()
	}

	// Send restaurants to the fetching pool
	for _, restaurant := range result.Restaurants {
		restaurantChan <- restaurant
	}

	// Close the restaurant channel to signal fetch goroutines to finish
	close(restaurantChan)
	fetchWg.Wait()

	oi := optimizeItems(menus)

	// Insert optimized items into the database
	fmt.Println("Inserting optimized items into the database...")
	err = insertOptimizedItems(db, oi)
	if err != nil {
		log.Fatal(err)
	}

	// Insert menus into the database
	fmt.Println("Inserting menus into the database...")
	for _, m := range menus {
		err = insertMenu(db, m, oi)
		if err != nil {
			log.Fatal(err)
		}
	}

	// reset the pragmas
	fmt.Println("Resetting database PRAGMAs...")
	_, err = db.Exec("PRAGMA synchronous = FULL")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("PRAGMA journal_mode = DELETE")
	if err != nil {
		log.Fatal(err)
	}
}

func createTables(db *sql.DB) error {
	// Create Menu table
	createMenuTable := `
	CREATE TABLE IF NOT EXISTS menus (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		restaurant_id INTEGER
	);`

	// Create ItemTypes table
	createItemTypesTable := `
	CREATE TABLE IF NOT EXISTS item_types (
		id INTEGER PRIMARY KEY,
		item_type TEXT UNIQUE
	);`

	// Create ItemCategories table
	createItemCategoriesTable := `
	CREATE TABLE IF NOT EXISTS item_categories (
		id INTEGER PRIMARY KEY,
		item_category TEXT UNIQUE
	);`

	// Create ItemNames table
	createItemNamesTable := `
	CREATE TABLE IF NOT EXISTS item_names (
		id INTEGER PRIMARY KEY,
		item_name TEXT UNIQUE
	);`

	// Create PrimaryFillingNames table
	createPrimaryFillingNamesTable := `
	CREATE TABLE IF NOT EXISTS primary_filling_names (
		id INTEGER PRIMARY KEY,
		primary_filling_name TEXT UNIQUE
	);`

	// Create Items table
	createItemsTable := `
	CREATE TABLE IF NOT EXISTS items (
		id TEXT PRIMARY KEY,
		item_type_id INTEGER,
		item_category_id INTEGER,
		item_name_id INTEGER,
		primary_filling_name_id INTEGER,
		FOREIGN KEY(item_type_id) REFERENCES item_types(id),
		FOREIGN KEY(item_category_id) REFERENCES item_categories(id),
		FOREIGN KEY(item_name_id) REFERENCES item_names(id),
		FOREIGN KEY(primary_filling_name_id) REFERENCES primary_filling_names(id)
	);`

	// Create Entree table
	createEntreeTable := `
	CREATE TABLE IF NOT EXISTS entrees (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		menu_id INTEGER,
		item_id TEXT,
		pos_id INTEGER,
		unit_price REAL,
		unit_delivery_price REAL,
		unit_count INTEGER,
		max_quantity INTEGER,
		eligible_for_delivery BOOLEAN,
		max_contents INTEGER,
		max_customizations INTEGER,
		max_on_the_side_customizations INTEGER,
		max_extras INTEGER,
		max_halfs INTEGER,
		max_extras_plus_halfs INTEGER,
		is_universal BOOLEAN,
		is_item_available BOOLEAN,
		FOREIGN KEY(menu_id) REFERENCES menus(id),
		FOREIGN KEY(item_id) REFERENCES items(id)
	);`

	// Create EntreeContentGroups table
	createEntreeContentGroupsTable := `
	CREATE TABLE IF NOT EXISTS entree_content_groups (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		entree_id INTEGER,
		content_group_id INTEGER,
		min_quantity INTEGER,
		max_quantity INTEGER,
		FOREIGN KEY(entree_id) REFERENCES entrees(id),
		FOREIGN KEY(content_group_id) REFERENCES content_groups(id)
	);`

	// Create ContentGroups table
	createContentGroupsTable := `
	CREATE TABLE IF NOT EXISTS content_groups (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content_group_name TEXT UNIQUE
	);`

	// Create Contents table
	createContentsTable := `
	CREATE TABLE IF NOT EXISTS contents (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		entree_id INTEGER,
		item_id TEXT,
		pos_id INTEGER,
		unit_price REAL,
		unit_delivery_price REAL,
		unit_count INTEGER,
		eligible_for_delivery BOOLEAN,
		pricing_reference_item_id INTEGER,
		count_towards_customization_max BOOLEAN,
		count_towards_content_max BOOLEAN,
		content_group_id INTEGER,
		default_content BOOLEAN,
		is_item_available BOOLEAN,
		FOREIGN KEY(entree_id) REFERENCES entrees(id),
		FOREIGN KEY(item_id) REFERENCES items(id),
		FOREIGN KEY(content_group_id) REFERENCES content_groups(id)
	);`

	// Create Drink table
	createDrinkTable := `
	CREATE TABLE IF NOT EXISTS drinks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		menu_id INTEGER,
		item_id TEXT,
		pos_id INTEGER,
		unit_price REAL,
		unit_delivery_price REAL,
		unit_count INTEGER,
		max_quantity INTEGER,
		eligible_for_delivery BOOLEAN,
		is_universal BOOLEAN,
		is_item_available BOOLEAN,
		FOREIGN KEY(menu_id) REFERENCES menus(id),
		FOREIGN KEY(item_id) REFERENCES items(id)
	);`

	// Create NonFoodItem table
	createNonFoodItemTable := `
	CREATE TABLE IF NOT EXISTS non_food_items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		menu_id INTEGER,
		item_id TEXT,
		pos_id INTEGER,
		unit_price REAL,
		unit_delivery_price REAL,
		unit_count INTEGER,
		max_quantity INTEGER,
		eligible_for_delivery BOOLEAN,
		is_universal BOOLEAN,
		is_item_available BOOLEAN,
		FOREIGN KEY(menu_id) REFERENCES menus(id),
		FOREIGN KEY(item_id) REFERENCES items(id)
	);`

	// Create Side table
	createSideTable := `
	CREATE TABLE IF NOT EXISTS sides (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		menu_id INTEGER,
		item_id TEXT,
		pos_id INTEGER,
		unit_price REAL,
		unit_delivery_price REAL,
		unit_count INTEGER,
		max_quantity INTEGER,
		eligible_for_delivery BOOLEAN,
		is_universal BOOLEAN,
		is_item_available BOOLEAN,
		FOREIGN KEY(menu_id) REFERENCES menus(id),
		FOREIGN KEY(item_id) REFERENCES items(id)
	);`

	// Create Restaurant table
	createRestaurantTable := `
	CREATE TABLE IF NOT EXISTS restaurants (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		restaurant_number INTEGER,
		restaurant_name TEXT,
		restaurant_location_type TEXT,
		restaurant_status TEXT,
		open_date TEXT,
		real_estate_category TEXT,
		operational_region TEXT,
		operational_sub_region TEXT,
		operational_patch TEXT,
		designated_market_area_name TEXT,
		distance REAL,
		directions_landmark TEXT,
		directions_cross_street1 TEXT,
		directions_cross_street2 TEXT,
		directions_pickup_instructions TEXT,
		timezone_current_timezone_offset INTEGER,
		timezone_timezone_offset INTEGER,
		timezone_timezone TEXT,
		timezone_timezone_id TEXT,
		timezone_observe_daylight_savings TEXT,
		timezone_daylight_savings_offset INTEGER,
		marketing_operations_market TEXT,
		marketing_special_menu_panel_instructions TEXT,
		marketing_feature_menu_panel TEXT,
		marketing_kids_menu_panel TEXT,
		marketing_calories_on_menu_panel TEXT,
		marketing_food_with_integrity_menu_board_width_id TEXT,
		marketing_menu_board_panel_height_id TEXT,
		marketing_menu_panel_type_id TEXT,
		marketing_alcohol_category TEXT,
		marketing_alcohol_category_description TEXT,
		marketing_marketing_alcohol_description TEXT,
		catering_enabled BOOLEAN,
		chipotlane_pickup_enabled BOOLEAN,
		experience_curbside_pickup_enabled BOOLEAN,
		experience_dining_room_open BOOLEAN,
		experience_digital_kitchen BOOLEAN,
		experience_walkup_window_enabled BOOLEAN,
		experience_pickup_inside_enabled BOOLEAN,
		experience_crew_tip_pickup_enabled BOOLEAN,
		experience_crew_tip_delivery_enabled BOOLEAN,
		experience_context_rest_exp_enabled BOOLEAN,
		sustainability_utensils_default_state TEXT,
		planned_subs_compl_date TEXT,
		actual_subs_compl_date TEXT,
		online_ordering_enabled BOOLEAN,
		online_ordering_dot_com_search_enabled TEXT,
		online_ordering_credit_cards_accepted BOOLEAN,
		online_ordering_gift_cards_accepted BOOLEAN,
		online_ordering_bulk_orders_accepted BOOLEAN,
		online_ordering_tax_assessed BOOLEAN,
		restaurant_terminal_site_id INTEGER
	);`

	// Create Address table
	createAddressTable := `
	CREATE TABLE IF NOT EXISTS addresses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		restaurant_id INTEGER,
		address_type TEXT,
		address_line1 TEXT,
		address_line2 TEXT,
		locality TEXT,
		administrative_area TEXT,
		postal_code TEXT,
		sub_administrative_area TEXT,
		country_code TEXT,
		latitude REAL,
		longitude REAL,
		accuracy_determination TEXT,
		FOREIGN KEY(restaurant_id) REFERENCES restaurants(id)
	);`

	// Create RealHours table
	createRealHoursTable := `
	CREATE TABLE IF NOT EXISTS real_hours (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		restaurant_id INTEGER,
		day_of_week TEXT,
		open_date_time TEXT,
		close_date_time TEXT,
		FOREIGN KEY(restaurant_id) REFERENCES restaurants(id)
	);`

	// Execute the table creation queries
	queries := []string{
		createMenuTable, createItemTypesTable, createItemCategoriesTable, createItemNamesTable, createPrimaryFillingNamesTable,
		createItemsTable, createEntreeTable, createEntreeContentGroupsTable, createContentGroupsTable, createContentsTable,
		createDrinkTable, createNonFoodItemTable, createSideTable, createRestaurantTable, createAddressTable, createRealHoursTable,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("error creating table: %v", err)
		}
	}

	return nil
}

func insertOptimizedItems(db *sql.DB, optimizedItems *optimizedItems) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	itemTypeStmt, err := tx.Prepare("INSERT INTO item_types (item_type, id) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("error preparing item type statement: %v", err)
	}
	defer itemTypeStmt.Close()

	for itemType, id := range optimizedItems.ItemTypes {
		_, err := itemTypeStmt.Exec(itemType, id)
		if err != nil {
			return fmt.Errorf("error inserting item type: %v, %v, %v", err, itemType, id)
		}
	}

	itemCategoryStmt, err := tx.Prepare("INSERT INTO item_categories (item_category, id) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("error preparing item category statement: %v", err)
	}
	defer itemCategoryStmt.Close()

	for item, id := range optimizedItems.ItemCategories {
		_, err := itemCategoryStmt.Exec(item, id)
		if err != nil {
			return fmt.Errorf("error inserting item category: %v", err)
		}
	}

	itemNameStmt, err := tx.Prepare("INSERT INTO item_names (item_name, id) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("error preparing item name statement: %v", err)
	}
	defer itemNameStmt.Close()

	for item, id := range optimizedItems.ItemNames {
		_, err := itemNameStmt.Exec(item, id)
		if err != nil {
			return fmt.Errorf("error inserting item name: %v", err)
		}
	}

	primaryFillingNameStmt, err := tx.Prepare("INSERT INTO primary_filling_names (primary_filling_name, id) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("error preparing primary filling name statement: %v", err)
	}
	defer primaryFillingNameStmt.Close()

	for item, id := range optimizedItems.PrimaryFillingNames {
		_, err := primaryFillingNameStmt.Exec(item, id)
		if err != nil {
			return fmt.Errorf("error inserting primary filling name: %v", err)
		}
	}

	itemStmt, err := tx.Prepare("INSERT INTO items (id, item_type_id, item_category_id, item_name_id, primary_filling_name_id) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("error preparing item statement: %v", err)
	}
	defer itemStmt.Close()

	for id, i := range optimizedItems.Items {
		_, err := itemStmt.Exec(id, i.Type, i.Category, i.Name, i.PrimaryFillingName)
		if err != nil {
			return fmt.Errorf("error inserting item: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func insertMenu(db *sql.DB, menu *menu.Menu, optimizedItems *optimizedItems) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	menuStmt, err := tx.Prepare("INSERT INTO menus (restaurant_id) VALUES (?)")
	if err != nil {
		return fmt.Errorf("error preparing menu statement: %v", err)
	}
	defer menuStmt.Close()

	result, err := menuStmt.Exec(menu.RestaurantID)
	if err != nil {
		return fmt.Errorf("error inserting menu: %v", err)
	}

	menuID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting menu ID: %v", err)
	}

	entreeStmt, err := tx.Prepare(`
		INSERT INTO entrees (menu_id, item_id, pos_id, unit_price, unit_delivery_price, unit_count, max_quantity, eligible_for_delivery,
		max_contents, max_customizations, max_on_the_side_customizations, max_extras, max_halfs, max_extras_plus_halfs, is_universal, is_item_available)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("error preparing entree statement: %v", err)
	}
	defer entreeStmt.Close()

	entreeContentGroupStmt, err := tx.Prepare(`
		INSERT INTO entree_content_groups (entree_id, content_group_id, min_quantity, max_quantity)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("error preparing entree content group statement: %v", err)
	}
	defer entreeContentGroupStmt.Close()

	contentStmt, err := tx.Prepare(`
		INSERT INTO contents (entree_id, item_id, pos_id, unit_price, unit_delivery_price, unit_count, eligible_for_delivery,
		pricing_reference_item_id, count_towards_customization_max, count_towards_content_max, content_group_id, default_content, is_item_available)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("error preparing content statement: %v", err)
	}
	defer contentStmt.Close()

	drinkStmt, err := tx.Prepare(`
		INSERT INTO drinks (menu_id, item_id, pos_id, unit_price, unit_delivery_price, unit_count, max_quantity, eligible_for_delivery, is_universal, is_item_available)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("error preparing drink statement: %v", err)
	}
	defer drinkStmt.Close()

	nonFoodItemStmt, err := tx.Prepare(`
		INSERT INTO non_food_items (menu_id, item_id, pos_id, unit_price, unit_delivery_price, unit_count, max_quantity, eligible_for_delivery, is_universal, is_item_available)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("error preparing non-food item statement: %v", err)
	}
	defer nonFoodItemStmt.Close()

	sideStmt, err := tx.Prepare(`
		INSERT INTO sides (menu_id, item_id, pos_id, unit_price, unit_delivery_price, unit_count, max_quantity, eligible_for_delivery, is_universal, is_item_available)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("error preparing side statement: %v", err)
	}
	defer sideStmt.Close()

	for _, entree := range menu.Entrees {
		result, err := entreeStmt.Exec(
			menuID, entree.ItemID, entree.PosID, entree.UnitPrice, entree.UnitDeliveryPrice, entree.UnitCount, entree.MaxQuantity, entree.EligibleForDelivery,
			entree.MaxContents, entree.MaxCustomizations, entree.MaxOnTheSideCustomizations, entree.MaxExtras, entree.MaxHalfs, entree.MaxExtrasPlusHalfs, entree.IsUniversal, entree.IsItemAvailable,
		)
		if err != nil {
			return fmt.Errorf("error inserting entree: %v", err)
		}

		entreeID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("error getting entree ID: %v", err)
		}

		for _, cg := range entree.ContentGroups {
			contentGroupID := optimizedItems.AddContentGroup(cg.ContentGroupName)
			_, err := entreeContentGroupStmt.Exec(entreeID, contentGroupID, cg.MinQuantity, cg.MaxQuantity)
			if err != nil {
				return fmt.Errorf("error inserting content group: %v", err)
			}
		}

		for _, content := range entree.Contents {
			contentGroupID := optimizedItems.AddContentGroup(content.ContentGroupName)
			_, err := contentStmt.Exec(
				entreeID, content.ItemID, content.PosID, content.UnitPrice, content.UnitDeliveryPrice, content.UnitCount, content.EligibleForDelivery,
				content.PricingReferenceItemID, content.CountTowardsCustomizationMax, content.CountTowardsContentMax, contentGroupID, content.DefaultContent, content.IsItemAvailable,
			)
			if err != nil {
				return fmt.Errorf("error inserting content: %v", err)
			}
		}
	}

	for _, drink := range menu.Drinks {
		_, err := drinkStmt.Exec(
			menuID, drink.ItemID, drink.PosID, drink.UnitPrice, drink.UnitDeliveryPrice, drink.UnitCount, drink.MaxQuantity, drink.EligibleForDelivery, drink.IsUniversal, drink.IsItemAvailable,
		)
		if err != nil {
			return fmt.Errorf("error inserting drink: %v", err)
		}
	}

	for _, nfi := range menu.NonFoodItems {
		_, err := nonFoodItemStmt.Exec(
			menuID, nfi.ItemID, nfi.PosID, nfi.UnitPrice, nfi.UnitDeliveryPrice, nfi.UnitCount, nfi.MaxQuantity, nfi.EligibleForDelivery, nfi.IsUniversal, nfi.IsItemAvailable,
		)
		if err != nil {
			return fmt.Errorf("error inserting non-food item: %v", err)
		}
	}

	for _, side := range menu.Sides {
		_, err := sideStmt.Exec(
			menuID, side.ItemID, side.PosID, side.UnitPrice, side.UnitDeliveryPrice, side.UnitCount, side.MaxQuantity, side.EligibleForDelivery, side.IsUniversal, side.IsItemAvailable,
		)
		if err != nil {
			return fmt.Errorf("error inserting side: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func insertRestaurant(db *sql.DB, restaurant restaurant.Restaurant) error {
	// Insert the restaurant
	result, err := db.Exec(`
		INSERT INTO restaurants (
			restaurant_number, restaurant_name, restaurant_location_type,
			restaurant_status, open_date, real_estate_category,
			operational_region, operational_sub_region, operational_patch,
			designated_market_area_name, distance, directions_landmark,
			directions_cross_street1, directions_cross_street2, directions_pickup_instructions,
			timezone_current_timezone_offset, timezone_timezone_offset, timezone_timezone,
			timezone_timezone_id, timezone_observe_daylight_savings, timezone_daylight_savings_offset,
			marketing_operations_market, marketing_special_menu_panel_instructions, marketing_feature_menu_panel,
			marketing_kids_menu_panel, marketing_calories_on_menu_panel,
			marketing_food_with_integrity_menu_board_width_id, marketing_menu_board_panel_height_id, marketing_menu_panel_type_id,
			marketing_alcohol_category, marketing_alcohol_category_description, marketing_marketing_alcohol_description,
			catering_enabled, chipotlane_pickup_enabled, experience_curbside_pickup_enabled,
			experience_dining_room_open, experience_digital_kitchen, experience_walkup_window_enabled,
			experience_pickup_inside_enabled, experience_crew_tip_pickup_enabled, experience_crew_tip_delivery_enabled,
			experience_context_rest_exp_enabled, sustainability_utensils_default_state,
			planned_subs_compl_date, actual_subs_compl_date, online_ordering_enabled,
			online_ordering_dot_com_search_enabled, online_ordering_credit_cards_accepted,
			online_ordering_gift_cards_accepted, online_ordering_bulk_orders_accepted,
			online_ordering_tax_assessed, restaurant_terminal_site_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
		restaurant.RestaurantNumber, restaurant.RestaurantName, restaurant.RestaurantLocationType,
		restaurant.RestaurantStatus, restaurant.OpenDate, restaurant.RealEstateCategory,
		restaurant.OperationalRegion, restaurant.OperationalSubRegion, restaurant.OperationalPatch,
		restaurant.DesignatedMarketAreaName, restaurant.Distance, restaurant.Directions.Landmark,
		restaurant.Directions.CrossStreet1, restaurant.Directions.CrossStreet2, restaurant.Directions.PickupInstructions,
		restaurant.Timezone.CurrentTimezoneOffset, restaurant.Timezone.TimezoneOffset, restaurant.Timezone.Timezone,
		restaurant.Timezone.TimezoneID, restaurant.Timezone.ObserveDaylightSavings, restaurant.Timezone.DaylightSavingsOffset,
		restaurant.Marketing.OperationsMarket, restaurant.Marketing.SpecialMenuPanelInstructions, restaurant.Marketing.FeatureMenuPanel,
		restaurant.Marketing.KidsMenuPanel, restaurant.Marketing.CaloriesOnMenuPanel,
		restaurant.Marketing.FoodWithIntegrityMenuBoardWidthID, restaurant.Marketing.MenuBoardPanelHeightID, restaurant.Marketing.MenuPanelTypeID,
		restaurant.Marketing.AlcoholCategory, restaurant.Marketing.AlcoholCategoryDescription, restaurant.Marketing.MarketingAlcoholDescription,
		restaurant.Catering.CateringEnabled, restaurant.Chipotlane.ChipotlanePickupEnabled, restaurant.Experience.CurbsidePickupEnabled,
		restaurant.Experience.DiningRoomOpen, restaurant.Experience.DigitalKitchen, restaurant.Experience.WalkupWindowEnabled,
		restaurant.Experience.PickupInsideEnabled, restaurant.Experience.CrewTipPickupEnabled, restaurant.Experience.CrewTipDeliveryEnabled,
		restaurant.Experience.ContextRestExpEnabled, restaurant.Sustainability.UtensilsDefaultState,
		restaurant.PlannedSubsComplDate, restaurant.ActualSubsComplDate, restaurant.OnlineOrdering.OnlineOrderingEnabled,
		restaurant.OnlineOrdering.OnlineOrderingDotComSearchEnabled, restaurant.OnlineOrdering.OnlineOrderingCreditCardsAccepted,
		restaurant.OnlineOrdering.OnlineOrderingGiftCardsAccepted, restaurant.OnlineOrdering.OnlineOrderingBulkOrdersAccepted,
		restaurant.OnlineOrdering.OnlineOrderingTaxAssessed, restaurant.OnlineOrdering.RestaurantTerminalSiteID)

	restaurantID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error inserting restaurant: %v", err)
	}

	// Insert addresses
	for _, address := range restaurant.Addresses {
		_, err := db.Exec("INSERT INTO addresses (restaurant_id, address_type, address_line1, address_line2, locality, administrative_area, postal_code, sub_administrative_area, country_code, latitude, longitude, accuracy_determination) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			restaurantID, address.AddressType, address.AddressLine1, address.AddressLine2, address.Locality, address.AdministrativeArea, address.PostalCode, address.SubAdministrativeArea, address.CountryCode, address.Latitude, address.Longitude, address.AccuracyDetermination)
		if err != nil {
			return fmt.Errorf("error inserting address: %v", err)
		}
	}

	// Insert real hours
	for _, hours := range restaurant.RealHours {
		_, err := db.Exec("INSERT INTO real_hours (restaurant_id, day_of_week, open_date_time, close_date_time) VALUES (?, ?, ?, ?)",
			restaurantID, hours.DayOfWeek, hours.OpenDateTime, hours.CloseDateTime)
		if err != nil {
			return fmt.Errorf("error inserting real hours: %v", err)
		}
	}

	return nil
}

type items map[string]item
type itemTypes map[string]int
type itemCategories map[string]int
type itemNames map[string]int
type primaryFillingNames map[string]int
type contentGroups map[string]int

type optimizedItems struct {
	Items               items
	ItemTypes           itemTypes
	ItemCategories      itemCategories
	ItemNames           itemNames
	PrimaryFillingNames primaryFillingNames
	ContentGroups       contentGroups
}

func NewOptimizedItems() *optimizedItems {
	return &optimizedItems{
		Items:               make(items),
		ItemTypes:           make(itemTypes),
		ItemCategories:      make(itemCategories),
		ItemNames:           make(itemNames),
		PrimaryFillingNames: make(primaryFillingNames),
		ContentGroups:       make(contentGroups),
	}
}

func (o *optimizedItems) AddItem(item item) item {
	if _, ok := o.Items[item.ID]; !ok {
		o.Items[item.ID] = item
	}

	return o.Items[item.ID]
}

func (o *optimizedItems) AddItemType(itemType string) int {
	if _, ok := o.ItemTypes[itemType]; !ok {
		index := len(o.ItemTypes)
		o.ItemTypes[itemType] = index + 1
		return index
	}

	return o.ItemTypes[itemType]
}

func (o *optimizedItems) AddItemCategory(itemCategory string) int {
	if _, ok := o.ItemCategories[itemCategory]; !ok {
		index := len(o.ItemCategories)
		o.ItemCategories[itemCategory] = index + 1
		return index
	}

	return o.ItemCategories[itemCategory]
}

func (o *optimizedItems) AddItemName(itemName string) int {
	if _, ok := o.ItemNames[itemName]; !ok {
		index := len(o.ItemNames)
		o.ItemNames[itemName] = index + 1
		return index
	}

	return o.ItemNames[itemName]
}

func (o *optimizedItems) AddPrimaryFillingName(primaryFillingName string) int {
	if _, ok := o.PrimaryFillingNames[primaryFillingName]; !ok {
		index := len(o.PrimaryFillingNames)
		o.PrimaryFillingNames[primaryFillingName] = index + 1
		return index
	}

	return o.PrimaryFillingNames[primaryFillingName]
}

func (o *optimizedItems) AddContentGroup(contentGroup string) int {
	if _, ok := o.ContentGroups[contentGroup]; !ok {
		index := len(o.ContentGroups)
		o.ContentGroups[contentGroup] = index + 1
		return index
	}

	return o.ContentGroups[contentGroup]
}

type item struct {
	ID                 string
	Type               int
	Category           int
	Name               int
	PrimaryFillingName int
}

func optimizeItems(menus []*menu.Menu) *optimizedItems {
	oi := NewOptimizedItems()
	for _, menu := range menus {
		for _, e := range menu.Entrees {
			oi.AddItem(item{
				ID:                 e.ItemID,
				Type:               oi.AddItemType(e.ItemType),
				Category:           oi.AddItemCategory(e.ItemCategory),
				Name:               oi.AddItemName(e.ItemName),
				PrimaryFillingName: oi.AddPrimaryFillingName(e.PrimaryFillingName),
			})

			for _, c := range e.Contents {
				oi.AddItem(item{
					ID:   c.ItemID,
					Type: oi.AddItemType(c.ItemType),
					Name: oi.AddItemName(c.ItemName),
				})
			}

			for _, cg := range e.ContentGroups {
				oi.AddContentGroup(cg.ContentGroupName)
			}
		}
		for _, s := range menu.Sides {
			oi.AddItem(item{
				ID:       s.ItemID,
				Type:     oi.AddItemType(s.ItemType),
				Category: oi.AddItemCategory(s.ItemCategory),
				Name:     oi.AddItemName(s.ItemName),
			})
		}
		for _, d := range menu.Drinks {
			oi.AddItem(item{
				ID:       d.ItemID,
				Type:     oi.AddItemType(d.ItemType),
				Category: oi.AddItemCategory(d.ItemCategory),
				Name:     oi.AddItemName(d.ItemName),
			})
		}
		for _, nfi := range menu.NonFoodItems {
			oi.AddItem(item{
				ID:       nfi.ItemID,
				Type:     oi.AddItemType(nfi.ItemType),
				Category: oi.AddItemCategory(nfi.ItemCategory),
				Name:     oi.AddItemName(nfi.ItemName),
			})
		}
	}

	return oi
}
