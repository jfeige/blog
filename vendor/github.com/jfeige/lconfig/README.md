# lconfig

安装:

	go get github.com/jfeige/lconfig
	
使用:

	配置文件格式参考myconfig.ini

	加载配置文件:
                config,err := NewConfig("myconfig.ini")
                //string
                v := config.String("host")
                or
                v := config.String("default:host")
                or
                v := config.String("redis::host")
                //int
                v,err := config.Int("port")
                or
                v,err := config.Int("redis::port")
                //[]string
                slave_addr,err := config.Int("mysql::slave_addr")
		or
		v,err := config.Sections("province")
                
API列表:

	String(key string) string
	
	Strings(key string) ([]string,error)
	
	Int(key string)(int,error)
	
	Int64(key string)(int64,error)
	
	Bool(key string)(bool,error)
	
	Float64(key string)(float64,error)
	
	Sections(key string)(map[string]string,error)
