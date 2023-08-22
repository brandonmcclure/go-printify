package go_printify

import (
	"fmt"
	"io"
	"net/http"
)

const (
	productsPath       = "shops/%d/products.json"
	productPath        = "shops/%d/products/%s.json"
	publishProductPath = "shops/%d/products/%s/publish.json"
	publishSuccessPath = "shops/%d/products/%s/publishing_succeeded.json"
	publishFailedPath  = "shops/%d/products/%s/publishing_failed.json"
	unpublishPath      = "shops/%d/products/%s/unpublish.json"
)

type ProductsResponse struct {
	CurrentPage  int       `json:"current_page"`
	Data         []Product `json:"data"`
	FirstPageUrl string    `json:"first_page_url"`
	LastPageUrl  string    `json:"last_page_url"`
	NextPageUrl  string    `json:"next_page_url"`
	From         int       `json:"from"`
	LastPage     int       `json:"last_page"`
	Path         string    `json:"path"`
	PerPage      int       `json:"per_page"`
	PrevPageUrl  string    `json:"prev_page_url"`
	To           int       `json:"to"`
	Total        int       `json:"total"`
}

type Product struct {
	Id                     string                   `json:"id,omitempty"`
	Title                  string                   `json:"title"`
	Description            string                   `json:"description"`
	Tags                   []string                 `json:"tags"`
	Options                []map[string]interface{} `json:"options"`
	Variants               []ProductVariant         `json:"variants"`
	Images                 []ProductMockUpImage     `json:"images"`
	CreatedAt              string                   `json:"created_at,omitempty"`
	UpdatedAt              string                   `json:"updated_at,omitempty"`
	Visible                bool                     `json:"visible"`
	BlueprintId            int                      `json:"blueprint_id"`
	PrintProviderId        int                      `json:"print_provider_id"`
	UserId                 int                      `json:"user_id"`
	ShopId                 int                      `json:"shop_id"`
	PrintAreas             []PrintArea              `json:"print_areas"`
	PrintDetails           []PrintDetails           `json:"print_details"`
	External               *External                `json:"external,omitempty"`
	IsLocked               bool                     `json:"is_locked"`
	SalesChannelProperties []string                 `json:"sales_channel_properties,omitempty"`
}

type ProductCreation struct {
	Title           string                   `json:"title"`
	Description     string                   `json:"description"`
	Variants        []ProductCreationVariant `json:"variants"`
	BlueprintId     int                      `json:"blueprint_id"`
	PrintProviderId int                      `json:"print_provider_id"`
	PrintAreas      []PrintArea              `json:"print_areas"`
}

type ProductCreationVariant struct {
	Id        int     `json:"id"`
	Price     float32 `json:"price"`
	IsEnabled bool    `json:"is_enabled"`
	IsDefault bool    `json:"is_default"`
}

type ProductVariant struct {
	Id          int     `json:"id"`
	Sku         string  `json:"sku"`
	Price       float32 `json:"price"`
	Cost        float32 `json:"cost"`
	Title       string  `json:"title"`
	Grams       int     `json:"grams"`
	IsEnabled   bool    `json:"is_enabled"`
	InStock     bool    `json:"in_stock"` // Deprecated
	IsDefault   bool    `json:"is_default"`
	IsAvailable bool    `json:"is_available"`
	Options     []int   `json:"options"`
}

type ProductMockUpImage struct {
	Src        string `json:"src"`
	VariantIds []int  `json:"variant_ids"`
	Position   string `json:"position"`
	IsDefault  bool   `json:"is_default"`
	IsPublish  bool   `json:"is_selected_for_publishing"`
}

type ProductPlaceholder struct {
	Position string         `json:"position"`
	Images   []ProductImage `json:"images"`
}

type ProductImage struct {
	Id     string  `json:"id"`
	Name   string  `json:"name"`
	Type   string  `json:"type"`
	Height int     `json:"height"`
	Width  int     `json:"width"`
	X      float32 `json:"x"`
	Y      float32 `json:"y"`
	Scale  float32 `json:"scale"`
	Angle  int     `json:"angle"`
}

type PrintArea struct {
	VariantIds   []int                `json:"variant_ids"`
	Placeholders []ProductPlaceholder `json:"placeholders"`
}

type PrintDetails struct {
	PrintOnSide string `json:"print_on_side"`
}

type PublishingProperties struct {
	Images      bool `json:"images"`
	Variants    bool `json:"variants"`
	Title       bool `json:"title"`
	Description bool `json:"description"`
	Tags        bool `json:"tags"`
}

type External struct {
	Id               string `json:"id"`
	Handle           string `json:"handle"`
	ShippingTemplate string `json:"shipping_template_id"`
}

/*
Retrieve a list of products
*/
func (c *Client) GetAllProducts(shopId int) ([]Product, error) {

	var allProducts []Product
	page := 1
	for {
		productResults, err := c.GetProducts(shopId, &page)
		if err != nil {
			fmt.Println("Received error from getProducts")
			return nil, err
		}

		allProducts = append(allProducts, productResults.Data...)

		if productResults.NextPageUrl == "" {
			break
		}
		page++
	}

	return allProducts, nil
}

func (c *Client) GetProducts(shopId int, page *int) (*ProductsResponse, error) {
	path := fmt.Sprintf(productsPath, shopId)
	query := ""
	if page != nil {
		query = fmt.Sprintf("page=%d", *page)
	}
	req, err := c.newRequest(http.MethodGet, path, query, nil)
	if err != nil {
		return nil, err
	}
	products := ProductsResponse{}
	_, err = c.do(req, &products)
	return &products, err
}

/*
Retrieve a product
*/
func (c *Client) GetProduct(shopId int, productId string) (*Product, error) {
	path := fmt.Sprintf(productPath, shopId, productId)
	req, err := c.newRequest(http.MethodGet, path, "", nil)
	if err != nil {
		return nil, err
	}
	product := &Product{}
	_, err = c.do(req, product)
	return product, err
}

/*
Create a new product
*/
func (c *Client) CreateProduct(shopId int, product ProductCreation) (*Product, error) {
	path := fmt.Sprintf(productsPath, shopId)
	req, err := c.newRequest(http.MethodPost, path, "", product)
	if err != nil {
		return nil, err
	}
	respProd := &Product{}
	_, err = c.do(req, respProd)
	return respProd, err
}

/*
Update a product
*/
func (c *Client) UpdateProduct(shopId int, product *Product) (*Product, error) {
	path := fmt.Sprintf(productPath, shopId, product.Id)
	req, err := c.newRequest(http.MethodPut, path, "", product)
	if err != nil {
		return nil, err
	}
	updatedProduct := &Product{}
	resp, err := c.do(req, updatedProduct)
	if err != nil {
		bb, _ := io.ReadAll(resp.Body)
		fmt.Println("RESP", string(bb))
	}
	return updatedProduct, err
}

/*
Delete a product
*/
func (c *Client) DeleteProduct(shopId int, productId string) error {
	path := fmt.Sprintf(productPath, shopId, productId)
	req, err := c.newRequest(http.MethodDelete, path, "", nil)
	if err != nil {
		return err
	}
	_, err = c.do(req, nil)
	return err
}

/*
Publish a product
*/
func (c *Client) PublishProduct(shopId int, productId string, publishProperties PublishingProperties) error {
	path := fmt.Sprintf(publishProductPath, shopId, productId)
	req, err := c.newRequest(http.MethodPost, path, "", publishProperties)
	if err != nil {
		return err
	}
	_, err = c.do(req, nil)
	return err
}

/*
Set product publish status to succeeded
*/
func (c *Client) SetProductPublishSuccess(shopId int, productId string, external External) error {
	path := fmt.Sprintf(publishSuccessPath, shopId, productId)
	req, err := c.newRequest(http.MethodPost, path, "", external)
	if err != nil {
		return err
	}
	_, err = c.do(req, nil)
	return err
}

/*
Set product publish status to failed
*/
func (c *Client) SetProductPublishFailre(shopId int, productId string, reason string) error {
	path := fmt.Sprintf(publishFailedPath, shopId, productId)
	req, err := c.newRequest(http.MethodPost, path, "", map[string]string{"reason": reason})
	if err != nil {
		return err
	}
	_, err = c.do(req, nil)
	return err
}

/*
Notify that a product has been unpublished
*/
func (c *Client) UnPublish(shopId int, productId string) error {
	path := fmt.Sprintf(unpublishPath, shopId, productId)
	req, err := c.newRequest(http.MethodPost, path, "", nil)
	if err != nil {
		return err
	}
	_, err = c.do(req, nil)
	return err
}
