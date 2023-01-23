package cmd

import (
	"fmt"
	"jpodivin/gomata/computer"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var random bool
var networked bool
var raddress string
var laddress string
var rightHanded bool

func init() {
	rootCmd.PersistentFlags().BoolVar(&random, "random", false, "Random initial state")
	rootCmd.PersistentFlags().BoolVar(&networked, "networked", false, "Run in networked mode")
	rootCmd.PersistentFlags().StringVar(&raddress, "raddress", ":1993", "Respond to communication on this address.")
	rootCmd.PersistentFlags().StringVar(&laddress, "laddress", ":1952", "Listen to communication on this address")
	rootCmd.PersistentFlags().BoolVar(&rightHanded, "right", true, "Set right site to networked value")
}

func PrintState(state []int8) string {
	strState := fmt.Sprintf("%v", state)
	return strState
}

func StrToInt(val string) int {
	c_val, err := strconv.Atoi(val)
	if err != nil {
		log.Fatal(fmt.Errorf("%v", err))
	}

	return c_val
}

var rootCmd = &cobra.Command{
	Use:   "gomata",
	Short: "Networked cellular automaton",
	Long:  "ditto",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var runtimeError error
		rule := StrToInt(args[0])
		time := StrToInt(args[1])
		var gameWorld computer.World

		if networked {
			gameWorld, runtimeError = computer.InitWorld(80, 0.2, random, rightHanded, raddress, laddress)

		} else {
			gameWorld, runtimeError = computer.InitWorld(80, 0.2, random, rightHanded)
		}

		ruleCode := computer.ComputeRule(rule)

		for i := 0; i < time; i++ {

			if runtimeError != nil {
				log.Fatal(fmt.Errorf("%v", runtimeError))
			}
			fmt.Println(PrintState(gameWorld.CurrentState))
			if networked {
				runtimeError = computer.ComputeState(gameWorld, ruleCode, laddress)
			} else {
				runtimeError = computer.ComputeState(gameWorld, ruleCode)
			}

		}
	},
}

func Execute() {
	log.SetFlags(1)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
