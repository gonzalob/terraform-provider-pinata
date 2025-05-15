package provider

import (
	"context"
	"os"
	"terraform-provider-pinata/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ provider.Provider = &pinataProvider{}
)

type pinataProvider struct {
	version string
}

type PinataProviderModel struct {
	Root  types.String `tfsdk:"root"`
	Token types.String `tfsdk:"token"`
}

func (p *pinataProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "pinata"
	resp.Version = p.version
}

func (p *pinataProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"root": schema.StringAttribute{
				Optional: true,
			},
			"token": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *pinataProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config PinataProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Root.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("root"),
			"Unknown Pinata API root",
			"The provider cannot create the Pinata client as there is an unknown configuration value for the Pinata API root. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the PINATA_ROOT environment variable.",
		)
	}

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown Pinata JWT token",
			"The provider cannot create the Pinata client as there is an unknown configuration value for the Pinata JWT token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the PINATA_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	root := os.Getenv("PINATA_ROOT")
	token := os.Getenv("PINATA_TOKEN")

	if !config.Root.IsNull() {
		root = config.Root.ValueString()
	}

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing Pinata JWT token",
			"The provider cannot create the Pinata client as there is a missing or empty value for the Pinata JWT token. "+
				"Set the token value in the configuration or use the PINATA_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "pinata_root", root)
	ctx = tflog.SetField(ctx, "pinata_token", token)
	tflog.Debug(ctx, "Initializing client")
	client, err := client.NewClient(&root, &token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Pinata client",
			"An unexpected error occurred when creating the Pinata client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Pinata client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *pinataProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewPinResource,
	}
}

func (p *pinataProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &pinataProvider{
			version: version,
		}
	}
}
