package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
)

func execute(name string, action Action, context map[string]interface{}) {
	err := action.Parse(context)
	if err != nil {
		fmt.Println(err)
	}
	output, err := action.Execute()
	if err != nil {
		fmt.Printf("Execution failed: %s\n", err)
	} else if output != nil {
		context[name] = output
	}
}

func main() {
	var confiFile string
	flag.StringVar(&confiFile, "config", "/etc/pipeline.d/config.json", "Configuration file")
	flag.Parse()

	c, err := GetConfig(confiFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Parsed configuration")
	context := make(map[string]interface{})
	context["env"] = GetEnvironment() // default environment output (TODO set as action?)

	//------TODO make dependency graph ( sort keys in the mean time )
	oa := make([]string, len(c.Actions))
	i := 0
	for k, _ := range c.Actions {
		oa[i] = k
		i++
	}
	sort.Strings(oa)

	for _, name := range oa {
		fmt.Printf("Executing: %s\n", name)
		execute(name, c.Actions[name], context)
	}
	fmt.Println("Finished")
	//-----

	// for name, action := range c.Actions {
	// 	fmt.Printf("Executing: %s\n", name)
	// 	execute(name, action, context)
	// }
}
