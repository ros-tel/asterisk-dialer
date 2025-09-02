package originate

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"

	"asterisk-dialer/config"

	"github.com/gin-gonic/gin"
	"github.com/heltonmarx/goami/ami"
)

type (
	originateRequest struct {
		Phone    string `form:"phone"`
		Template string `form:"template"`

		Channel     string `form:"channel"`
		Exten       string `form:"exten"`
		Context     string `form:"context"`
		Priority    int    `form:"priority" binding:"omitempty,min=1,max=1000"`
		Application string `form:"application"`
		Data        string `form:"data"`
		Timeout     int    `form:"timeout" binding:"omitempty,min=5000,max=180000"`
		CallerID    string `form:"callerid"`
		Variable    string `form:"variable"`
		Account     string `form:"account"`
		Codecs      string `form:"codecs"`
		Oid         string `form:"oid"`

		Raw []string `form:"raw"`
		Num []string `form:"num"`
	}

	originateResponse struct {
		ActionID string `json:"action_id"`
		Message  string `json:"message"`
		Response string `json:"response"`
	}

	task struct {
		Raw0       string
		Raw1       string
		Raw2       string
		Raw3       string
		Raw4       string
		SayNumber0 string
		SayNumber1 string
		SayNumber2 string
		SayNumber3 string
		SayNumber4 string
	}
)

func Originate(c *gin.Context) {
	ami_socket := c.MustGet("ami").(*ami.Socket)
	cnf := c.MustGet("cnf").(*config.Api)

	var r originateRequest
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("%+v", r)

	if r.Channel != "" && r.Phone != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Channel and Phone cannot be used at the same time"})
		return
	}

	uuid, err := ami.GetUUID()
	if err != nil {
		log.Fatalf("Get UUID: %v\n", err)
	}

	originate_data := ami.OriginateData{
		Channel:     "Local/" + r.Phone + "@" + cnf.OriginateData.OutgoingContext,
		Exten:       cnf.OriginateData.Exten,
		Context:     cnf.OriginateData.Context,
		Priority:    cnf.OriginateData.Priority,
		Application: cnf.OriginateData.Application,
		Data:        cnf.OriginateData.Data,
		Timeout:     cnf.OriginateData.Timeout,
		CallerID:    cnf.OriginateData.CallerID,
		Account:     cnf.OriginateData.Account,
		Codecs:      cnf.OriginateData.Codecs,
		Async:       "true",
	}
	if r.Channel != "" {
		originate_data.Channel = r.Channel
	}
	if r.Exten != "" {
		originate_data.Exten = r.Exten
	}
	if r.Context != "" {
		originate_data.Context = r.Context
	}
	if r.Priority != 0 {
		originate_data.Priority = r.Priority
	}
	if r.Application != "" {
		originate_data.Application = r.Application
	}
	if r.Data != "" {
		originate_data.Data = r.Data
	}
	if r.Timeout != 0 {
		originate_data.Timeout = r.Timeout
	}
	if r.CallerID != "" {
		originate_data.CallerID = r.CallerID
	}
	if r.Account != "" {
		originate_data.Account = r.Account
	}
	if r.Codecs != "" {
		originate_data.Codecs = r.Codecs
	}

	if r.Phone != "" {
		originate_data.Variable = append(originate_data.Variable, "PHONE="+r.Phone)
	}
	if r.Oid != "" {
		originate_data.Variable = append(originate_data.Variable, "OID="+r.Oid)
	}
	if cnf.OriginateData.Variable != "" {
		originate_data.Variable = append(originate_data.Variable, strings.Split(cnf.OriginateData.Variable, ",")...)
	}

	if r.Template != "" && (len(r.Raw) > 0 || len(r.Num) > 0) {
		var t task
		for k, v := range r.Raw {
			switch k {
			case 0:
				t.Raw0 = v
			case 1:
				t.Raw1 = v
			case 2:
				t.Raw2 = v
			case 3:
				t.Raw3 = v
			case 4:
				t.Raw4 = v
			}

		}
		for k, v := range r.Num {
			switch k {
			case 0:
				t.SayNumber0 = say(v)
			case 1:
				t.SayNumber1 = say(v)
			case 2:
				t.SayNumber2 = say(v)
			case 3:
				t.SayNumber3 = say(v)
			case 4:
				t.SayNumber4 = say(v)
			}

		}

		log.Printf("%+v", t)

		var buf bytes.Buffer
		writer := bufio.NewWriter(&buf)

		err = cnf.Template.ExecuteTemplate(writer, r.Template+".tpl", t)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = writer.Flush()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		sound, err := buf.ReadString(0)
		if err != nil {
			if err != io.EOF {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		sound = strings.Replace(sound, "&amp;", "&", -1)
		originate_data.Variable = append(originate_data.Variable, "SOUND="+strings.TrimSpace(sound))
	}

	log.Printf("%+v", originate_data)

	originate, err := ami.Originate(c, ami_socket, uuid, originate_data)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Printf("%+v", originate)

	code := http.StatusOK
	response := originate.Get("Response")
	if response != "Success" {
		code = http.StatusInternalServerError
	}
	c.JSON(code, originateResponse{
		ActionID: originate.Get("ActionID"),
		Message:  originate.Get("Message"),
		Response: response,
	})
}

func say(num string) string {
	var say_buf []string
	for _, s := range []rune(num) {
		say_buf = append(say_buf, "digits/"+string(s))
	}

	return strings.Join(say_buf, "&")
}
