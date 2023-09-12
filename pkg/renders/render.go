package renders

import (
	"bytes"
	"github.com/satish218/hotel_bookings/pkg/config"
	"github.com/satish218/hotel_bookings/pkg/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var app *config.AppConfig

// NewTemplates sets the config for the new template
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreteTemplateCache()
	}
	//get the template cache from app config

	//create a template cache
	//tc, err := CreteTemplateCache()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("could not get template from template cache")
	}

	buf := new(bytes.Buffer)
	td = AddDefaultData(td)
	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}
	//render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}

	//parsedTemplate, _ := template.ParseFiles("./templates/"+tmpl, "./templates/base.layout.gohtml")
	//err := parsedTemplate.Execute(w, nil)
	//if err != nil {
	//	fmt.Println("error parsing the file")
	//	return
	//}
}

//create template cache

func CreteTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	//get all the files named *.page.gohtml from ./templates folder
	pages, err := filepath.Glob("./templates/*.page.gohtml")
	if err != nil {
		return myCache, err
	}
	//range through all files ending with *page.gohtml

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob("./templates/*layout.gohtml")
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.gohtml")
		}
		myCache[name] = ts
	}
	return myCache, nil
}

//var tc = make(map[string]*template.Template)
//
//func RenderTemplate(w http.ResponseWriter, t string) {
//	var tmpl *template.Template
//	var err error
//	_, inMap := tc[t]
//	if !inMap {
//		//Need to create template
//		fmt.Println("creating cached template")
//		err = createTemplateCache(t)
//		if err != nil {
//			fmt.Println(err)
//		}
//
//	} else {
//		//use cached template
//		fmt.Println("using cached template")
//	}
//	tmpl = tc[t]
//	err = tmpl.Execute(w, nil)
//	if err != nil {
//		fmt.Println(err)
//	}
//}
//
//func createTemplateCache(t string) error {
//	templates := []string{
//		fmt.Sprintf("./templates/%s", t),
//		"./templates/base.layout.gohtml",
//	}
//	tmpl, err := template.ParseFiles(templates...)
//	if err != nil {
//		return err
//	}
//	tc[t] = tmpl
//	return nil
//}
