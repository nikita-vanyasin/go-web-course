package common

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql" // importing driver
	log "github.com/sirupsen/logrus"
	"os"
)

const defaultServerPort = "8000"
const logDir = "logs/"

type EnvironmentSettings struct {
	SQLConnectionString string
	ContentFolderPath   string
	ServerPort          string
}

func OpenSQLConnection(connectionString string) *sql.DB {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func GetEnvSettings() *EnvironmentSettings {

	contentFolderPath := os.Getenv("CONTENT_FOLDER_PATH")
	if contentFolderPath == "" {
		log.Fatal("You need to specify content folder path!")
	}

	connectionString := os.Getenv("CONNECTION_STRING")
	if connectionString == "" {
		log.Fatal("You need to specify connection string!")
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = defaultServerPort
	}

	return &EnvironmentSettings{
		ContentFolderPath:   contentFolderPath,
		SQLConnectionString: connectionString,
		ServerPort:          serverPort,
	}
}

func SetupLogging(logFileName string) *os.File {
	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile(logDir+logFileName+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
	}
	return file
}
