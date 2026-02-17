package dashboard

import "github.com/crossplane/upjet/v2/pkg/config"

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("signoz_dashboard", func(r *config.Resource) {
		r.ShortGroup = "dashboard"
		r.Kind = "Dashboard"
	})
}
