package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/nabishec/referal_links/internal/middleware/logger"
	"github.com/nabishec/referal_links/internal/storage/postgesql/db"
)

func main() {
	//TODO: init logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	debug := flag.Bool("debug", false, "set log level to debug")
	easyReading := flag.Bool("read", false, "set console writer")
	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	//for easy reading
	if *easyReading {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	//TODO: init config
	err := godotenv.Load("./configs/configuration.env")
	if err != nil {
		log.Error().Err(err).Msg("don't found configuration")
		os.Exit(1)
	}

	//TODO: init storage postgresql
	log.Info().Msg("Init storage")
	storage, err := db.NewDatabase()
	if err != nil {
		log.Error().AnErr(ErrReader(err)).Msg("Failed init storage")
		os.Exit(1)
	}
	log.Info().Msg("Storage init successful")
	_ = storage

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(logger.New())
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	_, _ = router, storage
	//TODO: run-server
	log.Warn().Msg("Program ended")
}

func ErrReader(err error) (f string, e error) {
	str := strings.Split(err.Error(), ":")
	f = str[0]
	e = fmt.Errorf(str[1])
	return f, e
}
