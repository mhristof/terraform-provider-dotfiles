package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCurl() *schema.Resource {
	return &schema.Resource{
		Create: resourceCurlCreate,
		Read:   resourceCurlRead,
		Update: resourceCurlUpdate,
		Delete: resourceCurlDelete,

		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"dest": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"extract": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceCurlCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg)

	log.Println(fmt.Sprintf("[DEBUG] config: %+v", config))

	url, dest := getDest(d, meta)

	d.SetId(dest)

	Wget(url, dest)

	extract := d.Get("extract").(string)

	if extract != "" {
		var err error
		switch {
		case strings.Contains(dest, ".zip"):
			_, err = ExtractZipToFile(dest, filepath.Join(config.root, extract))
		case strings.Contains(dest, ".tar.gz"):
			_, err = ExtractTarGzToFile(dest, filepath.Join(config.root, extract))
		default:
			err = errors.New("unsupported archive type")
		}
		os.Remove(dest)
		return err
	}

	log.Println(fmt.Sprintf("[DEBUG] wget %s %s", url, dest))
	return resourceCurlRead(d, meta)
}

func ExtractTarGzToFile(src string, file string) (string, error) {
	f, err := os.Open(src)
	if err != nil {
		return "", err
	}

	uncompressedStream, err := gzip.NewReader(f)
	if err != nil {
		return "", err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return "", err
		}

		switch header.Typeflag {
		case tar.TypeReg:
			log.Println(fmt.Sprintf("header.Name: %+v", header.Name))
			log.Println(fmt.Sprintf("file: %+v", file))

			if path.Base(header.Name) != path.Base(file) {
				continue
			}

			outFile, err := os.Create(file)
			if err != nil {
				return "", err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return "", err
			}
			outFile.Close()
			return file, nil
		default:
			errors.New(fmt.Sprintf(
				"ExtractTarGz: uknown type: %+v in %s",
				header.Typeflag,
				header.Name))
		}

	}
	return "", errors.New("File not found")
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func ExtractZipToFile(src string, file string) (string, error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		return "", err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name != path.Base(file) {
			continue
		}

		outFile, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return "", err
		}

		rc, err := f.Open()
		if err != nil {
			return "", err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		return "", nil

	}

	return "", errors.New("File not found")
}

func Wget(URL, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}
	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func getDest(d *schema.ResourceData, meta interface{}) (string, string) {
	config := meta.(*cfg)

	url := d.Get("url").(string)
	dest := d.Get("dest").(string)
	if dest == "" {
		dest = path.Base(url)
	}

	return url, filepath.Join(config.root, dest)
}

func resourceCurlRead(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] reading resource event")
	return nil
}

func resourceCurlUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceCurlCreate(d, meta)
}

func resourceCurlDelete(d *schema.ResourceData, meta interface{}) error {
	_, dest := getDest(d, meta)
	os.Remove(dest)
	return nil
}
