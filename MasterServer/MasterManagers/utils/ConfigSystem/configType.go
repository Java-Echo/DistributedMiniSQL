package config

//创建一个结构体
type Config struct {
	Etcd_region_register_catalog string
	Etcd_table_catalog           string
	Rpc_m2r_port                 string
	Rpc_m2c_port                 string
}
