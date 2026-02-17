package alert

import "github.com/crossplane/upjet/v2/pkg/config"

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("signoz_alert", func(r *config.Resource) {
		r.ShortGroup = "alert"
		r.Kind = "Alert"
	})
}
