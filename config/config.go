package config

import (
	"html/template"

	"github.com/gin-gonic/gin"
)

type (
	Api struct {
		OriginateData OriginateData `yaml:"originate_data"`
		AmoCrm        AmoCrm        `yaml:"amocrm"`
		Template      *template.Template
	}

	OriginateData struct {
		OutgoingContext string `yaml:"outgoing_context"`
		Exten           string `yaml:"exten"`
		Context         string `yaml:"context"`
		Priority        int    `yaml:"priority"`
		Application     string `yaml:"application"`
		Data            string `yaml:"data"`
		Timeout         int    `yaml:"timeout"`
		CallerID        string `yaml:"callerid"`
		Variable        string `yaml:"variable"`
		Account         string `yaml:"account"`
		Codecs          string `yaml:"codecs"`
	}

	AmoCrm struct {
		Login    string `yaml:"login"`
		Password string `yaml:"password"`
	}
)

func Inject(cnf *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("cnf", cnf)
	}
}

func LoadTemplate(cnf *Api, pattern string) {
	cnf.Template = template.Must(template.New("").ParseGlob(pattern))
}
