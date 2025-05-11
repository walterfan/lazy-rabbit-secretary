package book

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Structs for parsing Google Books API response
type GoogleBooksResponse struct {
	Items []struct {
		VolumeInfo struct {
			Title               string   `json:"title"`
			Authors             []string `json:"authors"`
			IndustryIdentifiers []struct {
				Type       string `json:"type"`
				Identifier string `json:"identifier"`
			} `json:"industryIdentifiers"`
			Publisher     string `json:"publisher"`
			PublishedDate string `json:"publishedDate"`
			ListPrice     struct {
				Amount       float64 `json:"amount"`
				CurrencyCode string  `json:"currencyCode"`
			} `json:"listPrice,omitempty"`
		} `json:"volumeInfo"`
		SaleInfo struct {
			ListPrice struct {
				Amount       float64 `json:"amount"`
				CurrencyCode string  `json:"currencyCode"`
			} `json:"listPrice,omitempty"`
		} `json:"saleInfo"`
	} `json:"items"`
}

// BookInfo struct to hold the book details
type BookInfo struct {
	Title         string
	Authors       string
	ISBN          string
	Publisher     string
	PublishedDate string
	Price         string
}

func searchBook(title string, author string) (BookInfo, error) {
	// Construct the API URL
	query := url.QueryEscape("intitle:" + title + "+inauthor:" + author)
	apiURL := "https://www.googleapis.com/books/v1/volumes?q=" + query

	// Make the HTTP GET request
	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("Error making API request:", err)
		return BookInfo{}, err
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return BookInfo{}, err
	}

	// Parse JSON
	var data GoogleBooksResponse
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return BookInfo{}, err
	}

	// Extract book information
	if len(data.Items) == 0 {
		return BookInfo{}, nil
	}

	item := data.Items[0]
	var isbn string
	for _, id := range item.VolumeInfo.IndustryIdentifiers {
		if id.Type == "ISBN_13" {
			isbn = id.Identifier
			break
		}
	}

	var price string
	if item.SaleInfo.ListPrice.Amount > 0 {
		price = fmt.Sprintf("%.2f %s", item.SaleInfo.ListPrice.Amount, item.SaleInfo.ListPrice.CurrencyCode)
	} else if item.VolumeInfo.ListPrice.Amount > 0 {
		price = fmt.Sprintf("%.2f %s", item.VolumeInfo.ListPrice.Amount, item.VolumeInfo.ListPrice.CurrencyCode)
	}

	return BookInfo{
		Title:         item.VolumeInfo.Title,
		Authors:       strings.Join(item.VolumeInfo.Authors, ", "),
		ISBN:          isbn,
		Publisher:     item.VolumeInfo.Publisher,
		PublishedDate: item.VolumeInfo.PublishedDate,
		Price:         price,
	}, nil
}

func main() {
	// Input: title and author
	title := "微服务之道: 度量驱动开发"
	author := "范亚敏"

	// Search for the book
	bookInfo, err := searchBook(title, author)

	if err != nil {
		fmt.Println("Error searching for book:", err)
		return
	}
	// Print out the book information
	if bookInfo.ISBN == "" {
		fmt.Println("No results found.")
		return
	}

	fmt.Println("Title:", bookInfo.Title)
	fmt.Println("Authors:", bookInfo.Authors)
	fmt.Println("ISBN:", bookInfo.ISBN)
	fmt.Println("Publisher:", bookInfo.Publisher)
	fmt.Println("Published Date:", bookInfo.PublishedDate)
	fmt.Println("Price:", bookInfo.Price)
}
