package config

//创建一个结构体
type Config struct {
	Etcd_ip                      string
	Etcd_port                    string
	Etcd_region_register_catalog string
	Etcd_master_address          string
}
