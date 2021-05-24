package hive

type StorageType string

const (
	TSV       StorageType = "tsv"
	OrcSnappy StorageType = "orc/snappy"
)