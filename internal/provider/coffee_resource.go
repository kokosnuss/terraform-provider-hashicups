package provider

import (
	"context"
	"fmt"
	"slices"
	"strconv"

	"github.com/hashicorp-demoapp/hashicups-client-go"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &coffeeResource{}
	_ resource.ResourceWithConfigure   = &coffeeResource{}
	_ resource.ResourceWithImportState = &coffeeResource{}
)

// NewcoffeeResource is a helper function to simplify the provider implementation.
func NewCoffeeResource() resource.Resource {
	return &coffeeResource{}
}

// coffeeResource is the resource implementation.
type coffeeResource struct {
	client *hashicups.Client
}

// coffeeResourceModel maps the resource schema data.
type coffeeResourceModel struct {
	ID          types.String      `tfsdk:"id"`
	Name        types.String      `tfsdk:"name"`
	Teaser      types.String      `tfsdk:"teaser"`
	Collection  types.String      `tfsdk:"collection"`
	Origin      types.String      `tfsdk:"origin"`
	Color       types.String      `tfsdk:"color"`
	Description types.String      `tfsdk:"description"`
	Price       types.Int64       `tfsdk:"price"`
	Image       types.String      `tfsdk:"image"`
	Ingredients []ingredientModel `tfsdk:"ingredients"`
}

// ingredientModel maps the nested ingredient schema data.
type ingredientModel struct {
	IngredientID types.Int64   `tfsdk:"ingredient_id"`
	Name         types.String  `tfsdk:"name"`
	Quantity     types.Float64 `tfsdk:"quantity"`
	Unit         types.String  `tfsdk:"unit"`
}

// Metadata returns the resource type name.
func (r *coffeeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_coffee"
}

// Schema defines the schema for the coffee resource.
func (r *coffeeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a coffee.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Numeric identifier of the coffee.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the coffee.",
				Required:    true,
			},
			"teaser": schema.StringAttribute{
				Description: "Short teaser text for the coffee.",
				Optional:    true,
			},
			"collection": schema.StringAttribute{
				Description: "Collection the coffee belongs to.",
				Optional:    true,
			},
			"origin": schema.StringAttribute{
				Description: "Origin or release season of the coffee.",
				Optional:    true,
			},
			"color": schema.StringAttribute{
				Description: "Color code associated with the coffee.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "Detailed description of the coffee.",
				Optional:    true,
			},
			"price": schema.Int64Attribute{
				Description: "Price of the coffee in USD.",
				Required:    true,
			},
			"image": schema.StringAttribute{
				Description: "URL or path to the coffee image.",
				Optional:    true,
			},
			"ingredients": schema.ListNestedAttribute{
				Description: "List of ingredients in the coffee.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ingredient_id": schema.Int64Attribute{
							Description: "Identifier of the ingredient.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the ingredient.",
							Optional:    true,
						},
						"quantity": schema.Float64Attribute{
							Description: "Quantity of the ingredient.",
							Required:    true,
						},
						"unit": schema.StringAttribute{
							Description: "Unit of measurement for the ingredient.",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

// Create a new resource.
func (r *coffeeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan coffeeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hashiCoffee := hashicups.Coffee{
		Name:   plan.Name.ValueString(),
		Teaser: plan.Teaser.ValueString(),
		Price:  float64(plan.Price.ValueInt64()),
		Image:  plan.Image.ValueString(),
	}
	c, err := r.client.CreateCoffee(hashiCoffee)
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	for _, ingredient := range plan.Ingredients {
		hashiIngredient := hashicups.Ingredient{
			Name:     ingredient.Name.ValueString(),
			Quantity: int(ingredient.Quantity.ValueFloat64()),
			Unit:     ingredient.Unit.ValueString(),
		}
		hi, err := r.client.CreateCoffeeIngredient(*c, hashiIngredient)
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
			return
		}
		ingredient.IngredientID = types.Int64Value(int64(hi.ID))
	}
	tflog.Info(ctx, fmt.Sprintf("c: %v", c))

	// Set state to fully populated data
	plan.ID = types.StringValue(strconv.Itoa(c.ID))
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *coffeeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state coffeeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	c, err := r.client.GetCoffee(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	ingredients, err := r.client.GetCoffeeIngredients(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	stateIngrediens := []ingredientModel{}
	for _, ingredient := range ingredients {

		stateIngrediens = append(stateIngrediens, ingredientModel{
			IngredientID: types.Int64Value(int64(ingredient.ID)),
			Name:         types.StringValue(ingredient.Name),
			Quantity:     types.Float64Value(float64(ingredient.Quantity)),
			Unit:         types.StringValue(ingredient.Unit),
		})

	}

	state.Ingredients = stateIngrediens

	state.ID = types.StringValue(strconv.Itoa(c.ID))
	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *coffeeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state coffeeResourceModel
	var plan coffeeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	// Get current ingredients and set to null if not present in plan
	req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, _ := strconv.Atoi(plan.ID.ValueString())
	hashiCoffe := hashicups.Coffee{
		ID:     id,
		Name:   plan.Name.ValueString(),
		Teaser: plan.Teaser.ValueString(),
		Price:  float64(plan.Price.ValueInt64()),
		Image:  plan.Image.ValueString(),
	}
	c, err := r.client.UpdateCoffee(hashiCoffe)
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	collectedIngredients := r.CollectIngredientModels(ctx, state.Ingredients, plan.Ingredients)
	for _, hashiIngredient := range collectedIngredients {
		hi, err := r.client.CreateCoffeeIngredient(*c, hashiIngredient)
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
			return
		}
		tflog.Info(ctx, fmt.Sprintf("hi: %v\n", hi))

		planIndex := slices.IndexFunc(plan.Ingredients, func(i ingredientModel) bool { return i.Name == types.StringValue(hi.Name) })
		if planIndex > -1 {
			plan.Ingredients[planIndex].IngredientID = types.Int64Value(int64(hi.ID))
		}

	}

	tflog.Info(ctx, fmt.Sprintf("c: %v", c))

	// Set state to fully populated data
	plan.ID = types.StringValue(strconv.Itoa(c.ID))
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *coffeeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state coffeeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	err := r.client.DeleteCoffee(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting HashiCups Coffee",
			"Could not delete coffee, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *coffeeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*hashicups.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *coffeeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *coffeeResource) CollectIngredientModels(ctx context.Context, stateIngredients, planIngredients []ingredientModel) []hashicups.Ingredient {

	collectedIngredients := []hashicups.Ingredient{}
	for _, stateIngredient := range stateIngredients {
		planIndex := slices.IndexFunc(planIngredients, func(i ingredientModel) bool { return i.Name == stateIngredient.Name })
		var ingredientToUse *ingredientModel
		if planIndex < 0 {
			// Delete
			ingredientToUse = &stateIngredient
			ingredientToUse.Quantity = types.Float64Value(0)
		} else {
			// Update
			ingredientToUse = &planIngredients[planIndex]
		}

		collectedIngredients = append(collectedIngredients, hashicups.Ingredient{
			Name:     ingredientToUse.Name.ValueString(),
			Quantity: int(ingredientToUse.Quantity.ValueFloat64()),
			Unit:     ingredientToUse.Unit.ValueString(),
		})
	}

	// Create New Ingredients
	for _, planIngredient := range planIngredients {
		stateIndex := slices.IndexFunc(stateIngredients, func(i ingredientModel) bool { return i.Name == planIngredient.Name })
		if stateIndex > 0 {
			// We handled this one already
			continue
		}
		collectedIngredients = append(collectedIngredients, hashicups.Ingredient{
			Name:     planIngredient.Name.ValueString(),
			Quantity: int(planIngredient.Quantity.ValueFloat64()),
			Unit:     planIngredient.Unit.ValueString(),
		})
	}

	return collectedIngredients
}
