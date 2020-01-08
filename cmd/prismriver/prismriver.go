package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ttpcodes/prismriver/assets"
	"gitlab.com/ttpcodes/prismriver/internal/app/constants"
	"gitlab.com/ttpcodes/prismriver/internal/app/server"
	"io"
	"os"
	"path"
)

func main() {
	// Set up configuration framework.
	viper.SetEnvPrefix("PRISMRIVER")
	viper.AutomaticEnv()

	viper.SetDefault(constants.DATA, "/var/lib/prismriver")
	viper.SetDefault(constants.DBHOST, "localhost")
	viper.SetDefault(constants.DBNAME, "prismriver")
	viper.SetDefault(constants.DBPASSWORD, "prismriver")
	viper.SetDefault(constants.DBPORT, "5432")
	viper.SetDefault(constants.DBUSER, "prismriver")
	viper.SetDefault(constants.VERBOSITY, "info")

	envVars := []string{
		constants.DBHOST,
		constants.DBNAME,
		constants.DBPASSWORD,
		constants.DBPORT,
		constants.DBUSER,
		constants.VERBOSITY,
	}

	for _, env := range envVars {
		if err := viper.BindEnv(env); err != nil {
			logrus.Warnf("error binding to variable %v: %v", env, err)
		}
	}

	verbosity := viper.GetString(constants.VERBOSITY)
	level, err := logrus.ParseLevel(verbosity)
	if err != nil {
		logrus.Errorf("Error reading verbosity level in configuration: %v", err)
	}
	logrus.SetLevel(level)
	dataDir := viper.GetString(constants.DATA)
	if err := os.MkdirAll(dataDir, os.ModeDir); err != nil {
		logrus.Fatalf("error creating data directory: %v", err)
	}
	if err := os.MkdirAll(dataDir+"/internal", os.ModeDir); err != nil {
		logrus.Fatalf("error creating internal directory: %v", err)
	}

	beQuiet, err := assets.HTTP.Open("bequiet.opus")
	if err != nil {
		logrus.Fatalf("Error reading bequiet.opus in internal filesystem (is this binary corrupted?): %v", err)
	}
	beQuietPath := path.Join(dataDir, "internal", "bequiet.opus")
	beQuietFile, err := os.Create(beQuietPath)
	if err != nil {
		logrus.Fatalf("Error creating application files: %v", err)
	}
	if _, err := io.Copy(beQuietFile, beQuiet); err != nil {
		logrus.Fatalf("error copying bequiet.opus: %v", err)
	}
	if err := beQuietFile.Close(); err != nil {
		logrus.Warnf("error closing reader on bequiet.opus: %v", err)
	}

	server.CreateRouter()
}
