package config

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// safeIdentifierFromProvider is like config.IdentifierFromProvider but
// returns ("", nil) instead of error when id is missing/empty in tfstate.
// This works around an upjet v2.2.0 bug where the framework client's
// setExternalName() doesn't guard against empty id before calling
// GetExternalNameFn (unlike the SDK client which does).
var safeIdentifierFromProvider = config.ExternalName{
	SetIdentifierArgumentFn: config.IdentifierFromProvider.SetIdentifierArgumentFn,
	GetExternalNameFn: func(tfstate map[string]any) (string, error) {
		id, ok := tfstate["id"].(string)
		if !ok || id == "" {
			return "", nil
		}
		return id, nil
	},
	GetIDFn:                config.IdentifierFromProvider.GetIDFn,
	OmittedFields:          config.IdentifierFromProvider.OmittedFields,
	DisableNameInitializer: config.IdentifierFromProvider.DisableNameInitializer,
	IdentifierFields:       config.IdentifierFromProvider.IdentifierFields,
}

// ExternalNameConfigs contains all external name configurations for this
// provider.
var ExternalNameConfigs = map[string]config.ExternalName{
	// SigNoz dashboard uses provider-assigned IDs (UUID from API)
	"signoz_dashboard": safeIdentifierFromProvider,
	// SigNoz alert uses provider-assigned IDs (integer from API)
	"signoz_alert": safeIdentifierFromProvider,
}

// TerraformPluginFrameworkExternalNameConfigs is for providers using
// terraform-plugin-framework (not SDK). The SigNoz provider uses framework.
var TerraformPluginFrameworkExternalNameConfigs = map[string]config.ExternalName{
	"signoz_dashboard": safeIdentifierFromProvider,
	"signoz_alert":     safeIdentifierFromProvider,
}

// ExternalNameConfigurations applies all external name configs listed in the
// table ExternalNameConfigs and sets the version of those resources to v1beta1
// assuming they will be tested.
func ExternalNameConfigurations() config.ResourceOption {
	return func(r *config.Resource) {
		if e, ok := ExternalNameConfigs[r.Name]; ok {
			r.ExternalName = e
		}
	}
}

// ExternalNameConfigured returns the list of all resources whose external name
// is configured manually.
func ExternalNameConfigured() []string {
	l := make([]string, len(ExternalNameConfigs))
	i := 0
	for name := range ExternalNameConfigs {
		// $ is added to match the exact string since the format is regex.
		l[i] = name + "$"
		i++
	}
	return l
}
