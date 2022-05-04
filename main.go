package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	const httpPort = "8000"

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Diag Service")
	})

	e.POST("/upload", upload)

	e.Logger.Fatal(e.Start(":" + httpPort))

}

func upload(c echo.Context) error {

	uploadDir := "diags"

	// Source
	file, err := c.FormFile("diag")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}

	//Vaidating format is *.tgz
	if !strings.HasSuffix(file.Filename, ".tgz") {
		return c.HTML(http.StatusUnsupportedMediaType, fmt.Sprintf("<p>Invalid file format for %s. This service only accepts *.tgz files.</p>", file.Filename))
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(filepath.Join(uploadDir, file.Filename))
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, fmt.Sprintf("<p>File %s uploaded.</p>", file.Filename))
}
