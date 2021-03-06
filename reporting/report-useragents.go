package reporting

import (
	"bytes"
	"html/template"
	"os"

	"github.com/ocmdev/rita/database"
	"github.com/ocmdev/rita/datatypes/useragent"
	"github.com/ocmdev/rita/reporting/templates"
)

func printUserAgents(db string, res *database.Resources) error {
	f, err := os.Create("useragents.html")
	if err != nil {
		return err
	}
	defer f.Close()
	out, err := template.New("useragents.html").Parse(templates.UserAgentsTempl)
	if err != nil {
		return err
	}

	var agents []useragent.UserAgent
	coll := res.DB.Session.DB(db).C(res.System.UserAgentConfig.UserAgentTable)
	coll.Find(nil).Sort("times_used").Limit(1000).All(&agents)

	w, err := getUserAgentsWriter(agents)
	if err != nil {
		return err
	}
	return out.Execute(f, &templates.ReportingInfo{DB: db, Writer: template.HTML(w)})
}

func getUserAgentsWriter(agents []useragent.UserAgent) (string, error) {
	tmpl := "<tr><td>{{.UserAgent}}</td><td>{{.TimesUsed}}</td></tr>\n"
	out, err := template.New("Agents").Parse(tmpl)
	if err != nil {
		return "", err
	}
	w := new(bytes.Buffer)
	for _, agent := range agents {
		err := out.Execute(w, agent)
		if err != nil {
			return "", err
		}
	}
	return w.String(), nil
}
