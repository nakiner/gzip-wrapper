package configs

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const ServiceName = "gzip-wrapper"

var options = []option{
	{"config", "string", "", "config file"},

	{"logger.level", "string", "info", "Level of logging. A string that correspond to the following levels: emerg, alert, crit, err, warning, notice, info, debug"},
	{"logger.time_format", "string", "2006-01-02T15:04:05.999999999", "Date format in logs"},

	{"s3.tenancy", "string", "", "s3 tenancy"},
	{"s3.region", "string", "", "s3 region"},
	{"s3.bucket_name", "string", "", "s3 bucket name"},
	{"s3.access_key_id", "string", "", "s3 access key ID"},
	{"s3.secret_access_key", "string", "", "s3 secret access key"},
}

type Config struct {
	Logger struct {
		Level      string
		TimeFormat string `mapstructure:"time_format"`
	}
	S3 S3
}

type S3 struct {
	Tenancy         string
	Region          string
	BucketName      string `mapstructure:"bucket_name"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
}

type option struct {
	name        string
	typing      string
	value       interface{}
	description string
}

// NewConfig returns and prints struct with config parameters
func NewConfig() *Config {
	return &Config{}
}

// read gets parameters from environment variables, flags or file.
func (c *Config) Read() error {
	viper.SetEnvPrefix(ServiceName)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	for _, o := range options {
		switch o.typing {
		case "string":
			pflag.String(o.name, o.value.(string), o.description)
		case "int":
			pflag.Int(o.name, o.value.(int), o.description)
		case "bool":
			pflag.Bool(o.name, o.value.(bool), o.description)
		case "float64":
			pflag.Float64(o.name, o.value.(float64), o.description)
		default:
			viper.SetDefault(o.name, o.value)
		}
	}

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	viper.BindPFlags(pflag.CommandLine)
	pflag.Parse()

	if fileName := viper.GetString("config"); fileName != "" {
		viper.SetConfigFile(fileName)
		viper.SetConfigType("toml")

		if err := viper.ReadInConfig(); err != nil {
			return errors.Wrap(err, "failed to read from file")
		}
	}

	if err := viper.Unmarshal(c); err != nil {
		return errors.Wrap(err, "failed to unmarshal")
	}
	return nil
}

func (c *Config) GenerateMdTable() error {
	var t string

	t += "Command line | Environment | Default |Description"
	t += fmt.Sprintln()
	t += "--- | --- | --- | ---"
	t += fmt.Sprintln()

	for _, o := range options {
		t += o.name
		t += " | "
		t += strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(ServiceName+"_"+o.name, ".", "_"), "-", "_"))
		t += " | "
		t += fmt.Sprintf("%v", o.value)
		t += " | "
		t += o.description
		t += fmt.Sprintln()
	}
	fmt.Fprintln(os.Stdout, t)
	return nil
}

func (c *Config) GenerateEnvironment() error {
	var t string

	for _, o := range options {
		t += strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(ServiceName+"_"+o.name, ".", "_"), "-", "_"))
		t += ": "
		t += fmt.Sprintf("%v", o.value)
		t += fmt.Sprintln()
	}
	fmt.Fprintln(os.Stdout, t)
	return nil
}

func (c *Config) GenerateFromTask() error {
	var t string

	t += "| Command line | Environment |"
	t += fmt.Sprintln()
	for _, o := range options {
		t += "| "
		t += strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(ServiceName+"_"+o.name, ".", "_"), "-", "_"))
		t += " | "
		t += fmt.Sprintf("%v", o.value)
		t += " |"
		t += fmt.Sprintln()
	}

	fmt.Fprintln(os.Stdout, t)
	return nil
}

func (c Config) Print() error {
	c.S3.SecretAccessKey = "******"

	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, string(b))
	return nil
}
