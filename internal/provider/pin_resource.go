package provider

import (
	"context"
	"fmt"
	"terraform-provider-pinata/internal/client"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &pinResource{}
	_ resource.ResourceWithConfigure   = &pinResource{}
	_ resource.ResourceWithImportState = &pinResource{}
)

func NewPinResource() resource.Resource {
	return &pinResource{}
}

type pinResourceModel struct {
	ID      types.String `tfsdk:"id"`
	CID     types.String `tfsdk:"cid"`
	Name    types.String `tfsdk:"name"`
	Version types.Number `tfsdk:"version"`
	Paths   []Path       `tfsdk:"paths"`
}

type pinResource struct {
	client *client.Client
}

type Path struct {
	Name types.String `tfsdk:"name"`
	Hash types.String `tfsdk:"hash"`
}

func (r *pinResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pin"
}

func (r *pinResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an IPFS pin",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Pinata ID for the pin",
				Computed:    true,
			},
			"cid": schema.StringAttribute{
				Description: "The pin's IPFS Content ID",
				Computed:    true,
			},
			"version": schema.NumberAttribute{
				Description: "The CID version to use",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Resource name",
				Optional:    true,
				Computed:    true,
			},
			"paths": schema.ListNestedAttribute{
				Description: "Local paths for the pin",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Path to each of the assets",
							Required:    true,
						},
						"hash": schema.StringAttribute{
							Description: "Resource checksum",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func (r *pinResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan pinResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name
	if name.IsNull() || name.IsUnknown() {
		name = types.StringValue(fmt.Sprintf("terraform-%d", time.Now().UnixMilli()))
	}

	var files []string
	for _, path := range plan.Paths {
		files = append(files, path.Name.ValueString())
	}

	pin, err := r.client.PinFolder(files, name.ValueString(), plan.Version.String())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error pinning",
			"Could not pin resources, unexpected error: "+err.Error(),
		)
		return
	}
	plan.ID = types.StringValue(pin.ID)
	plan.CID = types.StringValue(pin.IPFSHash)
	plan.Name = types.StringValue(pin.Name)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *pinResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state pinResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	pin, err := r.client.GetPinById((state.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Pinata pin",
			"Could not read Pinata pin with ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state.ID = types.StringValue(pin.Data.ID)
	state.CID = types.StringValue(pin.Data.CID)
	state.Name = types.StringValue(pin.Data.Name)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *pinResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan pinResourceModel
	var state pinResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Unpin(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error unpinning pin",
			"Could not unpin, unexpected error: "+err.Error(),
		)
		return
	}

	var files []string
	var name types.String
	name = plan.Name
	if name.IsNull() {
		name = state.Name
	}

	pin, err := r.client.PinFolder(files, name.ValueString(), plan.Version.String())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error pinning",
			"Could not pin resources, unexpected error: "+err.Error(),
		)
		return
	}
	plan.ID = types.StringValue(pin.ID)
	plan.CID = types.StringValue(pin.IPFSHash)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *pinResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state pinResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Unpin(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error unpinning pin",
			"Could not unpin, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *pinResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected pinata client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	ctx = tflog.SetField(ctx, "client", client)
	tflog.Info(ctx, "configured")
	r.client = client
}

func (r *pinResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
