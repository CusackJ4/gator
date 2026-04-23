package main

import (
	"fmt"
	"os"

	"database/sql"

	"github.com/CusackJ4/gator/internal/config"
	"github.com/CusackJ4/gator/internal/database"

	_ "github.com/lib/pq"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}

	// db stuff
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		fmt.Printf("An error! %v\n", err)
	}
	dbQueries := database.New(db)

	s := state{cfg: &cfg, db: dbQueries}
	cmdsInst := commands{cmds: map[string]func(*state, command) error{}}

	cmdsInst.register("login", handlerLogin)
	cmdsInst.register("register", handlerRegister)
	cmdsInst.register("reset", reset)
	cmdsInst.register("users", getNames)
	cmdsInst.register("addfeed", middlewareFunc(addFeed))
	cmdsInst.register("agg", agg)
	cmdsInst.register("feeds", viewFeeds)
	cmdsInst.register("follow", middlewareFunc(follow))
	cmdsInst.register("following", middlewareFunc(following))
	cmdsInst.register("unfollow", middlewareFunc(unfollow))
	cmdsInst.register("browse", middlewareFunc(browse))

	input := os.Args
	if len(input) < 2 {
		fmt.Println("Insufficient args")
		os.Exit(1)
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]
	cmdInst := command{name: cmdName,
		args: cmdArgs}
	err = cmdsInst.run(&s, cmdInst)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

//ex: go run . register paul
