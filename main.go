package main

import (
	"fmt"
	"os"
	"github.com/urfave/cli"
	"github.com/evcraddock/article-importer/config"
	"github.com/evcraddock/article-importer/tasks"
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
	      Destination: &configSettings.UserName,
	    },
	    cli.StringFlag{
	      Name:        "password",
	      Value:       "",
	      Usage:       "",
	      Destination: &configSettings.Password,
	    },
	    cli.StringFlag{
	      Name:        "serviceUrl",
	      Value:       configSettings.ServiceUrl,
	      Usage:       "",
	      Destination: &configSettings.ServiceUrl,
	    },
	}

	app.Commands = []cli.Command{
	    {
	      Name:  "new-article",
	      Usage: "create a new article",
	      Action: func(c *cli.Context) error {
	      	task := tasks.NewTask(configSettings)
	        article, err := task.CreateNewArticle()
	        if err != nil {
	        	return cli.NewExitError(err.Error(), 86)
	        }

	        fmt.Printf("Successfull Created Article %s on %v\n", article.Title, article.PublishDate)
	        return nil
	      },
	    },
	    {
	      Name:  "update-article",
	      Usage: "update an existing article",
	      Action: func(c *cli.Context) error {
	      	task := tasks.NewTask(configSettings)
	        article, err := task.UpdateArticle()
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
	        articleId, err := task.DeleteArticle()
	        if err != nil {
	        	return cli.NewExitError(err.Error(), 86)
	        }

	        fmt.Printf("Successfull Deleted Article %s \n", articleId)
	        return nil
	      },
	    },
	}

	app.Run(os.Args)
	
}

