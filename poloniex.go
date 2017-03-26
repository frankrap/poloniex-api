// Poloniex commons
package poloniex

import (
	"encoding/json"
	"fmt"
	"log"
)

func PrettyPrintJson(msg interface{}) {

	jsonstr, err := json.MarshalIndent(msg, "", "  ")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", string(jsonstr))
}
