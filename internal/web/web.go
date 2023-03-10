package web

import (
	"fmt"
	"html"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sarunask/triviadb-gui/internal/env"
	"github.com/sarunask/triviadb-gui/internal/triviadb"
)

var db []triviadb.Result

func Run() error {
	r := gin.Default()
	r.LoadHTMLFiles(
		fmt.Sprintf("%s/main.tmpl", env.Settings.TemplatesDir),
		fmt.Sprintf("%s/error.tmpl", env.Settings.TemplatesDir),
		fmt.Sprintf("%s/question.tmpl", env.Settings.TemplatesDir),
		fmt.Sprintf("%s/answer.tmpl", env.Settings.TemplatesDir),
		fmt.Sprintf("%s/summary.tmpl", env.Settings.TemplatesDir),
	)
	// r.StaticFile("/", fmt.Sprintf("%s/css/main.tmpl", env.Settings.TemplatesDir))
	r.Static("/js", filepath.Join(env.Settings.TemplatesDir, "js"))
	r.GET("/", root)
	r.POST("/start", start)
	r.POST("/answer", answer)
	r.POST("/question", question)
	return r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func root(c *gin.Context) {
	c.HTML(http.StatusOK, "main.tmpl", gin.H{})
}

func start(c *gin.Context) {
	url := c.PostForm("url")
	if len(url) == 0 {
		c.HTML(http.StatusUnprocessableEntity, "error.tmpl", gin.H{
			"err": "Empty URL",
		})
		return
	}
	var err error
	db, err = triviadb.GetResults(url)
	if err != nil {
		c.HTML(http.StatusUnprocessableEntity, "error.tmpl", gin.H{
			"err": err.Error(),
		})
		return
	}
	c.HTML(http.StatusOK, "question.tmpl", gin.H{
		"nr":       0,
		"question": html.UnescapeString(db[0].Question),
		"answers":  db[0].GetAnswers(),
	})
}

func summary(c *gin.Context) {
	correct := 0
	incorrect := 0
	for _, question := range db {
		if question.CustomerAnswer == question.CorrectAnswer {
			correct++
		} else {
			incorrect++
		}
	}
	c.HTML(http.StatusOK, "summary.tmpl", gin.H{
		"total":     len(db),
		"correct":   correct,
		"incorrect": incorrect,
	})
}

func question(c *gin.Context) {
	sNr := c.PostForm("number")
	nr, err := strconv.ParseUint(sNr, 10, 16)
	if err != nil {
		c.HTML(http.StatusUnprocessableEntity, "error.tmpl", gin.H{
			"err": fmt.Sprintf("incorrect question number %s", sNr),
		})
		return
	}
	if nr >= uint64(len(db)) {
		summary(c)
		return
	}
	c.HTML(http.StatusOK, "question.tmpl", gin.H{
		"nr":       nr,
		"question": html.UnescapeString(db[nr].Question),
		"answers":  db[nr].GetAnswers(),
	})
}

func answer(c *gin.Context) {
	qNr := c.PostForm("number")
	nr, err := strconv.ParseUint(qNr, 10, 16)
	if err != nil {
		c.HTML(http.StatusUnprocessableEntity, "error.tmpl", gin.H{
			"err": fmt.Sprintf("incorrect question number %s", qNr),
		})
		return
	}
	db[nr].CustomerAnswer = c.PostForm("options")
	if nr >= uint64(len(db)) {
		c.HTML(http.StatusUnprocessableEntity, "error.tmpl", gin.H{
			"err": fmt.Sprintf("we have only %d answers, you asked for %d", len(db), nr+1),
		})
		return
	}
	c.HTML(http.StatusOK, "answer.tmpl", gin.H{
		"nr":       nr,
		"question": html.UnescapeString(db[nr].Question),
		"answer":   db[nr].GetCorrectAnswer(),
	})
}
