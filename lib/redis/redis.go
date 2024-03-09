package lib


type RedisLib struct {
	ConfPath string 
}

var RedisLibInstance *RedisLib


func (rL *RedisLib) InitConf () error {
	return nil
}

func (rL *RedisLib) SetPath(fileName string, ConfEnvPath string)  {
	rL.ConfPath = ConfEnvPath + "/" + fileName + ".toml"
}