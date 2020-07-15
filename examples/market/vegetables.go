package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
)

// FruitIdentityParams represents the parameters that
// are required to identity a unique fruit in the market.
type VegetableIdentityParams struct {
	Name   string `path:"name"`
	APIKey string `header:"X-Api-Key" validate:"required"`
}

// ListFruitsParams represents the parameters that can
// be used to filter the fruit's market listing.
type ListVegetablesParams struct {
	Origin   *string  `query:"origin" description:"filter by fruit origin"`
	PriceMin *float64 `query:"price_min" description:"filter by minimum inclusive price" validate:"omitempty,min=1"`
	PriceMax *float64 `query:"price_max" description:"filter by maximum inclusive price" validate:"omitempty,max=15"`
}

// CreateFruit add a new fruit to the market.
func CreateVegetable(c *gin.Context,  v *Vegetable) (*Vegetable, error) {
	market.Lock()
	defer market.Unlock()

	n := strings.ToLower(v.Name)
	if _, ok := market.vegetables[n]; ok {
		return nil, errors.AlreadyExistsf("fruit")
	}
	v.AddedAt = time.Now()

	market.vegetables[n] = v

	return v, nil
}

// DeleteFruit removes a fruit from the market.
func DeleteVegetable(c *gin.Context, params *VegetableIdentityParams) error {
	if params.APIKey == "" {
		return errors.Forbiddenf("invalid api key")
	}
	market.Lock()
	defer market.Unlock()

	n := strings.ToLower(params.Name)
	if _, ok := market.fruits[n]; !ok {
		return errors.NotFoundf("fruit")
	}
	delete(market.fruits, n)

	return nil
}

// ListFruits lists the fruits of the market.
// Parameters can be used to filter the fruits.
func ListVegetables(c *gin.Context, params *ListVegetablesParams) ([]*Fruit, error) {
	basket := make([]*Fruit, 0)

	market.Lock()
	for _, f := range market.fruits {
		tobasket := true
		if params.Origin != nil && f.Origin != *params.Origin {
			tobasket = false
		}
		if params.PriceMin != nil && f.Price < *params.PriceMin {
			tobasket = false
		}
		if params.PriceMax != nil && f.Price > *params.PriceMax {
			tobasket = false
		}
		// If all conditions validates, add the
		// fruit to the returned basket.
		if tobasket {
			basket = append(basket, f)
		}
	}
	market.Unlock()

	c.Header("X-Market-Listing-Size", strconv.Itoa(len(basket)))

	return basket, nil
}


