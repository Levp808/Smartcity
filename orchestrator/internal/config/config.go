package config

type ConfigData struct {
	GRPCServerPort int    `env:"GRPC_SERVER_PORT" default:"50051" split_words:"true"`
	AIapi          string `env:"OPENROUTER_API_KEY"`
	KafkaBrokers   []string
}
