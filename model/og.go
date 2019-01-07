package model

import (
	"net/http"
	"messenger/helper"
	"github.com/graphql-go/graphql"
	"strings"
	"errors"
)

type Og struct {
	Title       string `meta:"og:title" json:"title"`
	Description string `meta:"og:description,description" json:"description"`
	Type        string `meta:"og:type" json:"type"`
	URL         string `meta:"og:url" json:"url"`
	Image       string `meta:"og:image" json:"image"`
}

var OgType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Og",
		Fields: graphql.Fields{
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"type": &graphql.Field{
				Type: graphql.String,
			},
			"image": &graphql.Field{
				Type: graphql.String,
			},
			"url": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

func GetOgTag(urlString string) (error, *Og) {

	if strings.HasPrefix(urlString, "http://localhost") ||
		strings.HasPrefix(urlString, "localhost") ||
		strings.HasPrefix(urlString, "https://localhost") ||
		strings.HasPrefix(urlString, "http://127.0.1") ||
		strings.HasPrefix(urlString, "https://127.0.1") ||
		strings.HasPrefix(urlString, "127.0.0.1") {
		return errors.New("invalid url"), nil
	}

	res, err := http.Get(urlString)
	if err != nil {
		return err, nil
	}
	data := new(Og)
	e := helper.Ogtag(res.Body, data)
	if e != nil {
		return e, nil
	}

	return nil, data

}
