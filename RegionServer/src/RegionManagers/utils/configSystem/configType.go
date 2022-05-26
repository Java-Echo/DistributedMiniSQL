package config

//创建一个结构体
type Config struct {
	Etcd_ip                      string
	Etcd_port                    string
	Etcd_region_register_catalog string
	Etcd_master_address          string
	Rpc_R2R_port                 string
	Rpc_M2R_port                 string
	Rpc_R2C_port                 string
	Etcd_table_catalog           string
	Minisql_table_store          string
}
