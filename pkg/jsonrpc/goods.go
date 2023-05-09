package jsonrpc

import (
	"fmt"
	"log"

	lamodatest "github.com/NikiTesla/lamoda_test"
	"github.com/NikiTesla/lamoda_test/pkg/database"
)

// Goods is struct to interact with goods at different warehouses, has environment with bd as a field
type Goods struct {
	db database.DataBase
}

// Create gets list of goods as []Good, ask database to create them
// put in reply successfully created goods
func (g *Goods) Create(args []lamodatest.Good, reply *Reply) error {
	created := make([]lamodatest.Good, 0, len(args))
	for _, good := range args {
		log.Printf("Creating good: %v", good)

		if err := g.db.CreateGood(good); err != nil {
			log.Print("Error while creating goods ", err.Error())
			return err
		}
		created = append(created, good)
	}
	*reply = Reply{fmt.Sprintf("Created goods %v: ", created)}

	return nil
}

// Add gets list of maps with good_code, warehouse_id and amount of goods to be added
// put in reply successfully added goods
func (g *Goods) Add(args []map[string]int, reply *Reply) error {
	log.Printf("Adding goods: %v", args)

	added := make([]map[string]int, 0, len(args))
	for _, arg := range args {
		good_code, warehouse_id, amount, err := parseGoods(arg)
		if err != nil {
			log.Print("Error occured while parsing args, error: ", err.Error())
			return err
		}

		if err := g.db.AddGood(good_code, warehouse_id, amount); err != nil {
			log.Print("Cannot add good, error: ", err.Error())
			continue
		}
		added = append(added, arg)
	}
	*reply = Reply{fmt.Sprintf("Added goods %v", added)}

	return nil
}

// Reserve gets list of maps with good_code, warehouse_id and amount of goods to be reserved
// put in reply successfully reserved goods
func (g *Goods) Reserve(args []map[string]int, reply *Reply) error {
	log.Print("Reserving: ", args)
	reserved := make([]map[string]int, 0, len(args))
	for _, arg := range args {
		good_code, warehouse_id, amount, err := parseGoods(arg)
		if err != nil {
			log.Print("Error while parsing args ", err.Error())
			return err
		}

		if err = g.db.ReserveGood(good_code, warehouse_id, amount); err != nil {
			log.Printf("error occured while reserving good with code %d, error: %s", good_code, err.Error())
			continue
		}
		reserved = append(reserved, arg)
	}
	*reply = Reply{fmt.Sprintf("Reserved goods: %v", reserved)}

	log.Print("reserved")

	return nil
}

// CancelReservation gets lsit of maps with good_code, warehouse_id and amount of goods to cancel reservation
// put in reply successfully cancelled reservations of goods
func (g *Goods) CancelReservation(args []map[string]int, reply *Reply) error {
	log.Print("Canceling reservations: ", args)
	cancelled := make([]map[string]int, 0, len(args))
	for _, arg := range args {
		good_code, warehouse_id, amount, err := parseGoods(arg)
		if err != nil {
			log.Print("Error while parsing args ", err.Error())
			return err
		}

		if err := g.db.CancelGoodReservation(good_code, warehouse_id, amount); err != nil {
			log.Printf("error occured while cancelling reservation of good with code %d, error: %s", good_code, err.Error())
			continue
		}
		cancelled = append(cancelled, arg)
	}

	*reply = Reply{fmt.Sprintf("Cancel reservations of goods: %v", cancelled)}

	log.Print("cancelled")

	return nil
}

// parseGoods is internal function to parse args where good_code, warehouse_id and amount all required
// gets map with these fields, returns them in presented above order and error
func parseGoods(arg map[string]int) (int, int, int, error) {
	good_code, ok := arg["good_code"]
	if !ok {
		return 0, 0, 0, fmt.Errorf("incorrect request, good_code not presented")
	}
	warehouse_id, ok := arg["warehouse_id"]
	if !ok {
		return 0, 0, 0, fmt.Errorf("incorrect request, warehouse_id not presented")
	}
	amount, ok := arg["amount"]
	if !ok {
		return 0, 0, 0, fmt.Errorf("incorrect request, amount not presented")
	}
	return good_code, warehouse_id, amount, nil
}