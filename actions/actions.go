package actions

var EXECUTORS Executors

type Executors struct {
	all map[string]Executor
}

func (e Executors) GetByName(name string) Executor {
	return e.all[name]
}

func (e Executors) GetNames() []string {
	keys := make([]string, 0, len(e.all))
	for k := range e.all {
		keys = append(keys, k)
	}
	return keys
}

func Init(jsonConfigFile string) error {

	err, cfg := loadConfigurationFromJson(jsonConfigFile)
	if err == nil {
		EXECUTORS.all = make(map[string]Executor)
		buildAllExecutors(cfg.Actions, cfg.Servers,
			func(name string, exe Executor) {
				EXECUTORS.all[name] = exe
			})
		return nil
	}
	return err
}
