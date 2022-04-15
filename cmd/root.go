package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const EnvPrefix = "S3TEST"

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "s3test",
	Short: "tests failing s3 call",
	Long:  `Executes problematic s3 calls for debugging purposes.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.s3test.yaml)")
	rootCmd.PersistentFlags().StringP("key", "k", "", "Access key (required)")
	rootCmd.PersistentFlags().StringP("secret", "s", "", "Access secret (required)")
	rootCmd.PersistentFlags().StringP("endpoint", "e", "", "S3 endpoint (required)")
	rootCmd.PersistentFlags().StringP("bucket", "b", "", "Bucket (required)")

	_ = viper.BindPFlag("key", rootCmd.PersistentFlags().Lookup("key"))
	_ = viper.BindPFlag("secret", rootCmd.PersistentFlags().Lookup("secret"))
	_ = viper.BindPFlag("endpoint", rootCmd.PersistentFlags().Lookup("endpoint"))
	_ = viper.BindPFlag("bucket", rootCmd.PersistentFlags().Lookup("bucket"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".s3test.yaml")
	}

	viper.SetEnvPrefix(EnvPrefix)
	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}
}
