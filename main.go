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
	//-----

	// for name, action := range c.Actions {
	// 	fmt.Printf("Executing: %s\n", name)
	// 	execute(name, action, context)
	// }
}

// ### Vault pki bootstrap
//
// vault mount -path=rca pki
// vault mount-tune -default-lease-ttl=43800h -max-lease-ttl=87600h rca
// vault mount -path=ica pki
// vault mount-tune -default-lease-ttl=35040h -max-lease-ttl=70080h ica
//
// vault write rca/root/generate/internal common_name=rootca ttl=87600h key_bits=4096
// vault write -field=csr ica/intermediate/generate/internal common_name=interca ttl=70080h > inter.csr
// cat inter.csr | vault write -field=certificate rca/root/sign-intermediate csr=- use_csr_values=true > inter.cert
// rm inter.csr
// cat inter.cert | vault write ica/intermediate/set-signed certificate=-
// rm inter.cert
//
// ####### role
// vault write ica/roles/default ttl=8760h allow_any_name=true
