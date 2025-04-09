package main

import (
	util "go-dev-template/mymodule/util"
)

func main() {
	// env_name := "ANY_TOKEN"
	// TOKEN, ok := os.LookupEnv(env_name)
	// if !ok {
	// 	fmt.Printf("%s is not set", env_name)
	// }
	// results, err := util.AnyFunction(TOKEN)
	// if err != nil {
	// 	util.OutLog(err)
	// 	panic(err)
	// }
	util.OutLog("main: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	// util.OutLog(results)
	// j, err := json.Marshal(results)
	// b := bytes.NewBuffer([]byte(j))
	// util.OutLog(b)

	// util.OutLog(results)
	l := util.NewBuiltinLogger("FILE")
	l.Debug("this is debug message")
	l.Info("this is info message")
	l.Warning("this is warning message")
	l.Error("this is error message")
	l.Fatal("this is fatal message")
}
