package main

import "prayersNetwork/server/routers"

func main() {

	router := routers.RegisterRouters()

	router.Run(":8085")
}
