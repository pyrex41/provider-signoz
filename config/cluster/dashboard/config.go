package dashboard

import "github.com/crossplane/upjet/v2/pkg/config"

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("signoz_dashboard", func(r *config.Resource) {
		r.ShortGroup = "dashboard"
		r.Kind = "Dashboard"

		// Prevent LateInitialize from adopting server-mutated JSON fields
		// into spec.forProvider. SigNoz's POST handler runs a v5 migration
		// that mutates widget JSON in breaking ways ($var → {{.var}},
		// injects #SIGNOZ_VALUE orderBy, changes op:"in" → op:"=").
		// Without this, Crossplane adopts the broken v5 state as the
		// desired state, making it impossible to fix via GitOps.
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{
				"widgets",
				"variables",
				"layout",
				"panel_map",
				"version",
			},
		}
	})
}
