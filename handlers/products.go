package handlers

import (
	"github.com/kevin/working/data"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		p.getProducts(writer, request)
		return
	}

	if request.Method == http.MethodPost {
		p.addProduct(writer, request)
		return
	}

	if request.Method == http.MethodPut {

		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(request.URL.Path, -1)

		if len(g) != 1 {
			p.l.Println("Invalid URI more than one id")
			http.Error(writer, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			p.l.Println("Invalid URI more than one capture group")
			http.Error(writer, "Invalid URI", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, error := strconv.Atoi(idString)

		if error != nil {
			p.l.Println("Invalid URI unable to convert to number")
			http.Error(writer, "Invalid URI", http.StatusBadRequest)
			return
		}

		p.updateProducts(id, writer, request)
		return
	}

	writer.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(writer http.ResponseWriter, request *http.Request) {
	p.l.Println("Handle GET Products")
	lp := data.GetProducts()
	err := lp.ToJSON(writer)
	//d,err := json.Marshal(lp)
	//writer.Write(d)

	if err != nil {
		http.Error(writer, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(writer http.ResponseWriter, request *http.Request) {
	p.l.Println("Handle Post Products")

	prod := &data.Product{}
	err := prod.FromJSON(request.Body)
	if err != nil {
		http.Error(writer, "Uable to unmarshal json", http.StatusBadRequest)
	}

	data.AddProduct(prod)
}

func (p *Products) updateProducts(id int, writer http.ResponseWriter, request *http.Request) {
	p.l.Println("Handle PUT Products")

	prod := &data.Product{}
	err := prod.FromJSON(request.Body)
	if err != nil {
		http.Error(writer, "Uable to unmarshal json", http.StatusBadRequest)
	}

	err = data.UpdateProduct(id, prod)

	if err == data.ErrProductNotFound {
		http.Error(writer, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(writer, "Product not found", http.StatusInternalServerError)
		return
	}
}
