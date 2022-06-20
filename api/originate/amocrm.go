package originate

import (
	"log"
	"net/http"
	"strings"

	"asterisk-dialer/config"

	"github.com/gin-gonic/gin"
	"github.com/heltonmarx/goami/ami"
)

type (
	AmoCrmRequest struct {
		Login  string `form:"_login" binding:"required"`
		Secret string `form:"_secret" binding:"required"`
		Action string `form:"_action" binding:"required"`

		CallerName string `form:"as"`
		From       string `form:"from"`
		To         string `form:"to"`

		Oid string `form:"rand"`
	}

	AmoCrmResponse struct {
		Status string `json:"status"`
		Action string `json:"action"`
		Data   string `json:"data"`
	}
)

// Эмуляция скриптов http://www.voxlink.ru/kb/asterisk-configuration/amocrm-asterisk/
func AmoCrm(c *gin.Context) {
	ami_socket := c.MustGet("ami").(*ami.Socket)
	cnf := c.MustGet("cnf").(*config.Api)

	var r AmoCrmRequest
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("%+v", r)

	if cnf.AmoCrm.Login != r.Login || cnf.AmoCrm.Password != r.Secret {
		c.Status(http.StatusForbidden)
		return
	}

	uuid, err := ami.GetUUID()
	if err != nil {
		log.Fatalf("Get UUID: %v\n", err)
	}

	var res ami.Response
	response := AmoCrmResponse{
		Status: "error",
	}
	switch r.Action {
	case "status":
		res, err = ami.Status(ami_socket, uuid, "", "")
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		if res.Get("Response") == "Success" {
			response.Status = "ok"
			response.Data = res.Get("Message")
		}
	case "call":
		r.To = cleanNumber(r.To)

		if r.To != "" && r.From != "" {
			originate_data := ami.OriginateData{
				Channel:  "Local/" + r.From + "@" + cnf.OriginateData.OutgoingContext,
				Exten:    r.To,
				Context:  cnf.OriginateData.Context,
				Priority: 1,
				Timeout:  cnf.OriginateData.Timeout,
				CallerID: `"` + r.CallerName + ` ` + r.To + `" <` + r.From + `>`,
				Variable: "FOO=1",
				Account:  cnf.OriginateData.Account,
				Codecs:   cnf.OriginateData.Codecs,
				Async:    "true",
			}

			originate_data.Variable += ",PHONE=" + r.To

			if r.Oid != "" {
				originate_data.Variable += ",OID=" + r.Oid
			}
			if cnf.OriginateData.Variable != "" {
				originate_data.Variable += "," + cnf.OriginateData.Variable
			}

			log.Printf("%+v", originate_data)

			res, err = ami.Originate(ami_socket, uuid, originate_data)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return
			}

			log.Printf("%+v", res)

			if res.Get("Response") == "Success" {
				response.Status = "ok"
				response.Data = res.Get("Message")
			}
		}
	}

	c.JSON(
		http.StatusOK,
		response,
	)
}

func cleanNumber(num string) string {
	num = strings.ReplaceAll(num, "(", "")
	num = strings.ReplaceAll(num, ")", "")
	num = strings.ReplaceAll(num, "-", "")
	num = strings.ReplaceAll(num, " ", "")

	return num
}
