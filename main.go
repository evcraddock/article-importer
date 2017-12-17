package main

import (
	"fmt"
	"os"

	"github.com/evcraddock/article-importer/config"
	"github.com/evcraddock/article-importer/tasks"
	"github.com/urfave/cli"
)

func main() {
	configSettings := config.NewConfiguration()
	app := cli.NewApp()
	app.Name = "Article Importer"
	app.Version = "1.0"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Erik Craddock",
			Email: "erik@craddock.org",
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "username",
			Value:       "",
			Usage:       "",
			Destination: &configSettings.Auth.UserName,
		},
		cli.StringFlag{
			Name:        "password",
			Value:       "",
			Usage:       "",
			Destination: &configSettings.Auth.Password,
		},
		cli.StringFlag{
			Name:        "serviceUrl",
			Value:       configSettings.Auth.ServiceURL,
			Usage:       "",
			Destination: &configSettings.Auth.ServiceURL,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "load-article",
			Usage: "load article from yaml file",
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "force, f"},
				cli.StringFlag{Name: "filename"},
			},
			Action: func(c *cli.Context) error {
				task := tasks.NewTask(configSettings)
				article, err := task.LoadArticle(c.String("filename"), c.Bool("force"))
				if err != nil {
					return cli.NewExitError("Error Message: "+err.Error(), 86)
				}

				fmt.Printf("Successfull Loaded Article %s (Id: %s) on %v\n", article.Title, article.ID, article.PublishDate)
				return nil
			},
		},
		{
			Name:  "update-article",
			Usage: "update an existing article",
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "force, f"},
			},
			Action: func(c *cli.Context) error {
				task := tasks.NewTask(configSettings)
				article, err := task.UpdateArticle(c.Bool("force"))
				if err != nil {
					return cli.NewExitError(err.Error(), 86)
				}

				fmt.Printf("Successfull Updated Article %s on %v\n", article.Title, article.PublishDate)
				return nil
			},
		},
		{
			Name:  "delete-article",
			Usage: "delete an existing article",
			Action: func(c *cli.Context) error {
				task := tasks.NewTask(configSettings)
				articleID, err := task.DeleteArticle()
				if err != nil {
					return cli.NewExitError(err.Error(), 86)
				}

				fmt.Printf("Successfull Deleted Article %s \n", articleID)
				return nil
			},
		},
		{
			Name:  "new-link",
			Usage: "create a new link",
			Action: func(c *cli.Context) error {
				task := tasks.NewTask(configSettings)
				link, err := task.CreateNewLink()
				if err != nil {
					return cli.NewExitError(err.Error(), 86)
				}

				fmt.Printf("Successfull Link %s \n", link.Title)
				return nil
			},
		},
		{
			Name:  "delete-link",
			Usage: "delete an existing link",
			Action: func(c *cli.Context) error {
				task := tasks.NewTask(configSettings)
				linkID, err := task.DeleteLink()
				if err != nil {
					return cli.NewExitError(err.Error(), 86)
				}

				fmt.Printf("Successfull Deleted link %s \n", linkID)
				return nil
			},
		},
	}

	app.Run(os.Args)

}
