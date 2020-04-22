package swaggergen

import (
	"errors"
	"log"
	"os"

	"github.com/anz-bank/sysl/pkg/exporter"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func GenerateSwagger(appname string, module *sysl.Module, fs afero.Fs) error {
	app := getApplication(module, appname)
	if app == nil {
		return errors.New("no application found")
	}

	for k, v := range app.Endpoints {
		log.Printf("k %s", k)
		log.Printf("v %s", v)
	}

	swaggerExporter := exporter.MakeSwaggerExporter(app, logrus.New())
	err := swaggerExporter.GenerateSwagger()
	if err != nil {
		return err
	}
	byteOutput, err := swaggerExporter.SerializeOutput("json")
	if err != nil {
		return err
	}

	err = afero.WriteFile(fs, "/swagger.json", byteOutput, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func getApplication(module *sysl.Module, appname string) *sysl.Application {
	return module.Apps[appname]
}
