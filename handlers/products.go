package handlers

import (
	"github.com/gorilla/mux"
	"github.com/kevin/working/data"
	"log"
	"net/http"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(writer http.ResponseWriter, request *http.Request) {
	p.l.Println("Handle GET Products")
	lp := data.GetProducts()
	err := lp.ToJSON(writer)
	//d,err := json.Marshal(lp)
	//writer.Write(d)

	if err != nil {
		http.Error(writer, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(writer http.ResponseWriter, request *http.Request) {
	p.l.Println("Handle Post Products")

	prod := &data.Product{}
	err := prod.FromJSON(request.Body)
	if err != nil {
		http.Error(writer, "Unable to unmarshal json", http.StatusBadRequest)
	}

	data.AddProduct(prod)
}

func (p *Products) UpdateProducts(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(writer, "Unable to get id", http.StatusBadRequest)
	}

	p.l.Println("Handle PUT Products")

	prod := &data.Product{}
	err = prod.FromJSON(request.Body)
	if err != nil {
		http.Error(writer, "Unable to unmarshal json", http.StatusBadRequest)
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
